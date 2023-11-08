package memory

import (
	domain "ilserver/domain/memory"
	"sync"
)

type Storage struct {
	rwMutex sync.RWMutex
	rooms   []domain.Room
}

// public
// -----------------------------------------------------------------------

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
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	room := domain.MakeEmptyRoomWithSearchingState(lang)
	s.rooms = append(s.rooms, room)
}

// TODO: Update...

// private
// -----------------------------------------------------------------------

func (s *Storage) findRoomIndexWithProfile(id string) (int, int, bool) {
	for roomIndex := range s.rooms {
		profileIndex, exist :=
			s.findProfileIndex(&s.rooms[roomIndex], id)

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
