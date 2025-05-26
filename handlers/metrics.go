package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func MetricsHandler() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		return ctx.Status(http.StatusInternalServerError).SendString("Not implemented")
	}
}
