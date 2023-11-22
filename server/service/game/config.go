package game

import "time"

type Config struct {
	RoomUpdateDuration                time.Duration
	TimeoutForUpdateRooms             time.Duration
	IntervalFromLastUpdateToNextState time.Duration
	ChattingStageDuration             time.Duration
	ChoosingStageDuration             time.Duration
	MaxProfileCountInRoom             int
}

func MakeConfig(
	roomUpdateDuration time.Duration,
	timeoutForUpdateRooms time.Duration,
	intervalFromLastUpdateToNextState time.Duration,
	chattingStageDuration time.Duration,
	choosingStageDuration time.Duration,
	maxProfileCountInRoom int,
) Config {
	return Config{
		RoomUpdateDuration:                roomUpdateDuration,
		TimeoutForUpdateRooms:             timeoutForUpdateRooms,
		IntervalFromLastUpdateToNextState: intervalFromLastUpdateToNextState,
		ChattingStageDuration:             chattingStageDuration,
		ChoosingStageDuration:             choosingStageDuration,
		MaxProfileCountInRoom:             maxProfileCountInRoom,
	}
}
