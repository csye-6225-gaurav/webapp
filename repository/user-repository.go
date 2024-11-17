package repository

import (
	"context"
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/csye-6225-gaurav/webapp/models"
	"github.com/csye-6225-gaurav/webapp/storage"
	"github.com/csye-6225-gaurav/webapp/utils"
	"github.com/gofiber/fiber/v2"
	zlog "github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(ctx *fiber.Ctx) error {
	userReq := models.RequestUser{}
	if len(ctx.Queries()) > 0 {
		zlog.Error().Str("endpoint", ctx.Path()).Msg("Query parameters are not allowed for CreateUser endpoint")
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}
	j := json.NewDecoder(strings.NewReader(string(ctx.Body())))
	j.DisallowUnknownFields()
	err := j.Decode(&userReq)

	if err != nil {
		zlog.Error().Err(err).Str("endpoint", ctx.Path()).Msg("Error decoding JSON")
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid request body"})
		return nil
	}
	if userReq.Email == "" || !isValidEmail(userReq.Email) {
		zlog.Error().Str("endpoint", ctx.Path()).Msg("Invalid or missing email")
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid or missing email"})
		return nil
	}
	if userReq.FirstName == "" {
		zlog.Error().Str("endpoint", ctx.Path()).Msg("FirstName is required")
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "FirstName required"})
		return nil
	}
	if userReq.LastName == "" {
		zlog.Error().Str("endpoint", ctx.Path()).Msg("LastName is required")
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "LastName required"})
		return nil
	}
	if userReq.Password == "" {
		zlog.Error().Str("endpoint", ctx.Path()).Msg("Password is required")
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Password is required"})
		return nil
	}
	if len(userReq.Password) < 8 {
		zlog.Error().Str("endpoint", ctx.Path()).Msg("Password should be more than 8 characters")
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Password should be more than 8 characters"})
		return nil
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		zlog.Error().Err(err).Str("endpoint", ctx.Path()).Msg("Error hashing password")
		ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "Error while hashing password"})
		return nil
	}
	user := models.User{
		Email:     userReq.Email,
		Password:  string(hashedPassword),
		FirstName: userReq.FirstName,
		LastName:  userReq.LastName,
	}
	start := time.Now()
	err = storage.DB.Create(&user).Error
	utils.Client.PrecisionTiming("db.createUser", time.Since(start))
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"uni_users_email\" ") {
			zlog.Error().Err(err).Str("endpoint", ctx.Path()).Str("email", user.Email).Msg("User already exists with the same email")
			ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "User already exists with the same email"})
			return nil
		}
		zlog.Error().Err(err).Str("endpoint", ctx.Path()).Msg("Failed to create user in database")
		ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "could not create user"})
		return nil
	}
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(os.Getenv("region")),
	)
	if err != nil {
		zlog.Error().Err(err).Msg("Faild to load aws config")
	}
	msg := models.Message{
		Email: user.Email,
		Token: user.Token,
	}
	message, err := json.Marshal(msg)
	snsClient := sns.NewFromConfig(cfg)
	publishInput := sns.PublishInput{TopicArn: aws.String(os.Getenv("sns_topic")), Message: aws.String(string(message))}
	_, err = snsClient.Publish(context.TODO(), &publishInput)
	if err != nil {
		zlog.Error().Err(err).Msg("Faild publish message")
	}
	resUser := models.ResponseUser{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	zlog.Info().Str("endpoint", ctx.Path()).Str("user_id", user.ID.String()).Msg("User created successfully")
	ctx.Status(fiber.StatusCreated).JSON(resUser)
	return nil
}

func isValidEmail(email string) bool {
	// Regular expression for validating email format
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func GetUser(ctx *fiber.Ctx) error {

	user := ctx.Locals("user").(models.User)
	if len(ctx.Queries()) > 0 {
		zlog.Error().Str("endpoint", ctx.Path()).Msg("Query parameters are not allowed for GetUser endpoint")
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}
	if len(ctx.Body()) > 0 {
		zlog.Error().Str("endpoint", ctx.Path()).Msg("Request body is not allowed for GetUser endpoint")
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}

	resUser := models.ResponseUser{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	zlog.Info().Str("endpoint", ctx.Path()).Str("user_id", user.ID.String()).Msg("User retrieved successfully")
	ctx.Status(fiber.StatusOK).JSON(resUser)
	return nil
}

func UpdateUser(ctx *fiber.Ctx) error {

	user := ctx.Locals("user").(models.User)
	if len(ctx.Queries()) > 0 {
		zlog.Error().Str("endpoint", ctx.Path()).Msg("Query parameters are not allowed for UpdateUser endpoint")
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}
	// var user models.User
	var updateUser models.UpdateUser
	j := json.NewDecoder(strings.NewReader(string(ctx.Body())))
	j.DisallowUnknownFields()
	err := j.Decode(&updateUser)
	if err != nil {
		zlog.Error().Err(err).Str("endpoint", ctx.Path()).Msg("Error decoding JSON in UpdateUser")
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid request body"})
		return nil
	}

	if updateUser.Password != "" {
		if len(updateUser.Password) < 8 {
			zlog.Error().Str("endpoint", ctx.Path()).Msg("Password should be more than 8 characters")
			ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Password should be more than 8 characters"})
			return nil
		}
		if updateUser.Password == "" {
			ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Password can't be empty"})
			return nil
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), bcrypt.DefaultCost)
		if err != nil {
			zlog.Error().Err(err).Str("endpoint", ctx.Path()).Msg("Error hashing password in UpdateUser")
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
	start := time.Now()
	err = storage.DB.Save(&user).Error
	utils.Client.PrecisionTiming("db.updateUser", time.Since(start))
	if err != nil {
		zlog.Error().Err(err).Str("endpoint", ctx.Path()).Msg("Error updating user in the database")
		ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "Error updating user"})
		return nil
	}
	zlog.Info().Str("endpoint", ctx.Path()).Str("user_id", user.ID.String()).Msg("User updated successfully")
	ctx.Status(fiber.StatusNoContent)
	return nil
}
