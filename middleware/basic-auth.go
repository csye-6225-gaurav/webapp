package middleware

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func BasicAuthMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		authHeader := ctx.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
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
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid credentials format",
			})
		}

		email := credentials[0]
		password := credentials[1]

		ctx.Locals("email", email)
		ctx.Locals("password", password)

		return ctx.Next()
	}
}
