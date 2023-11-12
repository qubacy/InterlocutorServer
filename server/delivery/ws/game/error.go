package game

import "errors"

var (
	ErrConnectionObjectNotFound = errors.New("Connection object not found")
	ErrUnknownMessageOperation  = errors.New("Unknown message operation")
)
