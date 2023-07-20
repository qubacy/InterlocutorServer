package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// TODO: ограничить время бездействия
	}

	// ***

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		websocket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Websocket Connected!")
		listen(websocket)
	})

	// ***

	http.ListenAndServe(":47777", nil)
}

/*

Requests:
- startSearching;
- stopSearching;
- sendMessage;
- leave (during the game);
- chooseUsers;


{"operation":operation_id_int,"body":{"error":{"message": error_text}}}

*/

const (
	SEARCHING_START int = iota
	SEARCHING_STOP
	SEARCHING_GAME_FOUND

	CHATTING_NEW_MESSAGE
	CHATTING_STAGE_IS_OVER

	CHOOSING_USERS_CHOSEN
	CHOOSING_STAGE_IS_OVER

	RESULTS_MATCHED_USERS_GOTTEN
)

// proto
// -----------------------------------------------------------------------

type SearchingStartBody struct {
	Profile struct {
		Username string `json:"username"`
		Contact  string `json:"contact"`
	} `json:"profile"`
}

// utils
// -----------------------------------------------------------------------

type WsPack struct {
	Operation int                    `json:"operation"`
	RawBody   map[string]interface{} `json:"body"`
}

type ChatMessage struct {
}

type ChatRoom struct {
}

type ErrBody struct {
	Err struct {
		Message string `json:"message"`
	} `json:"error"`
}

// -----------------------------------------------------------------------

func listen(conn *websocket.Conn) {
	for {
		messageType, messageContent, err := conn.ReadMessage()
		if messageType != websocket.TextMessage {
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

		var pack WsPack
		err = json.Unmarshal(messageContent, &pack)
		if err != nil {
			// TODO: правильная обработка ошибки синтаксиса пакета
			log.Println(conn.RemoteAddr(), err)
			conn.Close()
			return
		}
		log.Println(conn.RemoteAddr(), pack)

		// *** to handle

		messageResponse, err := routeWsPack(pack)
		if err != nil {
			log.Println(conn.RemoteAddr(), err)
			conn.Close()
			return
		}

		// ***

		if err := conn.WriteMessage(messageType, []byte(messageResponse)); err != nil {
			log.Println(conn.RemoteAddr(), err)
			conn.Close()
			return
		}

	}
}

// -----------------------------------------------------------------------

func routeWsPack(pack WsPack) ([]byte, error) {
	if pack.Operation == SEARCHING_START {
		bytes, _ := json.Marshal(pack.RawBody)
		var ssBody SearchingStartBody
		err := json.Unmarshal(bytes, &ssBody)

		if err != nil {
			return []byte{}, fmt.Errorf("body is invalid")
		}

		log.Println(ssBody)

	} else if pack.Operation == SEARCHING_STOP {

	}

	return []byte(""), nil
}
