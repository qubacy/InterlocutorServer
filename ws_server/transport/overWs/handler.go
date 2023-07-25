package overWs

import (
	"encoding/json"
	"fmt"
	"ilserver/domain"
	"ilserver/service"
	"ilserver/transport/overWsDto"
	"time"

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

func (h *CommonHandler) AddConn(conn *websocket.Conn) {
	profileId := uuid.New().String()
	remoteAddr := conn.RemoteAddr().String()

	// ***

	// profile linked to conn pointer!
	h.RemoteAddrAndProfileId[remoteAddr] = profileId
	h.ProfileIdAndConn[profileId] = conn
}

func (h *CommonHandler) RemoveConn(conn *websocket.Conn) {
	available, profileId := h.ProfileIdByConn(conn)
	if !available {
		return
	}

	h.RoomService.RemoveProfileByIdBlocking(profileId)
}

func (h *CommonHandler) RemoveConnAndClose(conn *websocket.Conn) error {
	h.RemoveConn(conn)

	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "reason"))
	return conn.Close()
}

func (h *CommonHandler) ProfileIdByConn(conn *websocket.Conn) (bool, string) {
	remoteAddr := conn.RemoteAddr().String()
	id, ok := h.RemoteAddrAndProfileId[remoteAddr]
	return ok, id
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
	if err := conn.WriteMessage(websocket.TextMessage, errPackBytes); err != nil {
		return err
	}

	return conn.Close()
}

// -----------------------------------------------------------------------

// TODO: что передать в этом метод body или raw?
func (h *CommonHandler) SearchingStart(
	conn *websocket.Conn, reqDto overWsDto.CliSearchingStartBodyClient) error {

	if !reqDto.IsValid() {
		h.Err(conn, SEARCHING_START, "body parameters are invalid")
		return fmt.Errorf("CommonHandler, SearchingStart, req dto is invalid")
	}

	// ***

	// TODO: создать метод в сервисе
	h.RoomService.Mx.Lock()

	available, room := h.RoomService.RoomWithSearchingState()
	if !available {
		h.RoomService.AddRoomWithSearchingState()
		room = &h.RoomService.Rooms[len(h.RoomService.Rooms)-1]
	}

	available, profileId := h.ProfileIdByConn(conn)
	if !available {
		return fmt.Errorf("Profile id is not available")
	}

	profile := overWsDto.MakeProfileFromReqDto(profileId, reqDto)
	room.Profiles = append(room.Profiles, profile)

	// ***

	var searchingState = room.State.(domain.SearchingStateRoom)
	searchingState.LaunchTime = time.Now()
	room.State = searchingState

	h.RoomService.Mx.Unlock()

	// ***

	packBytes := overWsDto.MakePackBytes(SEARCHING_START, overWsDto.SvrSearchingStartBody{})
	return conn.WriteMessage(websocket.TextMessage, packBytes)
}

func (h *CommonHandler) ChattingNewMessage(
	conn *websocket.Conn, reqDto overWsDto.CliChattingNewMessageBody) error {

	if !reqDto.IsValid() {
		h.Err(conn, CHATTING_NEW_MESSAGE, "body parameters are invalid")
		return fmt.Errorf("CommonHandler, ChattingNewMessage, req dto is invalid")
	}

	// ***

	available, profileId := h.ProfileIdByConn(conn)
	if !available {
		return fmt.Errorf("Profile id is not available")
	}

	h.RoomService.Mx.Lock()

	available, room := h.RoomService.RoomWithProfileById(profileId)
	if !available {
		h.RoomService.Mx.Unlock()
		return fmt.Errorf("Profile does not belong to the room")
	}

	var localProfileId int = 0
	for i := range room.Profiles {
		if room.Profiles[i].Id == profileId {
			localProfileId = i
		}
	}

	// ***

	resDto := overWsDto.SvrChattingNewMessageBody{
		Message: overWsDto.SvrMessage{
			SenderId: localProfileId,
			Text:     reqDto.Message.Text,
		},
	}

	packBytes := overWsDto.MakePackBytes(CHATTING_NEW_MESSAGE, resDto)
	for i := range room.Profiles {
		c := h.ProfileIdAndConn[room.Profiles[i].Id]
		if err := c.WriteMessage(websocket.TextMessage, packBytes); err != nil {
			h.RemoveConnAndClose(c)
		}
	}

	h.RoomService.Mx.Unlock()
	return nil
}
