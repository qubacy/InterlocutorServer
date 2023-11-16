package game

import "errors"

var (
	ErrInvalidClientPackBody                = errors.New("Invalid client pack body")
	ErrRoomNotFound                         = errors.New("Room not found")
	ErrFailedToAddProfileToRoom             = errors.New("Failed to add profile to room")
	ErrProfileIsNotLinkedToRoom             = errors.New("The profile is not linked to the room")
	ErrRoomIsNotInChattingState             = errors.New("The room is not in chatting state")
	ErrRoomIsNotInChoosingState             = errors.New("The room is not in choosing state")
	ErrRoomInUnknownState                   = errors.New("Room in unknown state")
	ErrFailedToUpdateRoomWithChoosingState  = errors.New("Failed to update room with choosing state")
	ErrFailedToUpdateRoomWithSearchingState = errors.New("Failed to update room with searching state")
)
