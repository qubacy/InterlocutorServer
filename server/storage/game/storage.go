package game

import (
	domain "ilserver/domain/memory"
	"sync"

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

// public
// -----------------------------------------------------------------------

func (s *Storage) Rooms() domain.RoomList {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	return s.rooms
}

func (s *Storage) ProfilesByRoomId(id string) domain.ProfileList {
	roomIndex, exist := s.findRoomById(id)
	if exist {
		return s.rooms[roomIndex].Profiles
	}
	return domain.ProfileList{}
}

func (s *Storage) RoomWithProfile(id string) (domain.Room, bool) {
	roomIndex, _, exist :=
		s.findRoomIndexWithProfile(id)
	if exist {
		return s.rooms[roomIndex], true
	}
	return domain.Room{}, false
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

func (s *Storage) InsertRoomWithSearchingState(lang int) {
	room := domain.MakeEmptyRoomWithSearchingState(lang)
	s.insertRoom(room)
}

// -----------------------------------------------------------------------

func (s *Storage) UpdateRoomToChattingState(roomId string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	roomIndex, exist := s.findRoomById(roomId)
	if exist {
		s.rooms[roomIndex].State =
			domain.MakeChattingRoomStateNow()
	}
}

func (s *Storage) UpdateRoomToChoosingState(roomId string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	roomIndex, exist := s.findRoomById(roomId)
	if exist {
		s.rooms[roomIndex].State =
			domain.MakeChattingRoomStateNow()
	}
}

// private
// -----------------------------------------------------------------------

func (s *Storage) insertRoom(room domain.Room) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	room.Id = uuid.New().String()
	s.rooms = append(s.rooms, room)
}

// utility
// -----------------------------------------------------------------------

func (s *Storage) findRoomById(id string) (int, bool) {
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
