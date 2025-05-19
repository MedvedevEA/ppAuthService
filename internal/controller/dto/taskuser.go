package dto

import (
	"github.com/google/uuid"
)

type AddTaskUser struct {
	TaskId *uuid.UUID `json:"taskId" validate:"required"`
	UserId *uuid.UUID `json:"userId" validate:"required"`
}
type GetTaskUsers struct {
	Offset int        `query:"offset" validate:"gte=0"`
	Limit  int        `query:"limit" validate:"gte=0"`
	UserId *uuid.UUID `query:"userId"`
	TaskId *uuid.UUID `query:"taskId"`
}
type RemoveTaskUser struct {
	TaskUserId *uuid.UUID `uri:"taskUserId" validate:"required"`
}
