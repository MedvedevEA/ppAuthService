package service

import (
	"log/slog"
	"ppAuthService/internal/repository/repostore"
)

type Service struct {
	store repostore.Repository
	lg    *slog.Logger
}

func New(store repostore.Repository, lg *slog.Logger) *Service {
	return &Service{
		store,
		lg,
	}
}
