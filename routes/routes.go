package routes

import (
	"github.com/csye-6225-gaurav/webapp/repository"
	"github.com/gofiber/fiber/v2"
)

type Repo struct {
	*repository.Repository
}

func (r *Repo) SetupRoutes(app *fiber.App) {
	api := app.Group("")
	api.All("/healthz", r.HealthCheck)
	api.All("/healthz/*", func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusNotFound)
		return nil
	})
}
