package models

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Completed   bool   `json:"completed" gorm:"default:false"`
	UserID      uint   `json:"user_id"`   // Link todo to user
}

type CreateTodoInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateTodoInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
