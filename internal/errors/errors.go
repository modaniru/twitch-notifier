package mye

import "errors"

var ErrFollowNotFound = errors.New("follow was not found")
var ErrStreamerNotFound = errors.New("streamer was not found")
var ErrFollowAlreadyExists = errors.New("follow already exists")
var UserWasNotFound = errors.New("user was not found")
