package handlers

import (
	"net/http"
	"todos/database"
	"todos/middleware"
	"todos/models"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.RegisterInput  true  "Register Input"
// @Success      201    {object}  models.AuthResponse
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /auth/register [post]
func Register(c *gin.Context) {
	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Email already registered",
		})
		return
	}

	// Create user
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to process password",
		})
		return
	}

	// Save to database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create user",
		})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully",
		"data": models.AuthResponse{
			Token: token,
			User: models.UserResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			},
		},
	})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.LoginInput  true  "Login Input"
// @Success      200    {object}  models.AuthResponse
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var input models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid email or password",
		})
		return
	}

	// Check password
	if !user.CheckPassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid email or password",
		})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data": models.AuthResponse{
			Token: token,
			User: models.UserResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			},
		},
	})
}

// GetProfile godoc
// @Summary      Get current user profile
// @Description  Get the authenticated user's profile
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.UserResponse
// @Failure      401  {object}  map[string]interface{}
// @Router       /auth/profile [get]
func GetProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}
