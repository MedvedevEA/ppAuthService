package apiserver

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"ppAuthService/internal/apiserver/middleware"
	"ppAuthService/internal/config"
	"ppAuthService/internal/controller"
	"ppAuthService/pkg/secure"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type ApiServer struct {
	app *fiber.App
	lg  *slog.Logger
	cfg *config.Api
}

func MustNew(ctrl *controller.Controller, lg *slog.Logger, cfg *config.Api) *ApiServer {
	const op = "apiserver.MustNew"
	publicKey, err := secure.LoadPublicKey(cfg.PublicKeyPath)
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}
	app := fiber.New(fiber.Config{
		AppName:      cfg.Name,
		WriteTimeout: cfg.WriteTimeout,
	})
	app.Use(recover.New(recover.ConfigDefault))
	app.Use(middleware.GetLoggerMiddlewareFunc(lg, cfg.Name))

	apiGroup := app.Group("/api")
	v1Group := apiGroup.Group("/v1")
	//public
	v1Group.Post("/login", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	v1Group.Post("/registration", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	//refresh
	v1Group.Post(
		"/refresh",
		middleware.GetAuthMiddlewareFunc(publicKey, "refresh"),
		func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) },
	)
	//access
	authGroup := v1Group.Group(
		"/auth",
		middleware.GetAuthMiddlewareFunc(publicKey, "access"),
	)

	authGroup.Post("/logout", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	authGroup.Post("/unregistration", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })

	authGroup.Post("/messages", ctrl.AddMessage)
	authGroup.Get("/messages/:messageId", ctrl.GetMessage)
	authGroup.Get("/messages", ctrl.GetMessages)
	authGroup.Patch("/messages/:messageId", ctrl.UpdateMessage)
	authGroup.Delete("/messages/:messageId", ctrl.RemoveMessage)

	authGroup.Post("/statuses", ctrl.AddStatus)
	authGroup.Get("/statuses/:statusId", ctrl.GetStatus)
	authGroup.Get("/statuses", ctrl.GetStatuses)
	authGroup.Patch("/statuses/:statusId", ctrl.UpdateStatus)
	authGroup.Delete("/statuses/:statusId", ctrl.RemoveStatus)

	authGroup.Post("/tasks", ctrl.AddTask)
	authGroup.Get("/tasks/:taskId", ctrl.GetTask)
	authGroup.Get("/tasks", ctrl.GetTasks)
	authGroup.Patch("/tasks/:taskId", ctrl.UpdateTask)
	authGroup.Delete("/tasks/:taskId", ctrl.RemoveTask)

	authGroup.Post("/taskusers", ctrl.AddTaskUser)
	authGroup.Get("/taskusers", ctrl.GetTaskUsers)
	authGroup.Delete("/taskusers/:taskUserId", ctrl.RemoveTaskUser)

	authGroup.Get("/users", ctrl.GetUsers)

	app.Use(middleware.BadRequest)

	return &ApiServer{
		app,
		lg,
		cfg,
	}
}

func (a *ApiServer) Run() {
	chError := make(chan error, 1)
	go func() {
		a.lg.Info(fmt.Sprintf("API Server '%s' is started in addr: %s", a.cfg.Name, a.cfg.Addr))
		if err := a.app.Listen(a.cfg.Addr); err != nil {
			chError <- err
		}
	}()
	go func() {
		chQuit := make(chan os.Signal, 1)
		signal.Notify(chQuit, syscall.SIGINT, syscall.SIGTERM)
		<-chQuit
		chError <- a.app.Shutdown()
	}()
	if err := <-chError; err != nil {
		a.lg.Error(fmt.Sprintf("API Server '%s' error", a.cfg.Name), slog.Any("error", err))
		return
	}
	a.lg.Info(fmt.Sprintf("API Server '%s' is stopped", a.cfg.Name))

}
