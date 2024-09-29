package routes

import (
	"github.com/csye-6225-gaurav/webapp/middleware"
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
	v1.All("/user", repository.CreateUser)
	v1.Get("/user/self", middleware.BasicAuthMiddleware(), repository.GetUser)
	v1.Put("/user/self", middleware.BasicAuthMiddleware(), repository.UpdateUser)
	v1.Post("/user/self", func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	})
	v1.Patch("/user/self", func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	})
	v1.Head("/user/self", func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	})
	v1.Delete("/user/self", func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	})
	v1.Options("/user/self", func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	})
}
