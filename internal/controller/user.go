package controller

import (
	"fmt"
	ctrlDto "ppAuthService/internal/controller/dto"
	ctrlErr "ppAuthService/internal/controller/err"
	repoStoreDto "ppAuthService/internal/repository/repostore/dto"
	"ppAuthService/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func (c *Controller) GetUsers(ctx *fiber.Ctx) error {
	req := &ctrlDto.GetUsers{
		Offset: 0,
		Limit:  10,
	}
	if err := ctx.QueryParser(req); err != nil {
		err := fmt.Errorf("controller.GetUsers: %w (%v)", ctrlErr.ErrQueryParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		err := fmt.Errorf("controller.GetUsers: %w (%v)", ctrlErr.ErrValidate, err)
		return ctx.Status(400).SendString(err.Error())
	}
	users, err := c.service.GetUsers(&repoStoreDto.GetUsers{
		Offset: req.Offset,
		Limit:  req.Limit,
		Name:   req.Name,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(users)
}
