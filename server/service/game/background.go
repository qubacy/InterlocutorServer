package game

import (
	"context"
	domain "ilserver/domain/memory"
	"ilserver/service/game/dto"
	"time"
)

func (s *Service) backgroundUpdates(ctx context.Context) {
	timeout := s.config.TimeoutForUpdateRooms
	for {
		select {
		case <-time.After(timeout):
			s.updateRooms()
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) updateRooms() {
	currentRooms := s.gameStorage.Rooms()
	for i := range currentRooms {
		room := &currentRooms[i] // <--- pointer because it's a heavy object.

		switch room.State.(type) {
		case domain.SearchingRoomState:
			s.updateRoomWithSearchingState(room)
		case domain.ChattingRoomState:
			s.updateRoomWithChattingState(room)
		case domain.ChoosingRoomState:
			s.updateRoomWithChoosingState(room)

		case nil:
			s.gameStorage.RemoveRoomById(room.Id) // ?
		}
	}
}

// concrete
// -----------------------------------------------------------------------

func (s *Service) updateRoomWithSearchingState(room *domain.Room) {
	if len(room.Profiles) < 2 {
		return
	}
	state := room.State.(domain.SearchingRoomState)

	// ***

	interval := time.Now().Sub(state.LastUpdateTime)
	if interval < s.config.IntervalFromLastUpdateToNextState {
		return
	}

	// ***

	topic, err := s.controlStorage.RandomTopic(
		context.TODO(), room.Language)
	if err != nil {
		s.handleError(room, err)
		return
	}

	foundGameDate := dto.MakeFoundGameData(0,
		s.config.ChattingStageDuration,
		s.config.ChoosingStageDuration,
		topic.Name,
	)

	foundGameDate.AddProfiles(room.Profiles)
	for i := range room.Profiles {
		foundGameDate.LocalProfileId = i
		serverBody := dto.MakeSvrSearchingGameFoundBody(
			foundGameDate)

		s.asyncResponseChan <- MakeAsyncResponse(
			room.Profiles[i].Id,
			serverBody,
		)
	}

	// *** work with storage.

	s.gameStorage.UpdateRoomToChattingState(room.Id)
}

func (s *Service) updateRoomWithChattingState(room *domain.Room) {
	state := room.State.(domain.SearchingRoomState)
	timeDifference := state.LaunchTime.Sub(time.Now())

	if timeDifference < s.config.ChattingStageDuration {
		return
	}

	// ***

	// TODO: make или new функция
	rs.Rooms[roomInx].State = &domain.ChoosingStateRoom{
		RoomState: domain.RoomState{
			Name:       domain.CHOOSING,
			LaunchTime: time.Now(),
		},
		ProfileIdAndMatchedIds: make(map[string][]string),
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

	// ***

	s.gameStorage.UpdateRoomToChoosingState(room.Id)
}

func (s *Service) updateRoomWithChoosingState(room *domain.Room) {

}

// -----------------------------------------------------------------------

func (s *Service) handleError(room *domain.Room, err error) {
	for i := range room.Profiles {
		s.asyncResponseAboutErrorChan <- MakeAsyncResponseAboutError(
			room.Profiles[i].Id,
			err,
		)
	}

	s.gameStorage.RemoveRoomById(room.Id)
}
