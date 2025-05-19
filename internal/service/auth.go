package service

import (
	svcDto "ppAuthService/internal/service/dto"

	"github.com/google/uuid"
)

func (s *Service) Registration(dto *svcDto.RegistrationRequest) error {
	return nil
}
func (s *Service) Unregistration(userId *uuid.UUID) error {
	return nil
}
func (s *Service) Login(dto *svcDto.LoginRequest) (*svcDto.LoginResponse, error) {
	return nil, nil
}
func (s *Service) Logout(dto *svcDto.LogoutRequest) error {
	return nil
}
func (s *Service) UpdatePassword(dto *svcDto.UpdatePasswordRequest) error {
	return nil
}
func (s *Service) RefreshToken(tokenId *uuid.UUID) (*svcDto.RefreshTokenResponse, error) {
	return nil, nil
}
