package service

import (
	"ppAuthService/internal/entity"
	repoStoreDto "ppAuthService/internal/repository/repostore/dto"

	"github.com/google/uuid"
)

func (s *Service) AddTaskUser(dto *repoStoreDto.AddTaskUser) (*entity.TaskUser, error) {
	return s.store.AddTaskUser(dto)
}
func (s *Service) GetTaskUsers(dto *repoStoreDto.GetTaskUsers) ([]*entity.TaskUser, error) {
	return s.store.GetTaskUsers(dto)
}
func (s *Service) RemoveTaskUser(TaskUserId *uuid.UUID) error {
	return s.store.RemoveTaskUser(TaskUserId)
}
