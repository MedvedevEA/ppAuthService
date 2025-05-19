package main

import (
	"context"

	"skillsRockTodo/internal/apiserver"
	"skillsRockTodo/internal/config"
	"skillsRockTodo/internal/controller"
	"skillsRockTodo/internal/infrastructure/store"
	"skillsRockTodo/internal/logger"
	"skillsRockTodo/internal/service"
)

func main() {
	cfg := config.MustLoad()

	lg := logger.MustNew(cfg.Env)

	store := store.MustNew(context.Background(), lg, &cfg.Store)

	service := service.New(store, lg)

	controller := controller.New(service, lg)

	apiServer := apiserver.MustNew(controller, lg, &cfg.Api)

	apiServer.Run()
}
