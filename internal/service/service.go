package service

import (
	"crypto/rsa"
	"log/slog"
	"ppAuthService/internal/config"
	"ppAuthService/internal/repository"
	"time"

	proto "github.com/MedvedevEA/ppProtos/gen/auth"
)

type Service struct {
	proto.UnimplementedAuthServiceServer
	store           repository.Repository
	privateKey      *rsa.PrivateKey
	accessLifetime  time.Duration
	refrashLifetime time.Duration
	lg              *slog.Logger
}

func MustNew(store repository.Repository, lg *slog.Logger, cfg *config.Token) *Service {

	return &Service{
		store:           store,
		privateKey:      nil,
		accessLifetime:  cfg.AccessLifetime,
		refrashLifetime: cfg.RefreshLifetime,
		lg:              lg,
	}
}
