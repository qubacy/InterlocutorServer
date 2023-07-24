package main

import (
	"encoding/json"
	"ilserver/config"
	"ilserver/transport/overWs"
	"ilserver/transport/overWsDto"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatal("config initialization has been failed with err:", err)
	}

	// ***

	handler := overWs.NewCommonHandler()
	overWs.BackgroundUpdateRooms(handler)

	// ***

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// TODO: ограничить время бездействия
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		websocket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Websocket Connected!")
		listen(handler, websocket) // app, over ws
	})

	// ***

	http.HandleFunc("/debug/rooms", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := json.Marshal(handler.RoomService.Rooms)
		w.Header().Add("Content-Type", "application/json")
		w.Write(bytes)
	})

	// ***

	http.ListenAndServe(":"+viper.GetString("port"), nil)
}

// -----------------------------------------------------------------------

func listen(handler *overWs.CommonHandler, conn *websocket.Conn) {
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

		var pack overWsDto.Pack
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
