package service

import (
	"ppAuthService/internal/entity"
	repoStoreDto "ppAuthService/internal/repository/repostore/dto"
)

func (s *Service) GetUsers(dto *repoStoreDto.GetUsers) ([]*entity.User, error) {
	return s.store.GetUsers(dto)
}
