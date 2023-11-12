package connection

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
