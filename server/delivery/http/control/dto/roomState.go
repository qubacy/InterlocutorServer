package dto

import (
	domain "ilserver/domain/memory"
)

// direct mapping
// -----------------------------------------------------------------------

type RoomState struct {
	Name       string  `json:"name"`
	LaunchTime float64 `json:"launch-time"`
}

func makeRoomState(state domain.RoomState) RoomState {
	return RoomState{
		Name:       state.Name.ToString(),
		LaunchTime: float64(state.LaunchTime.Unix()),
	}
}

type SearchingRoomState struct {
	RoomState      `json:"room-state"`
	LastUpdateTime float64 `json:"last-update-time"`
}

func makeSearchingRoomState(state domain.SearchingRoomState) SearchingRoomState {
	return SearchingRoomState{
		RoomState:      makeRoomState(state.RoomState),
		LastUpdateTime: float64(state.LastUpdateTime.Unix()),
	}
}

type ChattingRoomState struct {
	RoomState `json:"room-state"`
}

func makeChattingRoomState(state domain.ChattingRoomState) ChattingRoomState {
	return ChattingRoomState{
		RoomState: makeRoomState(state.RoomState),
	}
}

type ChoosingRoomState struct {
	RoomState                   `json:"room-state"`
	MatchedProfileIdsForProfile map[string][]string `json:"matched-profile-ids-for-profile"`
}

func makeChoosingRoomState(state domain.ChoosingRoomState) ChoosingRoomState {
	return ChoosingRoomState{
		RoomState:                   makeRoomState(state.RoomState),
		MatchedProfileIdsForProfile: state.MatchedProfileIdsForProfile,
	}
}

// no check... Checks the top level
func MakeRoomState(i interface{}) interface{} {
	switch value := i.(type) {

	case domain.SearchingRoomState:
		return makeSearchingRoomState(value)

	case domain.ChattingRoomState:
		return makeChattingRoomState(value)

	case domain.ChoosingRoomState:
		return makeChoosingRoomState(value)

	}
	return nil
}
