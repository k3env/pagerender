package handlers

import "github.com/gofiber/fiber/v3"

func HealthHandler() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		return ctx.SendString("OK")
	}
}
