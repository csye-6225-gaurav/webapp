package repository

import (
	"log"
	"strings"

	"github.com/csye-6225-gaurav/webapp/models"
	"github.com/csye-6225-gaurav/webapp/storage"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(ctx *fiber.Ctx) error {
	user := models.User{}

	err := ctx.BodyParser(&user)

	if err != nil {
		ctx.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
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
			ctx.Status(fiber.StatusConflict).JSON(&fiber.Map{"message": "User already exists with the same email"})
			return nil
		}
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "could not create user"})
		return nil
	}

	ctx.Status(fiber.StatusCreated)
	return nil
}
