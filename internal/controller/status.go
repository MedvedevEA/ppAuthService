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

func (c *Controller) AddStatus(ctx *fiber.Ctx) error {
	req := new(ctrlDto.AddStatus)
	if err := ctx.BodyParser(req); err != nil {
		err := fmt.Errorf("controller.AddStatus: %w (%v)", ctrlErr.ErrBodyParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.AddStatus: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	status, err := c.service.AddStatus(req.Name)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(201).JSON(status)
}
func (c *Controller) GetStatus(ctx *fiber.Ctx) error {
	req := new(ctrlDto.GetStatus)
	if err := ctx.ParamsParser(req); err != nil {
		err := fmt.Errorf("controller.GetStatus: %w (%v)", ctrlErr.ErrParamsParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.GetStatus: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	status, err := c.service.GetStatus(req.StatusId)
	if err != nil {
		if errors.Is(err, repoStoreErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(status)
}
func (c *Controller) GetStatuses(ctx *fiber.Ctx) error {
	statuses, err := c.service.GetStatuses()
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(statuses)
}
func (c *Controller) UpdateStatus(ctx *fiber.Ctx) error {
	req := new(ctrlDto.UpdateStatus)
	if err := ctx.BodyParser(req); err != nil {
		err := fmt.Errorf("controller.UpdateStatus: %w (%v)", ctrlErr.ErrBodyParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := ctx.ParamsParser(req); err != nil {
		err := fmt.Errorf("controller.UpdateStatus: %w (%v)", ctrlErr.ErrParamsParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.UpdateStatus: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	Status, err := c.service.UpdateStatus(&repoStoreDto.UpdateStatus{
		StatusId: req.StatusId,
		Name:     req.Name,
	})
	if err != nil {
		if errors.Is(err, repoStoreErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(Status)
}
func (c *Controller) RemoveStatus(ctx *fiber.Ctx) error {
	req := new(ctrlDto.RemoveStatus)
	if err := ctx.ParamsParser(req); err != nil {
		err := fmt.Errorf("controller.RemoveStatus: %w (%v)", ctrlErr.ErrParamsParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.RemoveStatus: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	err := c.service.RemoveStatus(req.StatusId)
	if err != nil {
		if errors.Is(err, repoStoreErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.SendStatus(204)
}
