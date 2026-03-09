package scheduler

import "errors"

var (
	ErrCloseChanNotFound     = errors.New("close channel not found")
	ErrCloseChanAlreadyExist = errors.New("close channel already exist")
)
