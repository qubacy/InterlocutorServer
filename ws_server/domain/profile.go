package domain

import "github.com/gorilla/websocket"

type Profile struct {
	Id uint64

	// TODO: плохой архитектурный паттерн, осторожно
	Conn     *websocket.Conn
	Username string
	Contact  string
}
