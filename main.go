package main

import (
	"fmt"
	"log"
	"os"

	"github.com/csye-6225-gaurav/webapp/models"
	"github.com/csye-6225-gaurav/webapp/repository"
	"github.com/csye-6225-gaurav/webapp/routes"
	"github.com/csye-6225-gaurav/webapp/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := storage.Config{
		Host:    os.Getenv("DB_Host"),
		Port:    os.Getenv("DB_Port"),
		Pass:    os.Getenv("DB_Pass"),
		User:    os.Getenv("DB_User"),
		DBname:  os.Getenv("DB_Name"),
		SSLMode: os.Getenv("DB_SSLMode"),
	}
	db, err := storage.NewConnection(&config)
	if err != nil {
		log.Println("Failed DB connection")
	}
	err = models.MigrateDB(db)
	if err != nil {
		log.Println(err)
	}
	tmp := repository.Repository{
		DB: db,
	}
	r := routes.Repo{
		Repository: &tmp,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	appPort := fmt.Sprintf(":%s", os.Getenv("APP_Port"))
	app.Listen(appPort)
}
