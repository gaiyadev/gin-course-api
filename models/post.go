package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title      string    `form:"title" json:"title" binding:"required,min=3"`
	Body       string    `form:"body" json:"body" binding:"required"`
	UserID     float64   `json:"user_id"`
	User       *User     `json:"user"`
	CategoryID float64   `json:"category_id" binding:"required"`
	Category   *Category `json:"category"`
}

type UpdatePost struct {
	Title string `form:"title" json:"title" binding:"required,min=3"`
	Body  string `form:"body" json:"body" binding:"required"`
}
