package n2

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func Exec() {
	upgrader := websocket.Upgrader{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err.Error())
			return
		}

		// ***

		conn.SetCloseHandler(func(code int, text string) error {
			log.Printf("close with `status` %v and `text` %v", code, text)

			data := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
			conn.WriteMessage(websocket.CloseMessage, data)

			err := conn.Close()
			if err != nil {
				log.Printf("close err: %v", err)
			}

			return nil
		})

		listen(conn)
	})

	http.ListenAndServe("127.0.0.1:45678", nil)
}

func listen(conn *websocket.Conn) {
	for {
		msgType, data, err := conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
			log.Println("Unexpected Close Error")
			return
		}

		if err != nil {
			log.Println("Read Message err:", err.Error())
			return
		}

		log.Println(msgType)
		log.Println(string(data))

		err = conn.WriteMessage(msgType, data)
		if err != nil {
			log.Println("Write Message err:", err.Error())
			return
		}
	}
}
