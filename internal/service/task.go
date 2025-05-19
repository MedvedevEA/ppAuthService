package service

import (
	"ppAuthService/internal/entity"
	repoStoreDto "ppAuthService/internal/repository/repostore/dto"

	"github.com/google/uuid"
)

func (s *Service) AddTask(dto *repoStoreDto.AddTask) (*entity.Task, error) {
	return s.store.AddTask(dto)
}
func (s *Service) GetTask(taskId *uuid.UUID) (*entity.Task, error) {
	return s.store.GetTask(taskId)
}
func (s *Service) GetTasks(dto *repoStoreDto.GetTasks) ([]*entity.Task, error) {
	return s.store.GetTasks(dto)
}
func (s *Service) UpdateTask(dto *repoStoreDto.UpdateTask) (*entity.Task, error) {
	return s.store.UpdateTask(dto)
}
func (s *Service) RemoveTask(taskId *uuid.UUID) error {
	return s.store.RemoveTask(taskId)
}
