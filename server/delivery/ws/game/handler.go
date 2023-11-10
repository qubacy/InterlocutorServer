package game

import (
	"encoding/json"
	"ilserver/delivery/ws/game/dto"
	"ilserver/service/game"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Handler struct {
	gameService       *game.Service
	connectionStorage *Storage
}

func NewHandler(gameService *game.Service) *Handler {
	return &Handler{
		gameService: gameService,
	}
}

func (h *Handler) Mux(pathStart string) *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", h.react)
	return serveMux
}

// react and listen
// -----------------------------------------------------------------------

func (h *Handler) react(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.listen(websocket)
}

// -----------------------------------------------------------------------

func (h *Handler) SearchingStart(pack dto.Pack) {

}

func (h *Handler) SearchingStop() {

}

func (h *Handler) ChattingNewMessage() {

}

func (h *Handler) ChoosingUsersChosen() {

}

// private
// -----------------------------------------------------------------------

func (h *Handler) listen(conn *websocket.Conn) {
	handler.AddConn(conn)

	connectionId := uuid.NewString()

	for {
		messageType, messageContent, err := conn.ReadMessage()
		if err != nil {

		}

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
