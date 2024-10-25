package repository

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/csye-6225-gaurav/webapp/models"
	"github.com/csye-6225-gaurav/webapp/storage"
	"github.com/gofiber/fiber/v2"
)

func SaveProfilePic(ctx *fiber.Ctx) error {

	user := ctx.Locals("user").(models.User)
	//TODO:check if profile pic already exists
	file, err := ctx.FormFile("profilePic")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse profile picture",
		})
	}
	fileContent, err := file.Open()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to open file",
		})
	}
	defer fileContent.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(fileContent); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to read file",
		})
	}
	bucketName := "csye6225-webapp"
	// Construct S3 file path
	fileName := fmt.Sprintf("%s/%s", user.ID.String(), file.Filename)

	// Upload to S3
	_, err = storage.S3Client.PutObject(ctx.Context(), &s3.PutObjectInput{
		Bucket: &bucketName, // replace with your bucket name
		Key:    &fileName,
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to upload file to S3",
		})
	}
	image := models.Image{
		FileName:   file.Filename,
		URL:        fmt.Sprintf("https://%s/%s", os.Getenv("S3_Bucket"), fileName),
		UploadDate: time.Now(),
		UserID:     user.ID,
	}
	if err := storage.DB.Create(&image).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save image metadata",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(image)

}

func GetProfilePic(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(models.User)
	var image models.Image
	err := storage.DB.Where("user_id = ?", user.ID).First(&image).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "Image not found"})
		}
		log.Println("Image not found:", err)
		return nil
	}

	ctx.Status(fiber.StatusOK).JSON(image)
	return nil
}
