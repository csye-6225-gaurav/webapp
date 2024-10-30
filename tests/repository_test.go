package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/csye-6225-gaurav/webapp/models"
	"github.com/csye-6225-gaurav/webapp/routes"
	"github.com/csye-6225-gaurav/webapp/storage"
	"github.com/csye-6225-gaurav/webapp/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestMain sets up and tears down the test environment
func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load("../.env")
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
	// Set up test database connection
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Pass, config.DBname, config.Port, config.SSLMode,
	)

	storage.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Could not connect to the test database: %v\n", err)
		os.Exit(1)
	}
	storage.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	storage.DB.AutoMigrate(&models.User{})
	utils.InitStatsD()
	// Run tests
	code := m.Run()

	// Clean up test database
	storage.DB.Exec("DROP TABLE users;")
	// Exit with the code from the test run
	os.Exit(code)
}

func setupApp() *fiber.App {
	app := fiber.New()
	routes.SetupRoutes(app)
	return app
}

// Test for creating a new user
func TestCreateUser(t *testing.T) {
	app := setupApp()

	payload := models.RequestUser{
		Email:     "testuser@example.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password123",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdUser models.User
	err = storage.DB.Where("email = ?", "testuser@example.com").First(&createdUser).Error
	assert.NoError(t, err)

	assert.Equal(t, "testuser@example.com", createdUser.Email)
	assert.Equal(t, "Test", createdUser.FirstName)
	assert.Equal(t, "User", createdUser.LastName)
}

// Test for fetching a user's info with basic authentication
func TestGetUser(t *testing.T) {
	app := setupApp()

	user := models.User{
		Email:     "testuser@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/user/self", nil)
	req.SetBasicAuth(user.Email, "password123") // Mocking the Basic Auth
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fetchedUser models.User
	err = json.NewDecoder(resp.Body).Decode(&fetchedUser)
	assert.NoError(t, err)

	assert.Equal(t, user.Email, fetchedUser.Email)
	assert.Equal(t, user.FirstName, fetchedUser.FirstName)
	assert.Equal(t, user.LastName, fetchedUser.LastName)
}

// Test for updating a user's info
func TestUpdateUser(t *testing.T) {
	app := setupApp()

	updatePayload := models.UpdateUser{
		FirstName: "UpdatedName",
		LastName:  "UpdatedLastName",
		Password:  "newpassword123",
	}
	body, _ := json.Marshal(updatePayload)

	req := httptest.NewRequest(http.MethodPut, "/v1/user/self", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("testuser@example.com", "password123")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	var updatedUser models.User
	storage.DB.Where("email = ?", "testuser@example.com").First(&updatedUser)

	assert.Equal(t, updatePayload.FirstName, updatedUser.FirstName)
	assert.Equal(t, updatePayload.LastName, updatedUser.LastName)
}
