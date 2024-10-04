package repository

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/csye-6225-gaurav/webapp/models"
	"github.com/csye-6225-gaurav/webapp/storage"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(ctx *fiber.Ctx) error {
	userReq := models.RequestUser{}
	if len(ctx.Queries()) > 0 {
		log.Println("Query parameters are not allowed for create user")
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}
	j := json.NewDecoder(strings.NewReader(string(ctx.Body())))
	j.DisallowUnknownFields()
	err := j.Decode(&userReq)

	if err != nil {
		log.Println("Error decoding JSON:", err)
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid request body"})
		return nil
	}
	if userReq.Email == "" || !isValidEmail(userReq.Email) {
		log.Println("Invalid or missing email")
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid or missing email"})
		return nil
	}
	if userReq.FirstName == "" {
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "FirstName required"})
		return nil
	}
	if userReq.LastName == "" {
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "LastName required"})
		return nil
	}
	if userReq.Password == "" {
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Password is required"})
		return nil
	}
	if len(userReq.Password) < 8 {
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Password should be more than 8 characters"})
		return nil
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "Error while hashing password"})
		return nil
	}
	user := models.User{
		Email:     userReq.Email,
		Password:  string(hashedPassword),
		FirstName: userReq.FirstName,
		LastName:  userReq.LastName,
	}
	err = storage.DB.Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"uni_users_email\" ") {
			ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "User already exists with the same email"})
			return nil
		}
		ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "could not create user"})
		return nil
	}
	user.Password = ""
	ctx.Status(fiber.StatusCreated).JSON(user)
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
	if len(ctx.Queries()) > 0 {
		log.Println("Query parameters are not allowed for get user")
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}
	if len(ctx.Body()) > 0 {
		log.Println("Request body is not allowed for get user endpoint")
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}
	var user models.User
	err := storage.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "could not fetch user"})
		}
		log.Println("User not found:", err)
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
	if len(ctx.Queries()) > 0 {
		log.Println("Query parameters are not allowed for update user")
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}
	var user models.User
	var updateUser models.UpdateUser
	j := json.NewDecoder(strings.NewReader(string(ctx.Body())))
	j.DisallowUnknownFields()
	err := j.Decode(&updateUser)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid request body"})
		return nil
	}
	err = storage.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "could not fetch user"})
			return nil
		}
		ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message": "User not found"})
		return nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message": "Invalid credentials"})
		return nil
	}
	if updateUser.Password != "" {
		if len(updateUser.Password) < 8 {
			ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Password should be more than 8 characters"})
			return nil
		}
		if updateUser.Password == "" {
			ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Password can't be empty"})
			return nil
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "Error while hashing password"})
			return nil
		}
		user.Password = string(hashedPassword)
	}

	if updateUser.FirstName != "" {
		user.FirstName = updateUser.FirstName
	}
	if updateUser.LastName != "" {
		user.LastName = updateUser.LastName
	}

	err = storage.DB.Save(&user).Error
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "Error updating user"})
		return nil
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}
