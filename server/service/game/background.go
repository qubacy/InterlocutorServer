package game

import (
	"context"
	domain "ilserver/domain/memory"
	"ilserver/service/game/dto"
	"time"
)

func (s *Service) backgroundUpdates(ctx context.Context) {
	timeout := s.config.TimeoutForUpdateRooms
	for { // <--- endless!
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
	state := room.State.(domain.SearchingRoomState)

	// ***

	interval := time.Now().Sub(state.LastUpdateTime)
	if interval < s.config.IntervalFromLastUpdateToNextState {
		return
	}

	if len(room.Profiles) == 0 {
		s.gameStorage.RemoveRoomById(room.Id)
		return
	}

	if len(room.Profiles) < 2 {
		return
	}

	// ***

	ctx, cancel := context.WithTimeout(context.Background(),
		s.config.RoomUpdateDuration)
	defer cancel()

	topic, err := s.controlStorage.RandomTopic(
		ctx, room.Language)
	if err != nil {
		s.handleError(room, err)
		return
	}

	// prepare response pack.
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
	if len(room.Profiles) == 0 {
		s.gameStorage.RemoveRoomById(room.Id)
		return
	}

	state := room.State.(domain.ChattingRoomState)
	timeDifference := time.Now().Sub(state.LaunchTime)

	if timeDifference < s.config.ChattingStageDuration {
		return
	}

	// ***

	serverBody := dto.MakeSvrChattingStageIsOverBodyEmpty()
	for i := range room.Profiles {
		s.asyncResponseChan <- MakeAsyncResponse(
			room.Profiles[i].Id, serverBody)
	}

	// ***

	s.gameStorage.UpdateRoomToChoosingState(room.Id)
}

func (s *Service) updateRoomWithChoosingState(room *domain.Room) {
	if len(room.Profiles) == 0 {
		s.gameStorage.RemoveRoomById(room.Id)
		return
	}

	state := room.State.(domain.ChoosingRoomState)
	timeDifference := time.Now().Sub(state.LaunchTime)

	if timeDifference < s.config.ChoosingStageDuration {
		return
	}

	// ***

	for _, currentProfile := range room.Profiles {
		matchedUsers := dto.MatchedUserList{}
		matchedProfileIds := state.MatchedProfileIdsForProfile[currentProfile.Id]

		for _, matchedId := range matchedProfileIds {

			// mutual selection check!
			if containsElementInList(
				currentProfile.Id,
				state.MatchedProfileIdsForProfile[matchedId], // <--- selected ids of another user
			) {
				profileLocalId, exist := profileIdToLocal(matchedId, room)
				if !exist {
					s.handleError(room, ErrProfileIsNotLinkedToRoom) // <--- if the user disconnects
					return
				}

				matchedProfile, exist := profileByLocalId(profileLocalId, room)
				if !exist {
					s.handleError(room, ErrProfileIsNotLinkedToRoom) // <--- impossible...
					return
				}

				matchedUsers.Add(
					dto.MakeMatchedUser(
						profileLocalId,
						matchedProfile.Contact,
					),
				)
			}
		}

		serverBody := dto.MakeSvrChoosingStageIsOverBody(matchedUsers)
		s.asyncResponseChan <- MakeAsyncResponse(
			currentProfile.Id, serverBody)

		// TODO: then disconnect the user? The client can be the initiator...
	}

	// ***

	s.gameStorage.UpdateRoomToNilState(room.Id)
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

func containsElementInList(value string, list []string) bool {
	for i := range list {
		if list[i] == value {
			return true
		}
	}
	return false
}
