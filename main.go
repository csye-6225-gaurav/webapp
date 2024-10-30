package main

import (
	"fmt"
	"os"
	"time"

	"github.com/csye-6225-gaurav/webapp/routes"
	"github.com/csye-6225-gaurav/webapp/storage"
	"github.com/csye-6225-gaurav/webapp/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	err := godotenv.Load(".env")
	if err != nil {
		zlog.Fatal().Msg(err.Error())
	}
	config := storage.Config{
		Host:    os.Getenv("DB_Host"),
		Port:    os.Getenv("DB_Port"),
		Pass:    os.Getenv("DB_Pass"),
		User:    os.Getenv("DB_User"),
		DBname:  os.Getenv("DB_Name"),
		SSLMode: os.Getenv("DB_SSLMode"),
	}
	err = storage.NewConnection(&config)
	if err != nil {
		zlog.Error().Msg(err.Error())
	}

	storage.ConnectToS3()
	utils.InitStatsD()
	app := fiber.New()
	routes.SetupRoutes(app)
	appPort := fmt.Sprintf(":%s", os.Getenv("APP_Port"))
	zlog.Info().Msg("application started successfully")
	app.Listen(appPort)
}
