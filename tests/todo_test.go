package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"todos/database"
	"todos/models"

	"github.com/stretchr/testify/assert"
)

// Helper: login and get token
func getAuthToken(t *testing.T, email, password string) string {
	body := models.LoginInput{
		Email:    email,
		Password: password,
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	return data["token"].(string)
}

// Helper: create a test user and return token
func createTestUserAndLogin(t *testing.T) string {
	database.DB.Exec("DELETE FROM todos")
	database.DB.Exec("DELETE FROM users")

	user := models.User{
		Name:     "Test User",
		Email:    "testuser@example.com",
		Password: "password123",
	}
	user.HashPassword()
	database.DB.Create(&user)

	return getAuthToken(t, "testuser@example.com", "password123")
}

func TestCreateTodo(t *testing.T) {
	setupTestRouter()
	token := createTestUserAndLogin(t)

	t.Run("Success - Create todo", func(t *testing.T) {
		body := models.CreateTodoInput{
			Title:       "Test Todo",
			Description: "Test Description",
		}

		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
	})

	t.Run("Fail - Missing title", func(t *testing.T) {
		body := map[string]string{"description": "No title"}

		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Fail - No auth token", func(t *testing.T) {
		body := models.CreateTodoInput{Title: "Test"}

		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestGetTodos(t *testing.T) {
	setupTestRouter()
	token := createTestUserAndLogin(t)

	// Seed some todos
	for i := 1; i <= 3; i++ {
		body := models.CreateTodoInput{
			Title:       fmt.Sprintf("Todo %d", i),
			Description: fmt.Sprintf("Description %d", i),
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
	}

	t.Run("Success - Get all todos", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/todos", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
		assert.Equal(t, float64(3), response["count"])
	})

	t.Run("Success - Filter by completed=false", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/todos?completed=false", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUpdateTodo(t *testing.T) {
	setupTestRouter()
	token := createTestUserAndLogin(t)

	// Create a todo
	body := models.CreateTodoInput{Title: "Original Title"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	data := createResp["data"].(map[string]interface{})
	todoID := int(data["ID"].(float64))

	t.Run("Success - Update todo", func(t *testing.T) {
		updateBody := models.UpdateTodoInput{
			Title:     "Updated Title",
			Completed: true,
		}

		jsonBody, _ := json.Marshal(updateBody)
		req, _ := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("/api/v1/todos/%d", todoID),
			bytes.NewBuffer(jsonBody),
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Fail - Todo not found", func(t *testing.T) {
		updateBody := models.UpdateTodoInput{Title: "Updated"}

		jsonBody, _ := json.Marshal(updateBody)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/todos/9999", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeleteTodo(t *testing.T) {
	setupTestRouter()
	token := createTestUserAndLogin(t)

	// Create a todo to delete
	body := models.CreateTodoInput{Title: "Todo to Delete"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	data := createResp["data"].(map[string]interface{})
	todoID := int(data["ID"].(float64))

	t.Run("Success - Delete todo", func(t *testing.T) {
		req, _ := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/v1/todos/%d", todoID),
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Fail - Delete non-existent todo", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/todos/9999", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
