package overWs

import (
	"fmt"
	"ilserver/domain"
	"ilserver/service"
	"ilserver/transport/overWsDto"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// TODO: должен отправлять пакеты сам?
type CommonHandler struct {
	Mx                     sync.RWMutex
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

	// profile linked to conn pointer
	h.Mx.Lock()
	h.RemoteAddrAndProfileId[remoteAddr] = profileId
	h.ProfileIdAndConn[profileId] = conn
	h.Mx.Unlock()
}

// private...
func (h *CommonHandler) removeConn(conn *websocket.Conn) {
	available, profileId := h.ProfileIdByConn(conn)
	if available {
		h.Mx.Lock()
		remoteAddr := conn.RemoteAddr().String()
		delete(h.RemoteAddrAndProfileId, remoteAddr)
		h.Mx.Unlock()

		h.RoomService.RemoveProfileByIdBlocking(profileId)
	}
	// may already be removed
}

func (h *CommonHandler) RemoveConnAndClose(conn *websocket.Conn) error {
	h.removeConn(conn)

	// ***

	closeBytes := websocket.FormatCloseMessage(1000, "reason")
	err := conn.WriteMessage(websocket.CloseMessage, closeBytes)
	if err == nil {
		return conn.Close()
	}

	return err
}

// TODO: изменить порядок выходных параметров (?)
func (h *CommonHandler) ProfileIdByConn(conn *websocket.Conn) (bool, string) {
	h.Mx.RLock()
	remoteAddr := conn.RemoteAddr().String()
	profileId, avb := h.RemoteAddrAndProfileId[remoteAddr]
	h.Mx.RUnlock()

	return avb, profileId
}

// -----------------------------------------------------------------------

const (
	OperationIsUnknown string = "operation is unknown"
)

func (h *CommonHandler) Err(conn *websocket.Conn, opCode int, errText string) error {
	srvErr := overWsDto.SvrErrBody{
		Err: overWsDto.Err{
			Message: errText,
		},
	}

	// ***

	errPackBytes := overWsDto.MakePackBytes(opCode, srvErr)
	if err := conn.WriteMessage(websocket.TextMessage, errPackBytes); err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------

// TODO: что передать в этом метод body или raw?
func (h *CommonHandler) SearchingStart(
	conn *websocket.Conn, reqDto overWsDto.CliSearchingStartBodyClient) error {

	if !reqDto.IsValid() {
		h.Err(conn, overWsDto.SEARCHING_START, "body parameters are invalid")
		return fmt.Errorf("CommonHandler, SearchingStart, req dto is invalid")
	}

	// ***

	// TODO: создать метод в сервисе
	h.RoomService.Mx.Lock()

	available, room := h.RoomService.AvailableRoomWithSearchingState()
	if !available {
		h.RoomService.AddRoomWithSearchingState()
		room = &h.RoomService.Rooms[len(h.RoomService.Rooms)-1]
	}

	available, profileId := h.ProfileIdByConn(conn)
	if !available {
		h.RoomService.Mx.Unlock()
		return fmt.Errorf("Profile id is not available")
	}

	profile := overWsDto.MakeProfileFromReqDto(profileId, reqDto)
	room.Profiles = append(room.Profiles, profile)

	// ***

	var searchingState = room.State.(*domain.SearchingStateRoom)
	searchingState.LaunchTime = time.Now()
	room.State = searchingState

	h.RoomService.Mx.Unlock()

	// ***

	packBytes := overWsDto.MakePackBytes(overWsDto.SEARCHING_START, overWsDto.SvrSearchingStartBody{})
	return conn.WriteMessage(websocket.TextMessage, packBytes)
}

func (h *CommonHandler) SearchingStop(
	conn *websocket.Conn, reqDto overWsDto.CliSearchingStopBody) error {

	// TODO: будет тут реализация?

	return nil
}

func (h *CommonHandler) ChattingNewMessage(
	conn *websocket.Conn, reqDto overWsDto.CliChattingNewMessageBody) error {

	if !reqDto.IsValid() {
		h.Err(conn, overWsDto.CHATTING_NEW_MESSAGE, "body parameters are invalid")
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

	packBytes := overWsDto.MakePackBytes(overWsDto.CHATTING_NEW_MESSAGE, resDto)
	for i := range room.Profiles {
		c := h.ProfileIdAndConn[room.Profiles[i].Id]
		if err := c.WriteMessage(websocket.TextMessage, packBytes); err != nil {
			h.RemoveConnAndClose(c)
		}
	}

	h.RoomService.Mx.Unlock()
	return nil
}

func (h *CommonHandler) ChoosingUsersChosen(
	conn *websocket.Conn, reqDto overWsDto.CliChoosingUsersChosenBody) error {

	if !reqDto.IsValid() {
		h.Err(conn, overWsDto.CHOOSING_USERS_CHOSEN, "body parameters are invalid")
		return fmt.Errorf("CommonHandler, ChoosingUsersChosen, req dto is invalid")
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

	localProfileId := findLocalIdByProfileId(&room.Profiles, profileId)
	if localProfileId == -1 {
		h.RoomService.Mx.Unlock()
		return fmt.Errorf("Profile was not found")
	}

	// ***

	switch room.State.(type) {
	case *domain.ChoosingStateRoom:
		choosingState := room.State.(*domain.ChoosingStateRoom)
		matchedIds := []string{}
		for _, userId := range reqDto.UserIdList {
			matchedIds = append(matchedIds, room.Profiles[userId].Id)
		}

		choosingState.ProfileIdAndMatchedIds[profileId] = matchedIds
	default:
		h.RoomService.Mx.Unlock()
		return fmt.Errorf("Room state is not choosing")
	}

	h.RoomService.Mx.Unlock()
	return nil
}

// utils
// -----------------------------------------------------------------------

func findLocalIdByProfileId(profiles *[]domain.Profile, profileId string) int {
	for i := range *profiles {
		if (*profiles)[i].Id == profileId {
			return i
		}
	}

	return -1 // err
}
