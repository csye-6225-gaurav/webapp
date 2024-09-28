package repository

import (
	"encoding/base64"
	"log"
	"regexp"
	"strings"

	"github.com/csye-6225-gaurav/webapp/models"
	"github.com/csye-6225-gaurav/webapp/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(ctx *fiber.Ctx) error {
	user := models.User{}
	user.ID = uuid.New()
	err := ctx.BodyParser(&user)

	if err != nil {
		ctx.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}
	if user.Email == "" || !isValidEmail(user.Email) {
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid or missing email"})
		return nil
	}
	if user.Password == "" {
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Password is required"})
		return nil
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "Error while hashing password"})
		return nil
	}
	user.Password = string(hashedPassword)
	err = storage.DB.Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"uni_users_email\" ") {
			ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "User already exists with the same email"})
			return nil
		}
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "could not create user"})
		return nil
	}

	ctx.Status(fiber.StatusCreated)
	return nil
}

func isValidEmail(email string) bool {
	// Regular expression for validating email format
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func GetUser(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
		ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message": "Missing or invalid Authorization header"})
		return nil
	}

	encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
	credentialsBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		log.Println("Error decoding base64 credentials:", err)
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid base64 encoding"})
		return nil
	}

	credentials := strings.SplitN(string(credentialsBytes), ":", 2)
	if len(credentials) != 2 {
		ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message": "Invalid credentials format"})
		return nil
	}
	email := credentials[0]
	password := credentials[1]

	var user models.User
	err = storage.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Println("User not found:", err)
		ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message": "Invalid email or password"})
		return nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message": "Invalid email or password"})
		return nil
	}

	user.Password = ""

	ctx.Status(fiber.StatusOK).JSON(user)
	return nil
}
