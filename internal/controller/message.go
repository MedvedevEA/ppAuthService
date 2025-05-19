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

func (c *Controller) AddMessage(ctx *fiber.Ctx) error {
	req := new(ctrlDto.AddMessage)
	if err := ctx.BodyParser(req); err != nil {
		err := fmt.Errorf("controller.AddMessage: %w (%v)", ctrlErr.ErrBodyParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.AddMessage: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	message, err := c.service.AddMessage(&repoStoreDto.AddMessage{
		TaskId: req.TaskId,
		UserId: req.UserId,
		Text:   req.Text,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(201).JSON(message)
}
func (c *Controller) GetMessage(ctx *fiber.Ctx) error {
	req := new(ctrlDto.GetMessage)
	if err := ctx.ParamsParser(req); err != nil {
		err := fmt.Errorf("controller.GetMessage: %w (%v)", ctrlErr.ErrParamsParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.GetMessage: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	message, err := c.service.GetMessage(req.MessageId)
	if err != nil {
		if errors.Is(err, repoStoreErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(message)
}
func (c *Controller) GetMessages(ctx *fiber.Ctx) error {
	req := &ctrlDto.GetMessages{
		Offset: 0,
		Limit:  10,
	}
	if err := ctx.QueryParser(req); err != nil {
		err := fmt.Errorf("controller.GetMessages: %w (%v)", ctrlErr.ErrQueryParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.GetMessages: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	Messages, err := c.service.GetMessages(&repoStoreDto.GetMessages{
		Offset: req.Offset,
		Limit:  req.Limit,
		TaskId: req.TaskId,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(Messages)
}
func (c *Controller) UpdateMessage(ctx *fiber.Ctx) error {
	req := new(ctrlDto.UpdateMessage)
	if err := ctx.BodyParser(req); err != nil {
		err := fmt.Errorf("controller.UpdateMessage: %w (%v)", ctrlErr.ErrBodyParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := ctx.ParamsParser(req); err != nil {
		err := fmt.Errorf("controller.UpdateMessage: %w (%v)", ctrlErr.ErrParamsParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.UpdateMessage: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	message, err := c.service.UpdateMessage(&repoStoreDto.UpdateMessage{
		MessageId: req.MessageId,
		Text:      req.Text,
	})
	if err != nil {
		if errors.Is(err, repoStoreErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(message)
}
func (c *Controller) RemoveMessage(ctx *fiber.Ctx) error {
	req := new(ctrlDto.RemoveMessage)
	if err := ctx.ParamsParser(req); err != nil {
		err := fmt.Errorf("controller.RemoveMessage: %w (%v)", ctrlErr.ErrParamsParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.RemoveMessage: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	err := c.service.RemoveMessage(req.MessageId)
	if err != nil {
		if errors.Is(err, repoStoreErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.SendStatus(204)
}
