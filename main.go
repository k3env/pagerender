package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recover2 "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/k3env/pagerender/handlers"
)

func main() {
	app := fiber.New(fiber.Config{})

	handlers.InitMetrics()

	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover2.New())
	app.Use(handlers.MetricsMiddleware())

	app.Post("/render", handlers.RenderHandler())
	app.Get("/health", handlers.HealthHandler())
	app.Get("/metrics", handlers.MetricsHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Listening on http://localhost:" + port)
	log.Fatal(app.Listen(":"+port, fiber.ListenConfig{DisableStartupMessage: true}))
}
