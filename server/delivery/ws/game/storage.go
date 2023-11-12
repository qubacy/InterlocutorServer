package game

import (
	"ilserver/delivery/ws/game/connection"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Storage struct {
	rwMutex     sync.RWMutex
	connections map[string]*connection.Connection
}

func NewStorage() *Storage {
	return &Storage{
		rwMutex:     sync.RWMutex{},
		connections: make(map[string]*connection.Connection),
	}
}

// methods
// -----------------------------------------------------------------------

func (s *Storage) AddOpenConnection(conn *websocket.Conn) *connection.Connection {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	id := uuid.NewString()
	connObj := connection.NewOpenConnection(id, conn)
	s.connections[id] = connObj
	return connObj
}

func (s *Storage) RemoveConnection(id string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	delete(s.connections, id)
}

func (s *Storage) GetConnection(id string) (*connection.Connection, bool) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	value, exist := s.connections[id]
	if exist {
		return value, true
	}
	return nil, false
}
