package transport

//TODO: нужны ли коды операций в ответе для клиента?

const (
	SEARCHING_START int = iota
	SEARCHING_STOP
	SEARCHING_GAME_FOUND

	CHATTING_NEW_MESSAGE
	CHATTING_STAGE_IS_OVER

	CHOOSING_USERS_CHOSEN
	CHOOSING_STAGE_IS_OVER
)
