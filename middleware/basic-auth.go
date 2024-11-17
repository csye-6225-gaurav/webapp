package middleware

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/csye-6225-gaurav/webapp/models"
	"github.com/csye-6225-gaurav/webapp/storage"
	"github.com/csye-6225-gaurav/webapp/utils"
	"github.com/gofiber/fiber/v2"
	zlog "github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func BasicAuthMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		authHeader := ctx.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			zlog.Error().
				Str("endpoint", ctx.Path()).
				Msg("Missing or invalid Authorization header")
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing or invalid Authorization header",
			})
		}

		encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")

		credentialsBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			zlog.Error().
				Err(err).
				Str("endpoint", ctx.Path()).
				Msg("Error decoding base64 credentials")
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid base64 encoding",
			})
		}

		credentials := strings.SplitN(string(credentialsBytes), ":", 2)
		if len(credentials) != 2 {
			zlog.Warn().
				Str("endpoint", ctx.Path()).
				Msg("Invalid credentials format")
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid credentials format",
			})
		}

		email := credentials[0]
		password := credentials[1]

		var user models.User
		start := time.Now()
		err = storage.DB.Where("email = ?", email).First(&user).Error
		utils.Client.PrecisionTiming("db.getUser", time.Since(start))
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message": "User not found"})
			}
			zlog.Error().
				Err(err).
				Str("endpoint", ctx.Path()).
				Str("email", email).
				Msg("Error retrieving user")
			return nil
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			zlog.Warn().
				Str("endpoint", ctx.Path()).
				Str("email", email).
				Msg("Invalid email or password")
			ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message": "Invalid email or password"})
			return nil
		}
		if !user.IsVerified {
			zlog.Warn().
				Str("endpoint", ctx.Path()).
				Str("email", email).
				Msg("user not verified")
			ctx.Status(fiber.StatusForbidden).JSON(&fiber.Map{"message": "user not verified"})
			return nil
		}
		ctx.Locals("user", user)
		zlog.Info().
			Str("endpoint", ctx.Path()).
			Str("email", email).
			Msg("User authenticated")
		return ctx.Next()
	}
}
