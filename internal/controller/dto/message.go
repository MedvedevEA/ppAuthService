package dto

import (
	"github.com/google/uuid"
)

type AddMessage struct {
	TaskId *uuid.UUID `json:"taskId" validate:"required"`
	UserId *uuid.UUID `json:"userId" validate:"required"`
	Text   string     `json:"text"`
}
type GetMessage struct {
	MessageId *uuid.UUID `uri:"messageId" validate:"required"`
}
type GetMessages struct {
	Offset int        `query:"offset" validate:"gte=0"`
	Limit  int        `query:"limit" validate:"gte=0"`
	TaskId *uuid.UUID `query:"taskId"`
}
type UpdateMessage struct {
	MessageId *uuid.UUID `uri:"MessageId" validate:"required"`
	Text      *string    `json:"text"`
}
type RemoveMessage struct {
	MessageId *uuid.UUID `uri:"MessageId" validate:"required"`
}
