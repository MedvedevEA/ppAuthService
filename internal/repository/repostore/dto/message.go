package dto

import "github.com/google/uuid"

type AddMessage struct {
	TaskId *uuid.UUID `json:"task_id"`
	UserId *uuid.UUID `json:"user_id"`
	Text   string     `json:"text"`
}
type GetMessages struct {
	Offset int        `json:"offset"`
	Limit  int        `json:"limit"`
	TaskId *uuid.UUID `json:"task_id"`
}
type UpdateMessage struct {
	MessageId *uuid.UUID `json:"message_id"`
	Text      *string    `json:"text"`
}
