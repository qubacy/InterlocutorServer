package server

import (
	"encoding/json"
	"ilserver/transport/overWs"
	"ilserver/transport/overWsDto"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

type App struct {
	httpServer *http.Server
}

func NewApp() *App {
	return &App{}
}

// -----------------------------------------------------------------------

func (s *App) Run() error {
	handler := overWs.NewCommonHandler()
	overWs.BackgroundUpdateRooms(handler)

	// *** ws

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		websocket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade err:", err)
			return
		}

		log.Println("Ws", websocket.RemoteAddr().String(), "connected")
		listen(handler, websocket)
	})

	// *** simple debug http

	mux.HandleFunc("/debug/rooms", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := json.Marshal(handler.RoomService.Rooms)
		w.Header().Add("Content-Type", "application/json")
		w.Write(bytes)
	})

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

func listen(handler *overWs.CommonHandler, conn *websocket.Conn) {
	for {
		messageType, messageContent, err := conn.ReadMessage()
		if messageType != websocket.TextMessage {
			log.Println(conn.RemoteAddr(), "message type is not text")
			conn.Close()
			return
		}

		if err != nil {
			log.Println(conn.RemoteAddr(), err)
			conn.Close()
			return
		}

		// ***

		var pack overWsDto.Pack
		err = json.Unmarshal(messageContent, &pack)
		if err != nil {
			log.Println(conn.RemoteAddr(), err)
			// TODO: отправить пакет с информацией об ошибки
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

func routeWsPack(handler *overWs.CommonHandler, conn *websocket.Conn, pack overWsDto.Pack) error {
	if pack.Operation == overWs.SEARCHING_START {
		bytes, err := json.Marshal(pack.RawBody)
		if err != nil {
			return err
		}

		var ssBody overWsDto.CliSearchingStartBodyClient
		err = json.Unmarshal(bytes, &ssBody)
		if err != nil {
			return err
		}

		handler.SearchingStart(conn, ssBody)
		return nil
	} else if pack.Operation == overWs.SEARCHING_STOP {

	} else if pack.Operation == overWs.CHATTING_NEW_MESSAGE {

	} else if pack.Operation == overWs.CHOOSING_USERS_CHOSEN {

	}

	return handler.Err(conn, pack.Operation, "operation is unknown")
}
