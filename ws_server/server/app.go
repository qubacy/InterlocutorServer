package server

import (
	"encoding/json"
	"ilserver/transport/overWs"
	"ilserver/transport/overWsDto"
	"log"
	"net/http"
	"sync"
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
		if err := listen(handler, websocket); err != nil {
			log.Println("One ws listen err:", err)
			websocket.Close()
			return
		}
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

func listen(handler *overWs.CommonHandler, conn *websocket.Conn) error {
	var once sync.Once
	for {
		messageType, messageContent, err := conn.ReadMessage()

		if err != nil {
			if err.(*websocket.CloseError) == nil {
				return handler.RemoveConnAndClose(conn)
			}

			return err
		}

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

		once.Do(func() {
			handler.AddConn(conn)
		})

		err = routeWsPack(handler, conn, pack)
		if err != nil {
			log.Println(conn.RemoteAddr(), err)
			return handler.RemoveConnAndClose(conn)
		}
	}
}

func routeWsPack(handler *overWs.CommonHandler, conn *websocket.Conn, pack overWsDto.Pack) error {

	// TODO: обработку перенести в handler, DTOs передавать сервисам

	bytes, err := json.Marshal(pack.RawBody)
	if err != nil {
		return err
	}

	// ***

	if pack.Operation == overWs.SEARCHING_START {
		var reqDto overWsDto.CliSearchingStartBodyClient
		err = json.Unmarshal(bytes, &reqDto)
		if err != nil {
			return err
		}

		return handler.SearchingStart(conn, reqDto)
	} else if pack.Operation == overWs.SEARCHING_STOP {
		var reqDto overWsDto.CliSearchingStopBody
		err = json.Unmarshal(bytes, &reqDto)
		if err != nil {
			return err
		}

		return handler.SearchingStop(conn, reqDto)
	} else if pack.Operation == overWs.CHATTING_NEW_MESSAGE {
		var reqDto overWsDto.CliChattingNewMessageBody
		err = json.Unmarshal(bytes, &reqDto)
		if err != nil {
			return err
		}

		return handler.ChattingNewMessage(conn, reqDto)
	} else if pack.Operation == overWs.CHOOSING_USERS_CHOSEN {
		var reqDto overWsDto.CliChoosingUsersChosenBody
		err = json.Unmarshal(bytes, &reqDto)
		if err != nil {
			return nil
		}

		return handler.ChoosingUsersChosen(conn, reqDto)
	}

	return handler.Err(conn, pack.Operation,
		overWs.OperationIsUnknown)
}
