package routes

import (
	"github.com/csye-6225-gaurav/webapp/repository"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("")
	api.All("/healthz", repository.HealthCheck)
	api.All("/healthz/*", func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusNotFound)
		return nil
	})
	v1 := api.Group("/v1")
	v1.Post("/user", repository.CreateUser)
	v1.Get("/user/self", repository.GetUser)
}
