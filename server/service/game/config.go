package game

import "time"

type Config struct {
	RoomUpdateDuration                time.Duration
	TimeoutForUpdateRooms             time.Duration
	IntervalFromLastUpdateToNextState time.Duration
	ChattingStageDuration             time.Duration
	ChoosingStageDuration             time.Duration
}
