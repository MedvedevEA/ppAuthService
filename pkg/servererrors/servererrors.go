package servererrors

import "errors"

var (
	RecordNotFound            = errors.New("record not found")
	InternalServerError       = errors.New("internal server error")
	InvalidUsernameOrPassword = errors.New("invalid username or password")
)
