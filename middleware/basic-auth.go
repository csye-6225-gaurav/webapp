package middleware

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/csye-6225-gaurav/webapp/models"
	"github.com/csye-6225-gaurav/webapp/storage"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func BasicAuthMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		authHeader := ctx.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			log.Println("Missing or invalid Authorization header")
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing or invalid Authorization header",
			})
		}

		encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")

		credentialsBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			log.Println("Error decoding base64 credentials:", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid base64 encoding",
			})
		}

		credentials := strings.SplitN(string(credentialsBytes), ":", 2)
		if len(credentials) != 2 {
			log.Println("Invalid credentials format")
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid credentials format",
			})
		}

		email := credentials[0]
		password := credentials[1]

		var user models.User
		err = storage.DB.Where("email = ?", email).First(&user).Error
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "User not found"})
			}
			log.Println("User not found:", err)
			return nil
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message": "Invalid email or password"})
			return nil
		}

		ctx.Locals("user", user)

		return ctx.Next()
	}
}
