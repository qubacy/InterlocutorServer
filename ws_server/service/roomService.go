package service

import (
	"encoding/json"
	"ilserver/domain"
	"ilserver/transport/overWsDto"
	"math/rand"
	"sync"
	"time"
)

type UpdateRoomMessage struct {
	ProfileId   string
	BytesResDto []byte
}

type RoomService struct {
	Mx sync.RWMutex

	// TODO: преобразовать в список
	Rooms []domain.Room

	// ***

	UpdateRoomMsgs chan UpdateRoomMessage
}

func NewRoomService() *RoomService {
	return &RoomService{
		Mx:             sync.RWMutex{},
		Rooms:          make([]domain.Room, 0),
		UpdateRoomMsgs: make(chan UpdateRoomMessage),
	}
}

// -----------------------------------------------------------------------

func (rs *RoomService) RemoveProfileByIdBlocking(profileId string) {
	rs.Mx.Lock()
	for i := range rs.Rooms {
		for j := range rs.Rooms[i].Profiles {
			if rs.Rooms[i].Profiles[j].Id == profileId {

				rs.Rooms[i].Profiles =
					append(rs.Rooms[i].Profiles[:j],
						rs.Rooms[i].Profiles[j+1:]...)
			}
		}
	}
	rs.Mx.Unlock()
}

func (rs *RoomService) RoomWithProfileById(profileId string) (bool, *domain.Room) {
	for i := range rs.Rooms {
		for j := range rs.Rooms[i].Profiles {
			if rs.Rooms[i].Profiles[j].Id == profileId {
				return true, &rs.Rooms[i]
			}
		}
	}
	return false, nil
}

func (rs *RoomService) RoomWithSearchingState() (bool, *domain.Room) {
	for i := range rs.Rooms {
		switch rs.Rooms[i].State.(type) {
		case domain.SearchingStateRoom:
			return true, &rs.Rooms[i]
		}
	}
	return false, nil
}

func (rs *RoomService) AddRoomWithSearchingState() {
	currentTime := time.Now()
	one := domain.Room{
		State: domain.SearchingStateRoom{
			RoomState: domain.RoomState{
				Name:       domain.SEARCHING,
				LaunchTime: currentTime,
			},
			LastUpdateTime: currentTime,
		},
	}
	rs.Rooms = append(rs.Rooms, one)
}

// -----------------------------------------------------------------------

func (rs *RoomService) BackgroundUpdateRoomsTick() {
	//log.Println("RoomService, BackgroundUpdateRoomsTick...")

	// ***

	rs.Mx.Lock()
	for i := range rs.Rooms {
		switch rs.Rooms[i].State.(type) {
		case domain.SearchingStateRoom:
			rs.updateRoomWithSearchingState(i)
		case domain.ChattingStateRoom:
			rs.updateRoomWithChattingState(i)
		case domain.ChoosingStateRoom:
			rs.updateRoomWithChoosingState(i)
		}
	}
	rs.Mx.Unlock()
}

func randChattingTopic() string {
	topics := []string{
		"кино", "музыка",
		"настольные игры",
		//...
	}

	index := rand.Intn(len(topics))
	return topics[index]
}

// TODO: параметры из конфига, попробовать viper
func (rs *RoomService) updateRoomWithSearchingState(roomInx int) {
	if len(rs.Rooms[roomInx].Profiles) <= 1 {
		return
	}

	searchingState := rs.Rooms[roomInx].State.(domain.SearchingStateRoom)
	dur := time.Now().Sub(searchingState.LastUpdateTime)

	if dur < 10*time.Second { // TODO: в конфиг
		return
	}

	rs.Rooms[roomInx].State = domain.ChattingStateRoom{
		RoomState: domain.RoomState{
			Name: domain.CHATTING,
		},
	}

	// ***

	startSessionTimeUnix := time.Now().Unix()

	var gameFoundBody overWsDto.SvrSearchingGameFoundBody
	gameFoundBody.FoundGameData.StartSessionTime = startSessionTimeUnix
	gameFoundBody.FoundGameData.ChattingStageDuration = 5 * 60 * 1000 // TODO: в конфиг
	gameFoundBody.FoundGameData.ChoosingStageDuration = 30 * 1000     // TODO: в конфиг
	gameFoundBody.FoundGameData.ChattingTopic = randChattingTopic()   // TODO: из БД получить на английском, потом перевод

	// TODO: объявить полноценную структуру
	for i := range rs.Rooms[roomInx].Profiles {
		gameFoundBody.FoundGameData.ProfilePublicList = append(
			gameFoundBody.FoundGameData.ProfilePublicList,
			struct {
				Id       int    "json:\"id\""
				Username string "json:\"username\""
			}{
				Id:       i,
				Username: rs.Rooms[roomInx].Profiles[i].Username,
			})
	}

	// ***

	for i := range rs.Rooms[roomInx].Profiles {
		gameFoundBody.FoundGameData.LocalProfileId = i

		// ***

		gameFoundBodyBytes, _ := json.Marshal(gameFoundBody)
		rawBody := make(map[string]interface{})
		json.Unmarshal(gameFoundBodyBytes, &rawBody)
		pack := overWsDto.Pack{
			Operation: 2,
			RawBody:   rawBody,
		}

		gameFoundPackBytes, _ := json.Marshal(pack)
		msg := UpdateRoomMessage{
			ProfileId:   rs.Rooms[roomInx].Profiles[i].Id,
			BytesResDto: gameFoundPackBytes,
		}

		rs.UpdateRoomMsgs <- msg
	}
}

func (rs *RoomService) updateRoomWithChattingState(roomInx int) {

}

func (rs *RoomService) updateRoomWithChoosingState(roomInx int) {

}
