package controller

import (
	"errors"
	"fmt"
	ctrlDto "ppAuthService/internal/controller/dto"
	ctrlErr "ppAuthService/internal/controller/err"
	repoStoreDto "ppAuthService/internal/repository/repostore/dto"
	repoStoreErr "ppAuthService/internal/repository/repostore/err"
	"ppAuthService/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func (c *Controller) AddTaskUser(ctx *fiber.Ctx) error {
	req := new(ctrlDto.AddTaskUser)
	if err := ctx.BodyParser(req); err != nil {
		err := fmt.Errorf("controller.AddTaskUser: %w (%v)", ctrlErr.ErrBodyParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.AddTaskUser: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	taskUser, err := c.service.AddTaskUser(&repoStoreDto.AddTaskUser{
		TaskId: req.TaskId,
		UserId: req.UserId,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(201).JSON(taskUser)
}
func (c *Controller) GetTaskUsers(ctx *fiber.Ctx) error {
	req := &ctrlDto.GetTaskUsers{
		Offset: 0,
		Limit:  10,
	}
	if err := ctx.QueryParser(req); err != nil {
		err := fmt.Errorf("controller.GetTaskUsers: %w (%v)", ctrlErr.ErrQueryParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.GetTaskUsers: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	taskUsers, err := c.service.GetTaskUsers(&repoStoreDto.GetTaskUsers{
		Offset: req.Offset,
		Limit:  req.Limit,
		TaskId: req.TaskId,
		UserId: req.UserId,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(taskUsers)
}
func (c *Controller) RemoveTaskUser(ctx *fiber.Ctx) error {
	req := new(ctrlDto.RemoveTaskUser)
	if err := ctx.ParamsParser(req); err != nil {
		err := fmt.Errorf("controller.RemoveTaskUser: %w (%v)", ctrlErr.ErrParamsParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.RemoveTaskUser: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	err := c.service.RemoveTaskUser(req.TaskUserId)
	if err != nil {
		if errors.Is(err, repoStoreErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.SendStatus(204)
}
