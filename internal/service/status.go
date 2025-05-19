package service

import (
	"ppAuthService/internal/entity"
	repoStoreDto "ppAuthService/internal/repository/repostore/dto"

	"github.com/google/uuid"
)

func (s *Service) AddStatus(name string) (*entity.Status, error) {
	return s.store.AddStatus(name)
}
func (s *Service) GetStatus(StatusId *uuid.UUID) (*entity.Status, error) {
	return s.store.GetStatus(StatusId)
}
func (s *Service) GetStatuses() ([]*entity.Status, error) {
	return s.store.GetStatuses()
}
func (s *Service) UpdateStatus(dto *repoStoreDto.UpdateStatus) (*entity.Status, error) {
	return s.store.UpdateStatus(dto)
}
func (s *Service) RemoveStatus(StatusId *uuid.UUID) error {
	return s.store.RemoveStatus(StatusId)
}
