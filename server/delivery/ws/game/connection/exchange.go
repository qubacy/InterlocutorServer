package connection

import "github.com/gorilla/websocket"

type Message struct {
	Type int
	Data []byte
}

func MakeMessage(messageType int, data []byte) Message {
	return Message{
		Type: messageType,
		Data: data,
	}
}

func MakeTextMessage(data []byte) Message {
	return MakeMessage(
		websocket.TextMessage,
		data,
	)
}
