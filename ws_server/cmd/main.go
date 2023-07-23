package main

import (
	"encoding/json"
	"ilserver/service"
	"ilserver/transport"
	"ilserver/transport/dto"
	"ilserver/transport/overWs"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// TODO: поместить в handler
var roomService []service.RoomService

func main() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// TODO: ограничить время бездействия
	}

	// ***

	var handler transport.Handler
	handler.RoomService.BackgroundWork(16 * time.Millisecond)

	// ***

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		websocket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Websocket Connected!")
		listen(&handler, websocket) // app, over ws
	})

	// ***

	http.HandleFunc("/debug/rooms", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := json.Marshal(handler.RoomService.Rooms)
		w.Header().Add("Content-Type", "application/json")
		w.Write(bytes)
	})

	// ***

	http.ListenAndServe(":47777", nil)
}

// -----------------------------------------------------------------------

func listen(handler *transport.Handler, conn *websocket.Conn) {
	for {
		messageType, messageContent, err := conn.ReadMessage()
		if messageType != websocket.TextMessage {
			log.Println(conn.RemoteAddr(), "message type is not text")
			conn.Close()
			return
		}

		// ***

		if err != nil {
			log.Println(conn.RemoteAddr(), err)
			conn.Close()
			return
		}

		// ***

		log.Println(conn.RemoteAddr(), string(messageContent))

		var pack overWs.Pack
		err = json.Unmarshal(messageContent, &pack)
		if err != nil {
			// TODO: правильная обработка ошибки синтаксиса пакета
			log.Println(conn.RemoteAddr(), err)
			conn.Close()
			return
		}
		log.Println(conn.RemoteAddr(), pack)

		// *** to handle
		handler.AddConn(conn)

		err = routeWsPack(handler, conn, pack)
		if err != nil {
			log.Println(conn.RemoteAddr(), err)
			conn.Close()
			return
		}
	}
}

// -----------------------------------------------------------------------

func routeWsPack(handler *transport.Handler, conn *websocket.Conn, pack overWs.Pack) error {
	if pack.Operation == transport.SEARCHING_START {
		bytes, err := json.Marshal(pack.RawBody)
		if err != nil {
			return err
		}

		var ssBody dto.CliSearchingStartBodyClient
		err = json.Unmarshal(bytes, &ssBody)
		if err != nil {
			return err
		}

		handler.SearchingStart(conn, ssBody)

	} else if pack.Operation == transport.SEARCHING_STOP {

	} else if pack.Operation == transport.CHATTING_NEW_MESSAGE {

	} else if pack.Operation == transport.CHOOSING_USERS_CHOSEN {

	}
	return handler.Err(conn, pack.Operation, "operation is unknown")
}
