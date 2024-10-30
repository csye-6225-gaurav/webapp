package middleware

import (
	"time"

	"github.com/csye-6225-gaurav/webapp/utils"
	"github.com/gofiber/fiber/v2"
)

func CountMetric() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()

		path := ctx.Path()
		utils.CountIncrement(path)
		err := ctx.Next()

		utils.CountTimer(path, start)

		return err
	}
}
