package repository

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/csye-6225-gaurav/webapp/models"
	"github.com/csye-6225-gaurav/webapp/storage"
	"github.com/csye-6225-gaurav/webapp/utils"
	"github.com/gofiber/fiber/v2"
	zlog "github.com/rs/zerolog/log"
)

func SaveProfilePic(ctx *fiber.Ctx) error {

	user := ctx.Locals("user").(models.User)
	endpoint := ctx.Path()
	form, err := ctx.MultipartForm()
	if err != nil {
		zlog.Error().Err(err).Str("endpoint", endpoint).Msg("Failed to parse form data")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse form data",
		})
	}
	for key := range form.File {
		if key != "profilePic" {
			ctx.Status(fiber.StatusBadRequest)
			return nil
		}
	}
	if len(form.Value) > 0 {
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}
	var image models.Image
	err = storage.DB.Where("user_id = ?", user.ID).First(&image).Error
	if err == nil {
		zlog.Warn().Str("endpoint", endpoint).Str("user_id", user.ID.String()).Msg("Profile picture already exists")
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Profile picture already exists",
		})
	}
	file, err := ctx.FormFile("profilePic")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse profile picture",
		})
	}
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		zlog.Warn().Str("endpoint", endpoint).Str("file_extension", ext).Msg("Invalid file format")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Only JPG, JPEG, and PNG formats are allowed",
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
	bucketName := os.Getenv("Bucket_Name")
	// Construct S3 file path
	fileName := fmt.Sprintf("%s/%s", user.ID.String(), file.Filename)

	// Upload to S3
	start := time.Now()
	_, err = storage.S3Client.PutObject(ctx.Context(), &s3.PutObjectInput{
		Bucket: &bucketName, // replace with your bucket name
		Key:    &fileName,
		Body:   bytes.NewReader(buf.Bytes()),
	})
	utils.Client.PrecisionTiming("s3.put", time.Since(start))
	if err != nil {
		zlog.Error().Err(err).Str("endpoint", endpoint).Msg("Failed to upload file to S3")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to upload file to S3",
		})
	}
	zlog.Info().Str("endpoint", endpoint).Str("user_id", user.ID.String()).Msg("File uploaded successfully")
	image = models.Image{
		FileName:   file.Filename,
		URL:        fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, os.Getenv("region"), fileName),
		UploadDate: time.Now(),
		UserID:     user.ID,
	}
	if err := storage.DB.Create(&image).Error; err != nil {
		zlog.Error().Err(err).Str("endpoint", endpoint).Msg("Failed to save image metadata")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save image metadata",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(image)

}

func GetProfilePic(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(models.User)
	endpoint := ctx.Path()
	var image models.Image
	err := storage.DB.Where("user_id = ?", user.ID).First(&image).Error
	if err != nil {
		zlog.Warn().Err(err).Str("endpoint", endpoint).Str("user_id", user.ID.String()).Msg("Image not found")
		if strings.Contains(err.Error(), "record not found") {
			ctx.Status(fiber.StatusNotFound).JSON(&fiber.Map{"message": "Image not found"})
		}
		log.Println("Image not found:", err)
		return nil
	}
	zlog.Info().Str("endpoint", endpoint).Str("user_id", user.ID.String()).Msg("Profile picture fetched successfully")
	ctx.Status(fiber.StatusOK).JSON(image)
	return nil
}

func DeleteProfilePic(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(models.User)
	endpoint := ctx.Path()
	var image models.Image

	err := storage.DB.Where("user_id = ?", user.ID).First(&image).Error
	if err != nil {
		zlog.Warn().Err(err).Str("endpoint", endpoint).Str("user_id", user.ID.String()).Msg("Failed to fetch image metadata")
		if strings.Contains(err.Error(), "record not found") {
			ctx.Status(fiber.StatusNotFound)
			return nil
		}
		log.Println("Failed to fetch image metadata:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch image metadata"})
	}

	bucketName := os.Getenv("Bucket_Name")
	filePath := fmt.Sprintf("%s/%s", user.ID.String(), image.FileName)
	start := time.Now()
	_, err = storage.S3Client.DeleteObject(ctx.Context(), &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &filePath,
	})
	utils.Client.PrecisionTiming("s3.Delete", time.Since(start))
	if err != nil {
		zlog.Error().Err(err).Str("endpoint", endpoint).Str("user_id", user.ID.String()).Msg("Failed to delete image from S3")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to delete image from S3"})
	}
	zlog.Info().Str("endpoint", endpoint).Str("user_id", user.ID.String()).Msg("Image deleted from S3")
	if err := storage.DB.Delete(&image).Error; err != nil {
		zlog.Error().Err(err).Str("endpoint", endpoint).Str("user_id", user.ID.String()).Msg("Failed to delete image metadata")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to delete image metadata"})
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}
