package repository

import (
	"ppAuthService/internal/entity"
	repoDto "ppAuthService/internal/repository/dto"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	AddUser(dto *repoDto.AddUser) (*uuid.UUID, error)
	GetUserByLogin(login string) (*entity.User, error)
	UpdateUser(dto *repoDto.UpdateUser) error
	RemoveUser(userId *uuid.UUID) error

	AddRefreshTokenWithRefreshTokenId(dto *repoDto.AddRefreshTokenWithRefreshTokenId) error
	GetRefreshToken(refreshTokenId *uuid.UUID) (*entity.RefreshToken, error)
	RevokeRefreshTokenByRefreshTokenId(refreshTokenId *uuid.UUID) error
	RevokeRefreshTokensByUserIdAndDeviceCode(dto *repoDto.RevokeRefreshTokensByUserIdAndDeviceCode) error
	RemoveRefreshTokensByExpirationAt(now time.Time) (int64, error)
}
