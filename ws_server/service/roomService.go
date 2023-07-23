package service

import (
	"encoding/json"
	"ilserver/domain"

	"ilserver/transport/dto"
	"ilserver/transport/overWs"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RoomService struct {
	Mx    sync.RWMutex
	Rooms []domain.Room
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

// -----------------------------------------------------------------------

// TODO: нужен ли контекст?
func (rs *RoomService) BackgroundWork(interval time.Duration) {
	go func() {
		for {
			select {
			case <-time.After(interval):
				rs.backgroundWorkIteration()
			}
		}
	}()
}

// TODO: внешняя сущность принимающая Handler
func (rs *RoomService) backgroundWorkIteration() {
	log.Println("RoomService, backgroundWorkIteration...")

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
		"кино", "музыка", "настольные игры", "уличные алкоголики", "vape nation",
	}
	return topics[rand.Intn(len(topics))]
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

	var gameFoundBody dto.SvrSearchingGameFoundBody
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

		c := rs.Rooms[roomInx].Profiles[i].Conn
		gameFoundBodyBytes, _ := json.Marshal(gameFoundBody)
		rawBody := make(map[string]interface{})
		json.Unmarshal(gameFoundBodyBytes, &rawBody)
		pack := overWs.Pack{
			Operation: 2,
			RawBody:   rawBody,
		}

		gameFoundPackBytes, _ := json.Marshal(pack)
		c.WriteMessage(websocket.TextMessage, gameFoundPackBytes)
	}
}

func (rs *RoomService) updateRoomWithChattingState(roomInx int) {

}

func (rs *RoomService) updateRoomWithChoosingState(roomInx int) {

}
