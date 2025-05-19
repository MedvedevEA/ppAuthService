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

func (c *Controller) AddTask(ctx *fiber.Ctx) error {
	req := new(ctrlDto.AddTask)
	if err := ctx.BodyParser(req); err != nil {
		err := fmt.Errorf("controller.AddTask: %w (%v)", ctrlErr.ErrBodyParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.AddTask: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	task, err := c.service.AddTask(&repoStoreDto.AddTask{
		StatusId:    req.StatusId,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(201).JSON(task)
}
func (c *Controller) GetTask(ctx *fiber.Ctx) error {
	req := new(ctrlDto.GetTask)
	if err := ctx.ParamsParser(req); err != nil {
		err := fmt.Errorf("controller.GetTask: %w (%v)", ctrlErr.ErrParamsParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.GetTask: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	task, err := c.service.GetTask(req.TaskId)
	if err != nil {
		if errors.Is(err, repoStoreErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(task)
}
func (c *Controller) GetTasks(ctx *fiber.Ctx) error {
	req := &ctrlDto.GetTasks{
		Offset: 0,
		Limit:  10,
	}
	if err := ctx.QueryParser(req); err != nil {
		err := fmt.Errorf("controller.GetTasks: %w (%v)", ctrlErr.ErrQueryParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.GetTasks: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	tasks, err := c.service.GetTasks(&repoStoreDto.GetTasks{
		Offset:   req.Offset,
		Limit:    req.Limit,
		StatusId: req.StatusId,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(tasks)
}
func (c *Controller) UpdateTask(ctx *fiber.Ctx) error {
	req := new(ctrlDto.UpdateTask)
	if err := ctx.BodyParser(req); err != nil {
		err := fmt.Errorf("controller.UpdateTask: %w (%v)", ctrlErr.ErrBodyParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := ctx.ParamsParser(req); err != nil {
		err := fmt.Errorf("controller.UpdateTask: %w (%v)", ctrlErr.ErrParamsParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.UpdateTask: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	task, err := c.service.UpdateTask(&repoStoreDto.UpdateTask{
		TaskId:      req.TaskId,
		StatusId:    req.StatusId,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, repoStoreErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(task)
}
func (c *Controller) RemoveTask(ctx *fiber.Ctx) error {
	req := new(ctrlDto.RemoveTask)
	if err := ctx.ParamsParser(req); err != nil {
		err := fmt.Errorf("controller.RemoveTask: %w (%v)", ctrlErr.ErrParamsParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.RemoveTask: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	err := c.service.RemoveTask(req.TaskId)
	if err != nil {
		if errors.Is(err, repoStoreErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.SendStatus(204)
}
