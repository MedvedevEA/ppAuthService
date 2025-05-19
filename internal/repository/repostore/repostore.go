package repostore

import (
	"ppAuthService/internal/entity"
	"ppAuthService/internal/repository/repostore/dto"

	"github.com/google/uuid"
)

type Repository interface {
	AddMessage(dto *dto.AddMessage) (*entity.Message, error)
	GetMessage(messageId *uuid.UUID) (*entity.Message, error)
	GetMessages(dto *dto.GetMessages) ([]*entity.Message, error)
	UpdateMessage(dto *dto.UpdateMessage) (*entity.Message, error)
	RemoveMessage(messageId *uuid.UUID) error

	AddStatus(name string) (*entity.Status, error)
	GetStatus(statusId *uuid.UUID) (*entity.Status, error)
	GetStatuses() ([]*entity.Status, error)
	UpdateStatus(dto *dto.UpdateStatus) (*entity.Status, error)
	RemoveStatus(statusId *uuid.UUID) error

	AddTask(dto *dto.AddTask) (*entity.Task, error)
	GetTask(taskId *uuid.UUID) (*entity.Task, error)
	GetTasks(dto *dto.GetTasks) ([]*entity.Task, error)
	UpdateTask(dto *dto.UpdateTask) (*entity.Task, error)
	RemoveTask(taskId *uuid.UUID) error

	AddTaskUser(dto *dto.AddTaskUser) (*entity.TaskUser, error)
	GetTaskUsers(dto *dto.GetTaskUsers) ([]*entity.TaskUser, error)
	RemoveTaskUser(taskUserID *uuid.UUID) error

	AddUserWithUserId(dto *dto.AddUser) (*entity.User, error)
	GetUsers(dto *dto.GetUsers) ([]*entity.User, error)
	RemoveUser(userId *uuid.UUID) error
}
