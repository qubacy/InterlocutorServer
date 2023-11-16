package game

import (
	domain "ilserver/domain/memory"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Storage struct {
	rwMutex sync.RWMutex
	rooms   domain.RoomList
}

func NewStorage() *Storage {
	return &Storage{
		rwMutex: sync.RWMutex{},
		rooms:   domain.RoomList{},
	}
}

var once = sync.Once{}
var instance *Storage = nil

func Instance() *Storage {
	once.Do(func() {
		instance = NewStorage()
	})
	return instance
}

// rooms
// -----------------------------------------------------------------------

func (s *Storage) Rooms() domain.RoomList {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	return s.rooms
}

func (s *Storage) RoomById(id string) (domain.Room, bool) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	roomIndex, exist := s.findRoomIndexById(id)
	if exist {
		return s.rooms[roomIndex], true
	}
	return domain.Room{}, false
}

func (s *Storage) RoomWithSearchingState(language int) (domain.Room, bool) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	for i := range s.rooms {
		room := &s.rooms[i]
		if room.Language == language {
			_, converted :=
				room.State.(domain.SearchingRoomState)

			if converted {
				return s.rooms[i], true
			}
		}
	}
	return domain.Room{}, false
}

func (s *Storage) RoomWithProfile(id string) (domain.Room, bool) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	roomIndex, _, exist :=
		s.findRoomIndexWithProfile(id)
	if exist {
		return s.rooms[roomIndex], true
	}
	return domain.Room{}, false
}

func (s *Storage) InsertRoomWithSearchingState(lang int) string {
	room := domain.MakeEmptyRoomWithSearchingState(lang)
	return s.blockingInsertRoom(room)
}

func (s *Storage) RemoveRoomById(id string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	roomIndex, exist := s.findRoomIndexById(id)
	if exist {
		s.rooms = append(
			s.rooms[:roomIndex],
			s.rooms[roomIndex+1:]...,
		)
	}
}

// -----------------------------------------------------------------------

func (s *Storage) UpdateRoomWithSearchingRoomState(roomId string,
	time time.Time,
) bool {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	roomIndex, exist := s.findRoomIndexById(roomId)
	if !exist {
		return false
	}

	room := &s.rooms[roomIndex]

	state, converted := room.State.(domain.SearchingRoomState)
	if !converted {
		return false
	}

	state.LastUpdateTime = time // !

	room.State = state
	return true
}

func (s *Storage) UpdateRoomWithChoosingState(roomId string,
	profileId string, matchedProfileIdsForProfile []string,
) bool {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	roomIndex, exist := s.findRoomIndexById(roomId)
	if !exist {
		return false
	}

	room := &s.rooms[roomIndex]
	if room.State == nil { // unnecessary check!
		return false
	}

	state, converted := room.State.(domain.ChoosingRoomState)
	if !converted {
		return false
	}

	state.MatchedProfileIdsForProfile[profileId] =
		matchedProfileIdsForProfile

	room.State = state
	return true
}

func (s *Storage) UpdateRoomToChattingState(roomId string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	roomIndex, exist := s.findRoomIndexById(roomId)
	if exist {
		s.rooms[roomIndex].State =
			domain.MakeChattingRoomStateNow()
	}
}

func (s *Storage) UpdateRoomToChoosingState(roomId string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	roomIndex, exist := s.findRoomIndexById(roomId)
	if exist {
		s.rooms[roomIndex].State =
			domain.MakeChoosingRoomStateNow()
	}
}

// bad code?
func (s *Storage) UpdateRoomToNilState(roomId string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	roomIndex, exist := s.findRoomIndexById(roomId)
	if exist {
		s.rooms[roomIndex].State = nil
	}
}

// profile
// -----------------------------------------------------------------------

func (s *Storage) ProfilesByRoomId(id string) domain.ProfileList {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	roomIndex, exist := s.findRoomIndexById(id)
	if exist {
		return s.rooms[roomIndex].Profiles
	}
	return domain.ProfileList{}
}

func (s *Storage) RemoveProfileById(id string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	roomIndex, profileIndex, exist :=
		s.findRoomIndexWithProfile(id)

	if exist {
		s.rooms[roomIndex].
			RemoveProfile(profileIndex)
	}
}

// TODO: or should the profile id only assign storage?
func (s *Storage) InsertProfileToRoomWithoutAssignId(profile domain.Profile, id string) bool {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	roomIndex, exist := s.findRoomIndexById(id)
	if !exist {
		return false
	}

	room := &s.rooms[roomIndex]
	room.Profiles = append(room.Profiles, profile)
	return true
}

// private
// -----------------------------------------------------------------------

func (s *Storage) blockingInsertRoom(room domain.Room) string {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	room.Id = uuid.New().String()
	s.rooms = append(s.rooms, room)
	return room.Id
}

// utility
// -----------------------------------------------------------------------

func (s *Storage) findRoomIndexById(id string) (int, bool) {
	for roomIndex := range s.rooms {
		if s.rooms[roomIndex].Id == id {
			return roomIndex, true
		}
	}
	return 0, false
}

func (s *Storage) findRoomIndexWithProfile(profileId string) (int, int, bool) {
	for roomIndex := range s.rooms {
		profileIndex, exist :=
			s.findProfileIndex(&s.rooms[roomIndex], profileId)

		if exist {
			return roomIndex, profileIndex, true
		}
	}
	return 0, 0, false
}

func (s *Storage) findProfileIndex(room *domain.Room, id string) (int, bool) {
	for i := range room.Profiles {
		if room.Profiles[i].Id == id {
			return i, true
		}
	}
	return 0, false
}
