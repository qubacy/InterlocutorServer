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

	// TODO: negative reaction if the user is already in the room

	roomLanguage := clientBody.Profile.Language
	room, exist := s.gameStorage.RoomWithSearchingState(roomLanguage)
	if !exist {
		insertedId := s.gameStorage.InsertRoomWithSearchingState(roomLanguage)
		room, exist = s.gameStorage.RoomById(insertedId)

		if !exist {
			return dto.SvrSearchingStartBody{}, ErrRoomNotFound
		}
	}

	// ***

	profile := clientBody.Profile.ToDomain(profileId)
	inserted := s.gameStorage.InsertProfileToRoomWithoutAssignId(profile, room.Id)
	if !inserted {
		return dto.SvrSearchingStartBody{}, ErrFailedToAddProfileToRoom
	}
	success := s.gameStorage.UpdateRoomWithSearchingRoomState(room.Id, time.Now())
	if !success {
		return dto.SvrSearchingStartBody{}, ErrFailedToUpdateRoomWithSearchingState
	}

	// ***

	return dto.SvrSearchingStartBody{}, nil
}

func (s *Service) SearchingStop(profileId string, clientBody dto.CliSearchingStopBody) error {
	s.gameStorage.RemoveProfileById(profileId) // repeated action?
	return nil
}

func (s *Service) ChattingNewMessage(profileId string, clientBody dto.CliChattingNewMessageBody) (
	dto.SvrChattingNewMessageBody, error,
) {
	if !clientBody.IsValid() {
		return dto.MakeSvrChattingNewMessageBodyEmpty(),
			ErrInvalidClientPackBody
	}

	// ***

	// TODO: negative reaction if the user is already in the room

	room, exist := s.gameStorage.RoomWithProfile(profileId)
	if !exist {
		return dto.MakeSvrChattingNewMessageBodyEmpty(),
			ErrProfileIsNotLinkedToRoom
	}
	roomStateName, exist := room.StateName()
	if !exist {
		return dto.MakeSvrChattingNewMessageBodyEmpty(),
			ErrRoomInUnknownState
	}
	if roomStateName != domain.CHATTING {
		return dto.MakeSvrChattingNewMessageBodyEmpty(),
			ErrRoomIsNotInChattingState
	}

	// ***

	profileLocalId, exist := profileIdToLocal(profileId, &room)
	if !exist {
		return dto.MakeSvrChattingNewMessageBodyEmpty(),
			ErrProfileIsNotLinkedToRoom // <--- here impossible!
	}
	serverBody := dto.MakeSvrChattingNewMessageBodyFromParts(
		profileLocalId, clientBody.Message.Text,
	)

	// ***

	for i := range room.Profiles {
		currentProfile := &room.Profiles[i]
		if currentProfile.Id == profileId {
			continue
		}

		s.asyncResponseChan <- MakeAsyncResponse(
			currentProfile.Id, serverBody)
	}

	return serverBody, nil
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
		return ErrRoomIsNotInChoosingState
	}

	// TODO: negative reaction if the user is already in the room

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
