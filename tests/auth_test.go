package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"todos/config"
	"todos/database"
	"todos/middleware"
	"todos/models"
	"todos/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var testRouter *gin.Engine

func setupTestRouter() {
	// Load test environment
	if err := godotenv.Load("../.env.test"); err != nil {
		// Use default test values if .env.test doesn't exist
		os.Setenv("APP_ENV", "test")
		os.Setenv("JWT_SECRET", "test_secret_key")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_USER", "postgres")
		os.Setenv("DB_PASSWORD", "yourpassword")
		os.Setenv("DB_NAME", "todos_test")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load test configuration:", err)
	}

	// Initialize JWT middleware
	middleware.InitJWT(cfg)

	// Set Gin test mode
	gin.SetMode(gin.TestMode)

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	// Clean up test data
	database.DB.Exec("DELETE FROM todos")
	database.DB.Exec("DELETE FROM users")

	// Setup router
	testRouter = gin.Default()
	routes.SetupRoutes(testRouter)
}

func TestRegister(t *testing.T) {
	setupTestRouter()

	t.Run("Success - Register new user", func(t *testing.T) {
		body := models.RegisterInput{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
		assert.NotEmpty(t, response["data"])
	})

	t.Run("Fail - Duplicate email", func(t *testing.T) {
		body := models.RegisterInput{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Fail - Missing required fields", func(t *testing.T) {
		body := map[string]string{"name": "John"}

		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLogin(t *testing.T) {
	setupTestRouter()

	// Create a user first
	user := models.User{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: "password123",
	}
	user.HashPassword()
	database.DB.Create(&user)

	t.Run("Success - Valid credentials", func(t *testing.T) {
		body := models.LoginInput{
			Email:    "jane@example.com",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		// Extract token for use in other tests
		data := response["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])
	})

	t.Run("Fail - Wrong password", func(t *testing.T) {
		body := models.LoginInput{
			Email:    "jane@example.com",
			Password: "wrongpassword",
		}

		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Fail - Non-existent email", func(t *testing.T) {
		body := models.LoginInput{
			Email:    "nobody@example.com",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
