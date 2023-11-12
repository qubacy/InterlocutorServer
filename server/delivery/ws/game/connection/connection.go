package connection

import (
	"context"
	"ilserver/pkg/utility"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	readingChan chan Message
	writingChan chan Message

	closed context.Context
	cancel context.CancelFunc

	id        string
	webSocket *websocket.Conn
}

func NewOpenConnection(id string, webSocket *websocket.Conn) *Connection {
	closed, cancel := context.WithCancel(context.Background())

	instance := &Connection{
		readingChan: make(chan Message),
		writingChan: make(chan Message),

		closed: closed,
		cancel: cancel,

		id:        id,
		webSocket: webSocket,
	}

	webSocket.SetCloseHandler(
		instance.closer)

	// ***

	go instance.reading()
	go instance.writing()

	return instance
}

// public
// -----------------------------------------------------------------------

func (c *Connection) Id() string {
	return c.id
}

func (c *Connection) Writer() chan<- Message {
	return c.writingChan
}

func (c *Connection) Reader() <-chan Message {
	return c.readingChan
}

func (c *Connection) Closed() context.Context {
	return c.closed
}

// -----------------------------------------------------------------------

func (c *Connection) CloseGracefully() {
	if c.IsClosed() {
		return
	}

	defer c.closeAll()

	c.writeCloseMessage()

	// server -----> client
	// ..close connection..
	// server <--x-- client
}

func (c *Connection) IsClosed() bool {
	select {
	case <-c.closed.Done():
		return true
	default:
	}
	return false
}

// private
// -----------------------------------------------------------------------

func (c *Connection) closer(code int, text string) error {
	log.Println("call close handler")

	defer c.closeAll()
	c.writeCloseMessage()

	// server <----- client
	// server -----> client
	// ..close connection..

	return nil
}

func (c *Connection) reading() {
	defer c.closeAll()

	for {
		messageType, data, err := c.webSocket.ReadMessage()
		if c.IsClosed() {
			return
		}

		if messageType == -1 { // something wrong.
			return
		}
		if messageType != websocket.TextMessage {
			log.Println(utility.CreateCustomError(c.reading, ErrMessageTypeNotText))
			return
		}

		if err != nil {
			log.Println(utility.CreateCustomError(c.reading, err))
			return
		}

		c.readingChan <- MakeMessage(messageType, data) // only text.
	}
}

func (c *Connection) writing() {
	defer c.closeAll()

	for {
		message := <-c.writingChan
		err := c.webSocket.WriteMessage(message.Type, message.Data)
		if c.IsClosed() {
			return
		}

		if err != nil {
			log.Println(utility.CreateCustomError(c.writing, err))
			return
		}
	}
}

func (c *Connection) closeAll() {
	log.Println("call close all")

	if c.IsClosed() {
		return
	}

	c.cancel()
	c.webSocket.Close() // ignore err.
}

func (c *Connection) writeCloseMessage() {
	deadline := time.Now().Add(time.Second)
	data := websocket.FormatCloseMessage(
		websocket.CloseNormalClosure, "")

	// The Close and WriteControl methods
	// can be called concurrently with all other methods.

	err := c.webSocket.WriteControl(websocket.CloseMessage, data, deadline)

	if err != nil {
		log.Println(utility.CreateCustomError(c.writeCloseMessage, err))
	}
}
