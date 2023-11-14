package game

import (
	"fmt"
	"ilserver/delivery/ws/game/connection"
	"ilserver/delivery/ws/game/dto"
	"ilserver/pkg/utility"
	"ilserver/service/game"
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

// can also pass key-value storage!
func (h *Handler) SearchingStart(conn *connection.Connection, pack dto.Pack) {
	cliBody, err := pack.AsCliSearchingStartBody()
	if err != nil {
		closeGracefully(conn, h.SearchingStart, err)
		return
	}

	svrBody, err := h.gameService.SearchingStart(conn.Id(), cliBody)
	if err != nil {
		// TODO: convert to sending error code!?
		closeGracefully(conn, h.SearchingStart, err)
		return
	}

	jsonBytes, err := dto.MakePackAsJsonBytes(dto.SEARCHING_START, svrBody)
	if err != nil {
		closeGracefully(conn, h.SearchingStart, err)
		return
	}

	// ***

	conn.Writer() <- connection.MakeTextMessage(jsonBytes)
}

func (h *Handler) SearchingStop(conn *connection.Connection, pack dto.Pack) {

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
		closeGracefully(conn, h.route, err)
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
			closeGracefully(conn, h.route,
				ErrUnknownMessageOperation)
		}
	}
}

// hidden functions
// -----------------------------------------------------------------------

func closeGracefully(
	conn *connection.Connection,
	i interface{}, err error,
) {
	err = utility.CreateCustomError(i, err)
	err = utility.UnwrapErrorsToLast(err) // since the control message is limited!
	conn.CloseGracefully(err.Error())
}
