package errors

import "errors"

var (
	ErrChatNotFound  = errors.New("user not found")
	ErrUserNotInChat = errors.New("user not in chat")
)
