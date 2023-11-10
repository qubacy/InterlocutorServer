package game

import "github.com/gorilla/websocket"

type Connection struct {
	id  string
	wsk *websocket.Conn
}

type Storage struct {
}
