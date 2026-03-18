package handlers

import (
	"net/http"
	"todos/database"
	"todos/models"

	"github.com/gin-gonic/gin"
)

// GetTodos godoc
// @Summary      Get all todos
// @Description  Get all todos for the authenticated user
// @Tags         Todos
// @Produce      json
// @Security     BearerAuth
// @Param        completed  query     bool  false  "Filter by completed status"
// @Success      200        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
// @Router       /todos [get]
func GetTodos(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var todos []models.Todo

	completed := c.Query("completed")
	query := database.DB.Where("user_id = ?", userID)

	switch completed {
	case "true":
		query = query.Where("completed = ?", true)
	case "false":
		query = query.Where("completed = ?", false)
	}

	if err := query.Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch todos",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(todos),
		"data":    todos,
	})
}

// GetTodo godoc
// @Summary      Get a single todo
// @Description  Get a todo by its ID
// @Tags         Todos
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Todo ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /todos/{id} [get]
func GetTodo(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var todo models.Todo

	if err := database.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).
		First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    todo,
	})
}

// CreateTodo godoc
// @Summary      Create a new todo
// @Description  Create a new todo for the authenticated user
// @Tags         Todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      models.CreateTodoInput  true  "Todo Input"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Router       /todos [post]
func CreateTodo(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var input models.CreateTodoInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	todo := models.Todo{
		Title:       input.Title,
		Description: input.Description,
		Completed:   false,
		UserID:      userID,
	}

	if err := database.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create todo",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Todo created successfully",
		"data":    todo,
	})
}

// UpdateTodo godoc
// @Summary      Update a todo
// @Description  Update a todo by its ID
// @Tags         Todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int                     true  "Todo ID"
// @Param        input  body      models.UpdateTodoInput  true  "Update Input"
// @Success      200    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Router       /todos/{id} [put]
func UpdateTodo(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var todo models.Todo

	if err := database.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).
		First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Todo not found",
		})
		return
	}

	var input models.UpdateTodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := database.DB.Model(&todo).Updates(map[string]interface{}{
		"title":       input.Title,
		"description": input.Description,
		"completed":   input.Completed,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update todo",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Todo updated successfully",
		"data":    todo,
	})
}

// DeleteTodo godoc
// @Summary      Delete a todo
// @Description  Delete a todo by its ID
// @Tags         Todos
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Todo ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /todos/{id} [delete]
func DeleteTodo(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var todo models.Todo

	if err := database.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).
		First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Todo not found",
		})
		return
	}

	if err := database.DB.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete todo",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Todo deleted successfully",
	})
}
