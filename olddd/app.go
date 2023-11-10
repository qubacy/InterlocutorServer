package app

import (
	"encoding/json"
	"fmt"
	"ilserver/delivery/http/control"

	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

// -----------------------------------------------------------------------

func Run() {

	gameHandler := overWs.NewHandler()
	controlHandler :=

		// ***

		overWs.BackgroundUpdateRooms(wsHandler)

	// ***

	var mux = http.NewServeMux()

	// *** websocket and control ***

	prepareWsServer(mux, wsHandler)
	prepareControlServer(mux, controlHandler)

	// *** simple debug http ***

	var useDebugSvr bool = viper.GetBool(
		"debug_server.use")
	if useDebugSvr {
		prepareDebugServer(mux, wsHandler)
	}

	// ***

	s.httpServer = &http.Server{
		Addr:           ":" + viper.GetString("port"),
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        mux,
	}

	return s.httpServer.ListenAndServe()
}

// high level route preparation
// -----------------------------------------------------------------------

func prepareControlServer(mux *http.ServeMux, handler *control.Handler) {
	mux.HandleFunc("/sign-in", handler.SignIn) // POST
	mux.HandleFunc("/topic", handler.Topic)    // POST
	mux.HandleFunc("/topics", handler.Topics)  // POST, GET
}

func prepareWsServer(mux *http.ServeMux, handler *overWs.Handler) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		websocket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade err:", err)
			return
		}

		log.Println("Ws", websocket.RemoteAddr().String(), "connected")
		if err := listen(handler, websocket); err != nil {
			log.Println("One ws listen err:", err)
			websocket.Close()
			return
		}
	})
}

// websocket server details
// -----------------------------------------------------------------------

func listen(handler *overWs.Handler, conn *websocket.Conn) error {
	handler.AddConn(conn)

	for {
		messageType, messageContent, err := conn.ReadMessage()

		// TODO: изучить закрытие веб-сокета
		if err != nil {
			switch err.(type) {
			case *websocket.CloseError:
				concreteErr := err.(*websocket.CloseError)
				log.Printf("Unexpected read message, close err %v", concreteErr)
			case *websocket.HandshakeError:
				concreteErr := err.(*websocket.HandshakeError)
				log.Printf("Unexpected read message, handshake err %v", concreteErr)
			}
			return handler.RemoveConnAndClose(conn)
		}

		// ***

		log.Println(string(messageContent))
		log.Println(messageType)

		if messageType == websocket.CloseMessage {
			return handler.RemoveConnAndClose(conn)
		}

		if messageType != websocket.TextMessage {
			log.Println(conn.RemoteAddr(), "message type is not text")
			return handler.RemoveConnAndClose(conn)
		}

		if err != nil {
			log.Println(conn.RemoteAddr(), err)
			return handler.RemoveConnAndClose(conn)
		}

		// ***

		var pack overWsDto.Pack
		err = json.Unmarshal(messageContent, &pack)
		if err != nil {
			log.Println(conn.RemoteAddr(), err)

			// TODO: отправить пакет с информацией об ошибки

			return handler.RemoveConnAndClose(conn)
		}
		log.Println(conn.RemoteAddr(), pack)

		// ***

		err = routeWsPack(handler, conn, pack)
		if err != nil {
			log.Println(conn.RemoteAddr(), err)
			return handler.RemoveConnAndClose(conn)
		}
	}
}

func routeWsPack(handler *overWs.Handler, conn *websocket.Conn, pack overWsDto.Pack) error {
	bytes, err := json.Marshal(pack.RawBody)
	if err != nil {
		return err
	}

	// ***

	if pack.Operation == overWsDto.SEARCHING_START {
		var reqDto overWsDto.CliSearchingStartBodyClient
		err = json.Unmarshal(bytes, &reqDto)
		if err != nil {
			return err
		}

		return handler.SearchingStart(conn, reqDto)
	} else if pack.Operation == overWsDto.SEARCHING_STOP {
		var reqDto overWsDto.CliSearchingStopBody
		err = json.Unmarshal(bytes, &reqDto)
		if err != nil {
			return err
		}

		return handler.SearchingStop(conn, reqDto)
	} else if pack.Operation == overWsDto.CHATTING_NEW_MESSAGE {
		var reqDto overWsDto.CliChattingNewMessageBody
		err = json.Unmarshal(bytes, &reqDto)
		if err != nil {
			return err
		}

		return handler.ChattingNewMessage(conn, reqDto)
	} else if pack.Operation == overWsDto.CHOOSING_USERS_CHOSEN {
		var reqDto overWsDto.CliChoosingUsersChosenBody
		err = json.Unmarshal(bytes, &reqDto)
		if err != nil {
			return nil
		}

		return handler.ChoosingUsersChosen(conn, reqDto)
	}

	// ***

	return fmt.Errorf("routeWsPack, operation is unknown")
}