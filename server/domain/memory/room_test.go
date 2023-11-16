package memory

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_DirectCopyRoom(t *testing.T) {
	baseRoom := MakeEmptyRoomWithSearchingState(1)
	baseRoom.Id = uuid.NewString()

	// ***

	copyRoom := baseRoom
	copyRoom.Id = uuid.NewString()

	time.Sleep(500 * time.Millisecond)

	state, ok := copyRoom.State.(SearchingRoomState) // copy state.
	if ok {
		state.LaunchTime = time.Now()
		copyRoom.State = state // reset.
	} else {
		t.Fatalf("Wrong type state")
	}

	fmt.Println(baseRoom)
	fmt.Println(copyRoom)
}

func Test_DirectCopyRoom_v1(t *testing.T) {
	baseRoom := MakeEmptyRoomWithSearchingState(1)

	baseRoom.Id = uuid.NewString()
	baseRoom.State = NewSearchingRoomState(time.Now())

	// *** deep copy very hard...

	copyRoom := baseRoom
	copyRoom.Id = uuid.NewString()
	copyRoom.State = new(SearchingRoomState)
	searchingRoomState := copyRoom.State.(*SearchingRoomState)
	*searchingRoomState = *baseRoom.State.(*SearchingRoomState)

	time.Sleep(500 * time.Millisecond)

	state, ok := baseRoom.State.(*SearchingRoomState)
	if ok {
		state.LaunchTime = time.Now()
	} else {
		t.Fatalf("Wrong type state")
	}

	fmt.Println(baseRoom.State)
	fmt.Println(copyRoom.State)
}

// experiments
// -----------------------------------------------------------------------

func Test_uuid_NewString(t *testing.T) {
	fmt.Println(uuid.NewString())
}

func Test_interface(t *testing.T) {
	var value interface{} = ChattingRoomState{}
	state, converted := value.(ChattingRoomState)
	fmt.Println(state, converted)

	value = nil
	state, converted = value.(ChattingRoomState)
	fmt.Println(state, converted)

	// ***

	switch value.(type) {
	case ChattingRoomState:
		fmt.Println("value is instance of ChattingRoomState")
	case nil:
		fmt.Println("value is nil")
	}
}
