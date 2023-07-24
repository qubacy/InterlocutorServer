package overWs

import (
	"encoding/json"
	"fmt"
	"ilserver/service"
	"ilserver/transport/overWsDto"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// TODO: должен отправлять пакеты сам?
type CommonHandler struct {
	RemoteAddrAndProfileId map[string]string
	ProfileIdAndConn       map[string]*websocket.Conn

	RoomService *service.RoomService
}

func NewCommonHandler() *CommonHandler {
	return &CommonHandler{
		ProfileIdAndConn:       make(map[string]*websocket.Conn),
		RemoteAddrAndProfileId: make(map[string]string),
		RoomService:            service.NewRoomService(),
	}
}

// -----------------------------------------------------------------------

func (h *CommonHandler) AddConn(conn *websocket.Conn) error {
	profileId := uuid.New().String()
	remoteAddr := conn.RemoteAddr().String()

	// ***

	h.RemoteAddrAndProfileId[remoteAddr] = profileId
	h.ProfileIdAndConn[profileId] = conn

	// TODO: добавить проверки на уникальность (?)

	h.ProfileIdAndConn[profileId] = conn
	return nil
}

func (h *CommonHandler) ProfileIdByConn(conn *websocket.Conn) string {
	remoteAddr := conn.RemoteAddr().String()
	return h.RemoteAddrAndProfileId[remoteAddr]
}

// -----------------------------------------------------------------------

func (h *CommonHandler) Err(conn *websocket.Conn, opCode int, errText string) error {
	srvErr := overWsDto.SvrErrBody{}
	srvErr.Err.Message = errText
	srvErrBytes, _ := json.Marshal(srvErr)

	// ***

	rawBody := make(map[string]interface{})
	json.Unmarshal(srvErrBytes, &rawBody)
	errPack := overWsDto.Pack{
		Operation: opCode,
		RawBody:   rawBody,
	}

	errPackBytes, _ := json.Marshal(errPack)
	conn.WriteMessage(websocket.TextMessage, errPackBytes)

	return conn.Close()
}

// -----------------------------------------------------------------------

// TODO: что передать в этом метод body или raw?
func (h *CommonHandler) SearchingStart(
	conn *websocket.Conn, reqDto overWsDto.CliSearchingStartBodyClient) error {

	if !reqDto.IsValid() {
		h.Err(conn, SEARCHING_START, "body parameters are invalid")
		return fmt.Errorf("SearchingStart, req dto is invalid")
	}

	// ***

	h.RoomService.Mx.Lock()
	available, room := h.RoomService.RoomWithSearchingState()
	if !available {
		h.RoomService.AddRoomWithSearchingState()
		room = &h.RoomService.Rooms[len(h.RoomService.Rooms)-1]
	}

	profile := overWsDto.MakeProfileFromReqDto(h.ProfileIdByConn(conn), reqDto)
	room.Profiles = append(room.Profiles, profile)
	h.RoomService.Mx.Unlock()

	// TODO: обновить время LastUpdateTime

	// ***

	pack := overWsDto.MakePack(SEARCHING_START, overWsDto.SvrSearchingStartBody{})
	packBytes, _ := json.Marshal(pack)
	conn.WriteMessage(websocket.TextMessage, packBytes)

	return nil
}
