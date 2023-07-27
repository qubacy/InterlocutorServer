package service

import (
	"ilserver/domain"
	"ilserver/transport/overWsDto"
	"math/rand"
	"sync"
	"time"

	"github.com/spf13/viper"
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
	isRemoved := false
	for i := range rs.Rooms {
		if isRemoved {
			break
		}

		for j := range rs.Rooms[i].Profiles {
			if rs.Rooms[i].Profiles[j].Id == profileId {

				rs.Rooms[i].Profiles =
					append(rs.Rooms[i].Profiles[:j],
						rs.Rooms[i].Profiles[j+1:]...)
				break
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
		case *domain.SearchingStateRoom:
			return true, &rs.Rooms[i]
		}
	}
	return false, nil
}

func (rs *RoomService) AddRoomWithSearchingState() {
	currentTime := time.Now()
	one := domain.Room{
		State: &domain.SearchingStateRoom{
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
		case *domain.SearchingStateRoom:
			rs.updateRoomWithSearchingState(i)
		case *domain.ChattingStateRoom:
			rs.updateRoomWithChattingState(i)
		case *domain.ChoosingStateRoom:
			rs.updateRoomWithChoosingState(i)
		case nil:
			// TODO: удалить комнату, возможно создать массив индексов
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

func (rs *RoomService) updateRoomWithSearchingState(roomInx int) {
	if len(rs.Rooms[roomInx].Profiles) <= 1 {
		return
	}

	searchingState := rs.Rooms[roomInx].State.(*domain.SearchingStateRoom)
	dur := time.Now().Sub(searchingState.LastUpdateTime)

	// TODO: сложный путь до параметра в файле конфигурации
	interval := viper.GetDuration(
		"background" +
			".update_rooms" +
			".with_searching_state" +
			".interval_from_last_update_to_next_state",
	)
	if dur < interval {
		return
	}

	// TODO: make или new функция
	rs.Rooms[roomInx].State = &domain.ChattingStateRoom{
		RoomState: domain.RoomState{
			Name:       domain.CHATTING,
			LaunchTime: time.Now(),
		},
	}

	// ***

	gameFoundBody := overWsDto.SvrSearchingGameFoundBody{
		FoundGameData: makeFoundGameData(),
	}

	for i := range rs.Rooms[roomInx].Profiles {
		current := &rs.Rooms[roomInx].Profiles[i]
		gameFoundBody.FoundGameData.ProfilePublicList = append(
			gameFoundBody.FoundGameData.ProfilePublicList,
			overWsDto.ProfilePublic{
				Id:       i,
				Username: current.Username,
			})
	}

	// ***

	for i := range rs.Rooms[roomInx].Profiles {
		current := &rs.Rooms[roomInx].Profiles[i]
		gameFoundBody.FoundGameData.LocalProfileId = i

		// ***

		gameFoundPackBytes := overWsDto.MakePackBytes(
			overWsDto.SEARCHING_GAME_FOUND, gameFoundBody)

		msg := UpdateRoomMessage{
			ProfileId:   current.Id,
			BytesResDto: gameFoundPackBytes,
		}

		rs.UpdateRoomMsgs <- msg
	}
}

func (rs *RoomService) updateRoomWithChattingState(roomInx int) {
	// TODO: удалить пустую комнату

	chattingState := rs.Rooms[roomInx].State.(*domain.ChattingStateRoom)
	launchTimeUnix := chattingState.LaunchTime.Unix() * 1000 // ms
	currentTimeUnix := time.Now().Unix() * 1000

	csDuration := int64(viper.GetDuration("found_game.chatting_stage_duration").Milliseconds())
	difference := (currentTimeUnix - launchTimeUnix) // ms

	if difference < csDuration {
		return
	}

	// TODO: make или new функция
	rs.Rooms[roomInx].State = &domain.ChoosingStateRoom{
		RoomState: domain.RoomState{
			Name:       domain.CHOOSING,
			LaunchTime: time.Now(),
		},
	}

	// ***

	for i := range rs.Rooms[roomInx].Profiles {
		current := &rs.Rooms[roomInx].Profiles[i]
		packBytes := overWsDto.MakePackBytes(
			overWsDto.CHATTING_STAGE_IS_OVER,
			overWsDto.SvrChattingStageIsOverBody{})

		msg := UpdateRoomMessage{
			ProfileId:   current.Id,
			BytesResDto: packBytes,
		}

		rs.UpdateRoomMsgs <- msg
	}
}

func (rs *RoomService) updateRoomWithChoosingState(roomInx int) {
	// TODO: удалить пустую комнату

	choosingState := rs.Rooms[roomInx].State.(*domain.ChoosingStateRoom)
	launchTimeUnix := choosingState.LaunchTime.Unix() * 1000 // ms
	currentTimeUnix := time.Now().Unix() * 1000

	csDuration := int64(viper.GetDuration("found_game.choosing_stage_duration").Milliseconds())
	difference := (currentTimeUnix - launchTimeUnix) // ms

	if difference < csDuration {
		return
	}

	// ***

	// TODO:
	rs.Rooms[roomInx].State = nil

	// ***

	for i := range rs.Rooms[roomInx].Profiles {
		current := &rs.Rooms[roomInx].Profiles[i]
		matchedUsers := []overWsDto.MatchedUser{}

		matchedIds := choosingState.ProfileIdAndMatchedIds[current.Id]
		for _, mid := range matchedIds {
			matchedUsers = append(matchedUsers, overWsDto.MatchedUser{
				Id: findLocalIdByProfileId(&rs.Rooms[roomInx].Profiles, mid),
			})
		}

		// ***

		// TODO: сделать проверку взаимного выбора

		// ***

		packBytes := overWsDto.MakePackBytes(
			overWsDto.CHOOSING_STAGE_IS_OVER,
			overWsDto.SvrChoosingStageIsOverBody{
				MatchedUsers: matchedUsers,
			})

		msg := UpdateRoomMessage{
			ProfileId:   current.Id,
			BytesResDto: packBytes,
		}

		rs.UpdateRoomMsgs <- msg
	}
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

func makeFoundGameData() overWsDto.FoundGameData {
	startSessionTimeUnix := time.Now().Unix()
	invalidLocalProfileId := -1

	return overWsDto.FoundGameData{
		LocalProfileId: invalidLocalProfileId,

		StartSessionTime:      startSessionTimeUnix,
		ChattingStageDuration: int64(viper.GetDuration("found_game.chatting_stage_duration").Milliseconds()),
		ChoosingStageDuration: int64(viper.GetDuration("found_game.choosing_stage_duration").Milliseconds()),

		ChattingTopic:     randChattingTopic(),
		ProfilePublicList: []overWsDto.ProfilePublic{},
	}
}
