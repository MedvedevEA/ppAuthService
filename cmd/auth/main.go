package main

import (
	"ppAuthService/internal/config"
	"ppAuthService/internal/logger"
	"ppAuthService/internal/server"
	"ppAuthService/internal/service"
	"ppAuthService/internal/store"
)

func main() {

	cfg := config.MustNew()
	lg := logger.MustNew(cfg.Env)
	store := store.MustNew(lg, &cfg.Store)
	service := service.MustNew(store, lg, &cfg.Token)
	server := server.MustNew(service, lg, &cfg.Server)

	server.Start()

}
