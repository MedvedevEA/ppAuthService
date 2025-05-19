package middleware

import (
	"crypto/rsa"
	"fmt"
	"log/slog"
	"time"

	"skillsRockTodo/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func GetLoggerMiddlewareFunc(lg *slog.Logger, appName string) func(c *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()

		err := ctx.Next()

		lg.Info(
			fmt.Sprintf("API server '%s' request", appName),
			slog.Any("method", ctx.Method()),
			slog.Any("path", ctx.Path()),
			slog.Any("statusCode", ctx.Response().StatusCode()),
			slog.Any("time", time.Since(start)),
		)
		return err
	}
}
func GetAuthMiddlewareFunc(publicKey *rsa.PublicKey, tokenType string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tokenString := ctx.Get("Authorization")
		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			return ctx.Status(401).SendString("authMiddleware: the JWT string is not a valid token format")
		}
		tokenClaims, err := jwt.ParseToken(tokenString[7:], publicKey)
		if err != nil {
			return ctx.Status(401).SendString(fmt.Sprintf("authMiddleware: token parsing error(%v)", err))
		}
		if tokenClaims.TokenType != tokenType {
			return ctx.Status(401).SendString("authMiddleware: the token type is not of the expected type for this route")
		}

		ctx.Locals("claims", tokenClaims)
		return ctx.Next()
	}
}
func BadRequest(ctx *fiber.Ctx) error {
	return ctx.Status(404).SendString("badRequestMiddleware: unregistered route")
}
