package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	MessageId *uuid.UUID `json:"messageId"`
	TaskId    *uuid.UUID `json:"taskId"`
	UserId    *uuid.UUID `json:"userId"`
	Text      string     `json:"text"`
	CreateAt  time.Time  `json:"createAt"`
	UpdateAt  *time.Time `json:"updateAt"`
}
type Status struct {
	StatusId *uuid.UUID `json:"statusId"`
	Name     string     `json:"name"`
}
type Task struct {
	TaskId      *uuid.UUID `json:"taskId"`
	StatusId    *uuid.UUID `json:"statusId"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
}
type TaskUser struct {
	TaskUserId *uuid.UUID `json:"taskUserId"`
	TaskId     *uuid.UUID `json:"taskId"`
	UserId     *uuid.UUID `json:"userId"`
}
type User struct {
	UserId *uuid.UUID `json:"userId"`
	Name   string     `json:"name"`
}
