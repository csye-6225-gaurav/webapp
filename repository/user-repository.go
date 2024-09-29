package repository

import (
	"encoding/json"
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
	if ctx.Method() != fiber.MethodPost {
		log.Println("Method not allowed")
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	}
	j := json.NewDecoder(strings.NewReader(string(ctx.Body())))
	j.DisallowUnknownFields()
	err := j.Decode(&user)

	if err != nil {
		log.Println("Error decoding JSON:", err)
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid request body"})
		return nil
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

	email := ctx.Locals("email").(string)
	password := ctx.Locals("password").(string)

	var user models.User
	err := storage.DB.Where("email = ?", email).First(&user).Error
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

func UpdateUser(ctx *fiber.Ctx) error {

	email := ctx.Locals("email").(string)
	password := ctx.Locals("password").(string)

	var user models.User
	err := storage.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		ctx.Status(fiber.StatusNotFound).JSON(&fiber.Map{"message": "User not found"})
		return nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message": "Invalid credentials"})
		return nil
	}

	updatedUser := models.User{}
	err = ctx.BodyParser(&updatedUser)
	if err != nil {
		ctx.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Invalid request body"})
		return err
	}

	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "Error while hashing password"})
			return nil
		}
		user.Password = string(hashedPassword)
	}

	if updatedUser.FirstName != "" {
		user.FirstName = updatedUser.FirstName
	}
	if updatedUser.LastName != "" {
		user.LastName = updatedUser.LastName
	}

	err = storage.DB.Save(&user).Error
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "Error updating user"})
		return nil
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}
