package game

import "time"

type Config struct {
	TimeoutForUpdateRooms             time.Duration
	IntervalFromLastUpdateToNextState time.Duration
	ChattingStageDuration             time.Duration
	ChoosingStageDuration             time.Duration
}
