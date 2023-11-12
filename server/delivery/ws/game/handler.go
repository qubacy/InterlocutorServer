package game

import (
	"fmt"
	"ilserver/delivery/ws/game/connection"
	"ilserver/delivery/ws/game/dto"
	"ilserver/pkg/utility"
	"ilserver/service/game"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler struct {
	gameService       *game.Service
	connectionStorage *Storage
}

func NewHandler(gameService *game.Service) *Handler {
	return &Handler{
		gameService:       gameService,
		connectionStorage: NewStorage(),
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

func (h *Handler) SearchingStart(conn *connection.Connection, pack dto.Pack) {
	fmt.Println("SearchingStart")
}

func (h *Handler) SearchingStop(conn *connection.Connection, pack dto.Pack) {
	fmt.Println("SearchingStop")
}

func (h *Handler) ChattingNewMessage(conn *connection.Connection, pack dto.Pack) {
	fmt.Println("ChattingNewMessage")
}

func (h *Handler) ChoosingUsersChosen(conn *connection.Connection, pack dto.Pack) {
	fmt.Println("ChoosingUsersChosen")
}

// private
// -----------------------------------------------------------------------

func (h *Handler) listen(conn *websocket.Conn) {
	connObj := h.connectionStorage.AddOpenConnection(conn)
	defer h.connectionStorage.RemoveConnection(connObj.Id())

	for {
		select {
		case msg := <-connObj.Reader():
			h.route(connObj, msg)
		case <-connObj.Closed().Done():
			return
		}
	}
}

func (h *Handler) route(conn *connection.Connection, message connection.Message) {
	pack, err := dto.MakePackFromJson(message.Data)
	if err != nil {
		log.Println(utility.CreateCustomError(h.route, err))
		conn.CloseGracefully()
	} else {
		switch pack.Operation {

		case dto.SEARCHING_START:
			h.SearchingStart(conn, pack)
		case dto.SEARCHING_STOP:
			h.SearchingStop(conn, pack)

		case dto.CHATTING_NEW_MESSAGE:
			h.ChattingNewMessage(conn, pack)
		case dto.CHOOSING_USERS_CHOSEN:
			h.ChoosingUsersChosen(conn, pack)

		default:
			log.Println(utility.CreateCustomError(h.route,
				ErrUnknownMessageOperation))
			conn.CloseGracefully()
		}
	}
}
