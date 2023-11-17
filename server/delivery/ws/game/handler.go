package game

import (
	"fmt"
	"ilserver/delivery/ws/game/connection"
	"ilserver/delivery/ws/game/dto"
	"ilserver/pkg/utility"
	"ilserver/service/game"
	serviceDto "ilserver/service/game/dto"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler struct {
	gameService       *game.Service
	connectionStorage *Storage
}

func NewHandler(gameService *game.Service) *Handler {
	instance := &Handler{
		gameService:       gameService,
		connectionStorage: NewStorage(),
	}

	go instance.processAsyncResponses()
	go instance.processAsyncErrors()

	return instance
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
		h.closeGracefullyWithError(conn, h.SearchingStart, err)
		return
	}

	svrBody, err := h.gameService.SearchingStart(conn.Id(), cliBody)
	if err != nil {
		// TODO: convert to sending error code!?
		h.closeGracefullyWithError(conn, h.SearchingStart, err)
		return
	}

	jsonBytes, err := dto.MakePackAsJsonBytes(dto.SEARCHING_START, svrBody)
	if err != nil {
		h.closeGracefullyWithError(conn, h.SearchingStart, err)
		return
	}

	// ***

	conn.Writer() <- connection.MakeTextMessage(jsonBytes)
}

func (h *Handler) SearchingStop(conn *connection.Connection, pack dto.Pack) {
	cliBody, err := pack.AsCliSearchingStopBody()
	if err != nil {
		h.closeGracefullyWithError(conn, h.SearchingStop, err)
		return
	}

	err = h.gameService.SearchingStop(conn.Id(), cliBody)
	if err != nil {
		h.closeGracefullyWithError(conn, h.SearchingStop, err)
		return
	}

	h.closeGracefully(conn)
}

func (h *Handler) ChattingNewMessage(conn *connection.Connection, pack dto.Pack) {
	cliBody, err := pack.AsCliChattingNewMessageBody()
	if err != nil {
		h.closeGracefullyWithError(conn, h.ChattingNewMessage, err)
		return
	}

	svrBody, err := h.gameService.ChattingNewMessage(conn.Id(), cliBody)
	if err != nil {
		h.closeGracefullyWithError(conn, h.ChattingNewMessage, err)
		return
	}

	jsonBytes, err := dto.MakePackAsJsonBytes(dto.CHATTING_NEW_MESSAGE, svrBody)
	if err != nil {
		h.closeGracefullyWithError(conn, h.ChattingNewMessage, err)
		return
	}

	// ***

	// message to myself...
	conn.Writer() <- connection.MakeTextMessage(jsonBytes)
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
			h.gameService.ProfileLeftGame(connObj.Id())
			return
		}
	}
}

func (h *Handler) route(conn *connection.Connection, message connection.Message) {
	pack, err := dto.MakePackFromJson(message.Data)
	if err != nil {
		h.closeGracefullyWithError(conn, h.route, err)
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
			h.closeGracefullyWithError(conn, h.route,
				ErrUnknownMessageOperation)
		}
	}
}

// -----------------------------------------------------------------------

func (h *Handler) processAsyncResponses() {
	for {
		select {
		case response := <-h.gameService.AsyncResponse():
			conn, exist := h.connectionStorage.GetConnection(response.ProfileId)
			if !exist {
				log.Println("connection does not exist")
				continue
			}

			operation, exist := operationByServerBody(response.ServerBody)
			if !exist {
				log.Println("unknown operation")
				continue
			}

			jsonBytes, err := dto.MakePackAsJsonBytes(
				operation, response.ServerBody)

			if err != nil {
				h.closeGracefullyWithError(conn, h.processAsyncResponses, err)
				continue
			}

			conn.Writer() <- connection.MakeTextMessage(jsonBytes)
		}
	}
}

func (h *Handler) processAsyncErrors() {
	for {
		select {
		case response := <-h.gameService.AsyncResponseAboutError():
			conn, exist := h.connectionStorage.GetConnection(response.ProfileId)
			if exist {
				h.closeGracefullyWithError(conn,
					h.processAsyncErrors, response.Err)
			} else {
				log.Println("connection does not exist")
			}
		}
	}
}

// hidden functions
// -----------------------------------------------------------------------

func (h *Handler) closeGracefully(conn *connection.Connection) {
	conn.CloseGracefully("")
}

func (h *Handler) closeGracefullyWithError(
	conn *connection.Connection,
	i interface{}, err error,
) {
	err = utility.CreateCustomError(i, err)
	err = utility.UnwrapErrorsToLast(err) // since the control message is limited!
	conn.CloseGracefully(err.Error())
}

// -----------------------------------------------------------------------

// bad solution?
func operationByServerBody(i interface{}) (int, bool) {
	switch i.(type) {
	case serviceDto.SvrSearchingStartBody:
		return dto.SEARCHING_START, true
	case serviceDto.SvrSearchingGameFoundBody:
		return dto.SEARCHING_GAME_FOUND, true

	case serviceDto.SvrChattingNewMessageBody:
		return dto.CHATTING_NEW_MESSAGE, true
	case serviceDto.SvrChattingStageIsOverBody:
		return dto.CHATTING_STAGE_IS_OVER, true

	case serviceDto.SvrChoosingUsersChosenBody:
		return dto.CHOOSING_USERS_CHOSEN, true
	case serviceDto.SvrChoosingStageIsOverBody:
		return dto.CHOOSING_STAGE_IS_OVER, true
	}
	return 0, false
}
