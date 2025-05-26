package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/k3env/pagerender/handlers"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	app.Post("/render", handlers.RenderHandler())
	app.Get("/health", handlers.HealthHandler())
	app.Get("/metrics", handlers.MetricsHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Listening on http://localhost:" + port)
	log.Fatal(app.Listen(":" + port))
}
