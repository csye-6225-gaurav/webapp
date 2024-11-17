package repository

import (
	"time"

	"github.com/csye-6225-gaurav/webapp/models"
	"github.com/csye-6225-gaurav/webapp/storage"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func VerifyUser(ctx *fiber.Ctx) error {
	email := ctx.Query("user")
	token := ctx.Query("token")
	if email == "" || token == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required query parameters: user and/or token",
		})
	}
	var user models.User
	if err := storage.DB.Where("email = ? AND token = ?", email, token).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found or invalid token",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	if time.Since(user.CreatedAt) > 2*time.Minute {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Verification link has expired",
		})
	}
	user.IsVerified = true
	if err := storage.DB.Save(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user verification status",
		})
	}
	return ctx.JSON(fiber.Map{
		"message": "User verified successfully",
	})
}
