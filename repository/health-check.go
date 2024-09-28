package repository

import (
	"log"

	"github.com/csye-6225-gaurav/webapp/storage"
	"github.com/gofiber/fiber/v2"
)

func HealthCheck(ctx *fiber.Ctx) error {

	ctx.Set("cache-control", "no-cache")
	if ctx.Method() != fiber.MethodGet {
		log.Println("Method not allowed")
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	}
	if len(ctx.Body()) > 0 {
		log.Println("Request body is not allowed for health check")
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}
	if len(ctx.Queries()) > 0 {
		log.Println("Query parameters are not allowed for health check")
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}
	sqlDB, err := storage.DB.DB()
	if err != nil {
		log.Println("failed to get DB instance")
	}
	if err := sqlDB.Ping(); err != nil {
		log.Printf("failed to ping database: %v", err)
		ctx.Status(fiber.StatusServiceUnavailable)
		return nil
	}
	log.Println("Successfully pinged the PostgreSQL database!")
	ctx.Status(fiber.StatusOK)
	return nil
}
