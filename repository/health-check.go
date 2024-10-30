package repository

import (
	"time"

	"github.com/csye-6225-gaurav/webapp/storage"
	"github.com/csye-6225-gaurav/webapp/utils"
	"github.com/gofiber/fiber/v2"
	zlog "github.com/rs/zerolog/log"
)

func HealthCheck(ctx *fiber.Ctx) error {

	ctx.Set("cache-control", "no-cache")
	if ctx.Method() != fiber.MethodGet {
		zlog.Warn().Str("endpoint", ctx.Path()).Msg("Method not allowed")
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	}
	if len(ctx.Body()) > 0 {
		zlog.Warn().Str("endpoint", ctx.Path()).Msg("Request body is not allowed for health check")
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}
	if len(ctx.Queries()) > 0 {
		zlog.Warn().Str("endpoint", ctx.Path()).Msg("Query parameters are not allowed for health check")
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}
	sqlDB, err := storage.DB.DB()
	if err != nil {
		zlog.Error().Str("endpoint", ctx.Path()).Msg("Failed to get database instance")
	}
	start := time.Now()
	if err := sqlDB.Ping(); err != nil {
		zlog.Error().Err(err).Str("endpoint", ctx.Path()).Msg("Failed to ping database")
		ctx.Status(fiber.StatusServiceUnavailable)
		return nil
	}
	utils.Client.PrecisionTiming("db.ping", time.Since(start))
	zlog.Info().Str("endpoint", ctx.Path()).Msg("Successfully pinged the PostgreSQL database")
	ctx.Status(fiber.StatusOK)
	return nil
}
