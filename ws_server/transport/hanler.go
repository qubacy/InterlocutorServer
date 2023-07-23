package transport

import (
	"encoding/json"
	"errors"
	"ilserver/service"
	"ilserver/transport/dto"
	"ilserver/transport/overWs"
	"unsafe"

	"github.com/gorilla/websocket"
)

// TODO: должен отправлять пакеты сам?
type Handler struct {
	Conns       map[uint64]*websocket.Conn
	RoomService service.RoomService
}

// -----------------------------------------------------------------------

func (h *Handler) AddConn(conn *websocket.Conn) error {
	ptr := unsafe.Pointer(conn)
	idr := *(*uint64)(ptr)

	if _, ok := h.Conns[idr]; ok {
		return errors.New("Handler, AddConn, key is already in use")
	}

	h.Conns[idr] = conn
	return nil
}

func (h *Handler) ConnIdr(conn *websocket.Conn) uint64 {
	ptr := unsafe.Pointer(conn)
	idr := *(*uint64)(ptr)
	return idr
}

// -----------------------------------------------------------------------

func (h *Handler) Err(conn *websocket.Conn, opCode int, errText string) error {
	srvErr := dto.SrvErrBody{}
	srvErr.Err.Message = errText
	srvErrBytes, _ := json.Marshal(srvErr)

	// ***

	rawBody := make(map[string]interface{})
	json.Unmarshal(srvErrBytes, &rawBody)
	errPack := overWs.Pack{
		Operation: opCode,
		RawBody:   rawBody,
	}

	errPackBytes, _ := json.Marshal(errPack)
	conn.WriteMessage(websocket.TextMessage, errPackBytes)

	return conn.Close()
}

// -----------------------------------------------------------------------

// TODO: что передать в этом метод body или raw?
func (h *Handler) SearchingStart(conn *websocket.Conn, body dto.CliSearchingStartBodyClient) error {
	h.RoomService.Mx.Lock()
	hasRoom, roomPtr := h.RoomService.RoomWithSearchingState()

	h.RoomService.Mx.Unlock()

	// ***

	return nil
}
