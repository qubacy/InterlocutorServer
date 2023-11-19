package game

import (
	"context"
	domain "ilserver/domain/memory"
	"ilserver/service/game/dto"
	"ilserver/storage/control"
	"ilserver/storage/game"
	"sync"
	"time"
)

type Service struct {
	config         Config
	gameStorage    *game.Storage // <--- impl
	controlStorage control.Storage
	transactionMx  sync.Mutex // will help eliminate rare moments?

	asyncResponseChan           chan AsyncResponse
	asyncResponseAboutErrorChan chan AsyncResponseAboutError
}

func NewService(
	ctx context.Context,
	config Config,
	gameStorage *game.Storage,
	controlStorage control.Storage,
) *Service {
	instance := &Service{
		config:         config,
		gameStorage:    gameStorage,
		controlStorage: controlStorage,
		transactionMx:  sync.Mutex{},

		asyncResponseChan:           make(chan AsyncResponse, 9),
		asyncResponseAboutErrorChan: make(chan AsyncResponseAboutError, 9),
	}

	go instance.backgroundUpdates(ctx)
	return instance
}

// public
// -----------------------------------------------------------------------

func (s *Service) AsyncResponse() <-chan AsyncResponse {
	return s.asyncResponseChan
}

func (s *Service) AsyncResponseAboutError() <-chan AsyncResponseAboutError {
	return s.asyncResponseAboutErrorChan
}

// -----------------------------------------------------------------------

func (s *Service) SearchingStart(profileId string, clientBody dto.CliSearchingStartBody) (
	dto.SvrSearchingStartBody, error,
) {
	if !clientBody.IsValid() {
		return dto.MakeSvrSearchingStartBodyEmpty(), ErrInvalidClientPackBody
	}

	// ***

	// user is already in the room...
	_, exist := s.gameStorage.RoomWithProfile(profileId)
	if exist {
		return dto.SvrSearchingStartBody{}, nil
	}

	roomLanguage := clientBody.Profile.Language
	room, exist := s.gameStorage.RoomWithSearchingState(roomLanguage)
	if !exist {
		insertedId := s.gameStorage.InsertRoomWithSearchingState(roomLanguage)
		room, exist = s.gameStorage.RoomById(insertedId)

		if !exist {
			return dto.MakeSvrSearchingStartBodyEmpty(), ErrRoomNotFound
		}
	}

	// ***

	profile := clientBody.Profile.ToDomain(profileId)
	inserted := s.gameStorage.InsertProfileToRoomWithoutAssignId(profile, room.Id)
	if !inserted {
		return dto.MakeSvrSearchingStartBodyEmpty(), ErrFailedToAddProfileToRoom
	}
	success := s.gameStorage.UpdateRoomWithSearchingRoomState(room.Id, time.Now())
	if !success {
		return dto.MakeSvrSearchingStartBodyEmpty(), ErrFailedToUpdateRoomWithSearchingState
	}

	// ***

	return dto.SvrSearchingStartBody{}, nil
}

func (s *Service) SearchingStop(profileId string, clientBody dto.CliSearchingStopBody) error {
	s.gameStorage.RemoveProfileById(profileId) // repeated action?
	return nil
}

func (s *Service) ChattingNewMessage(profileId string, clientBody dto.CliChattingNewMessageBody) error {
	if !clientBody.IsValid() {
		return ErrInvalidClientPackBody
	}

	// ***

	room, exist := s.gameStorage.RoomWithProfile(profileId)
	if !exist {
		return ErrProfileIsNotLinkedToRoom
	}
	roomStateName, exist := room.StateName()
	if !exist {
		return ErrRoomInUnknownState
	}
	if roomStateName != domain.CHATTING {
		return nil // or ErrRoomIsNotInChattingState?
	}

	// ***

	profileLocalId, exist := profileIdToLocal(profileId, &room)
	if !exist {
		return ErrProfileIsNotLinkedToRoom // <--- here impossible!
	}
	serverBody := dto.MakeSvrChattingNewMessageBodyFromParts(
		profileLocalId, clientBody.Message.Text,
	)

	// ***

	for i := range room.Profiles {
		currentProfile := &room.Profiles[i]
		s.asyncResponseChan <- MakeAsyncResponse(
			currentProfile.Id, serverBody)
	}

	return nil
}

func (s *Service) ChoosingUsersChosen(
	profileId string, clientBody dto.CliChoosingUsersChosenBody,
) error {
	room, exist := s.gameStorage.RoomWithProfile(profileId)
	if !exist {
		return ErrProfileIsNotLinkedToRoom
	}
	roomStateName, exist := room.StateName()
	if !exist {
		return ErrRoomInUnknownState
	}
	if roomStateName != domain.CHOOSING {
		return nil // or ErrRoomIsNotInChoosingState?
	}

	// ***

	matchedProfileIdsForProfile := []string{}
	for _, userLocalId := range clientBody.UserIdList { // index.
		profile := &room.Profiles[userLocalId]
		matchedProfileIdsForProfile = append(
			matchedProfileIdsForProfile,
			profile.Id,
		)
	}

	success := s.gameStorage.UpdateRoomWithChoosingState(
		room.Id, profileId, matchedProfileIdsForProfile) // <--- influence!
	if !success {
		return ErrFailedToUpdateRoomWithChoosingState
	}

	return nil
}

// -----------------------------------------------------------------------

func (s *Service) ProfileLeftGame(profileId string) {
	s.gameStorage.RemoveProfileById(profileId)
}

// hidden functions
// -----------------------------------------------------------------------

func profileByLocalId(id int, room *domain.Room) (domain.Profile, bool) {
	if id < 0 || id >= len(room.Profiles) {
		return domain.Profile{}, false
	}

	return room.Profiles[id], true
}

func profileIdToLocal(id string, room *domain.Room) (int, bool) {
	for i := range room.Profiles {
		if room.Profiles[i].Id == id {
			return i, true
		}
	}
	return 0, false
}
