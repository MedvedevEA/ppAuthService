package service

import (
	"ppAuthService/internal/entity"
	repoStoreDto "ppAuthService/internal/repository/repostore/dto"

	"github.com/google/uuid"
)

func (s *Service) AddMessage(dto *repoStoreDto.AddMessage) (*entity.Message, error) {
	return s.store.AddMessage(dto)
}
func (s *Service) GetMessage(MessageId *uuid.UUID) (*entity.Message, error) {
	return s.store.GetMessage(MessageId)
}
func (s *Service) GetMessages(dto *repoStoreDto.GetMessages) ([]*entity.Message, error) {
	return s.store.GetMessages(dto)
}
func (s *Service) UpdateMessage(dto *repoStoreDto.UpdateMessage) (*entity.Message, error) {
	return s.store.UpdateMessage(dto)
}
func (s *Service) RemoveMessage(MessageId *uuid.UUID) error {
	return s.store.RemoveMessage(MessageId)
}
