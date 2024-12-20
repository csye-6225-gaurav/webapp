package routes

import (
	"github.com/csye-6225-gaurav/webapp/middleware"
	"github.com/csye-6225-gaurav/webapp/repository"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("")
	api.Use(func(ctx *fiber.Ctx) error {
		if ctx.Path() == "/healthz" && ctx.Method() != fiber.MethodGet {
			ctx.Set("cache-control", "no-cache")
			ctx.Status(fiber.StatusMethodNotAllowed)
			return nil
		}
		if ctx.Path() == "/v1/user" && ctx.Method() != fiber.MethodPost {
			ctx.Set("cache-control", "no-cache")
			ctx.Status(fiber.StatusMethodNotAllowed)
			return nil
		}
		if ctx.Path() == "/v1/user/self" {
			ctx.Set("cache-control", "no-cache")
		}
		return ctx.Next()
	})
	api.Get("/healthz", middleware.CountMetric(), repository.HealthCheck)
	api.All("/healthz/*", middleware.CountMetric(), func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusNotFound)
		return nil
	})
	v1 := api.Group("/v1")
	v1.Post("/user", middleware.CountMetric(), repository.CreateUser)
	v1.Get("/user/self", middleware.CountMetric(), middleware.BasicAuthMiddleware(), repository.GetUser)
	v1.Put("/user/self", middleware.CountMetric(), middleware.BasicAuthMiddleware(), repository.UpdateUser)
	v1.Post("/user/self", middleware.CountMetric(), func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	})
	v1.Patch("/user/self", middleware.CountMetric(), func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	})
	v1.Head("/user/self", middleware.CountMetric(), func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	})
	v1.Delete("/user/self", middleware.CountMetric(), func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	})
	v1.Options("/user/self", middleware.CountMetric(), func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	})
	v1.Post("/user/self/pic", middleware.CountMetric(), middleware.BasicAuthMiddleware(), repository.SaveProfilePic)
	v1.Get("/user/self/pic", middleware.CountMetric(), middleware.BasicAuthMiddleware(), repository.GetProfilePic)
	v1.Delete("/user/self/pic", middleware.CountMetric(), middleware.BasicAuthMiddleware(), repository.DeleteProfilePic)
	v1.Get("/user/self/verify", middleware.CountMetric(), repository.VerifyUser)
}
