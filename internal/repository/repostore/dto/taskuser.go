package dto

import "github.com/google/uuid"

type AddTaskUser struct {
	TaskId *uuid.UUID `json:"task_id"`
	UserId *uuid.UUID `json:"user_id"`
}
type GetTaskUsers struct {
	Offset int        `json:"offset"`
	Limit  int        `json:"limit"`
	TaskId *uuid.UUID `json:"task_id"`
	UserId *uuid.UUID `json:"user_id"`
}
