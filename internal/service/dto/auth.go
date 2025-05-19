package dto

import (
	"github.com/google/uuid"
)

type RegistrationRequest struct {
	Name     string
	Password string
}
type LoginRequest struct {
	Name       string
	DeviceCode string
	Password   string
}
type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}
type LogoutRequest struct {
	UserId     *uuid.UUID
	DeviceCode string
}
type UpdatePasswordRequest struct {
	UserId      *uuid.UUID
	NewPassword string
}
type RefreshTokenResponse struct {
	AccessToken  string
	RefreshToken string
}
