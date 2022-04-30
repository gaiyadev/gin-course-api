package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	//ID     uint    `gorm:"primaryKey"json:"id"`
	Title  string  `form:"title" json:"title" binding:"required,min=3"`
	Body   string  `form:"body" json:"body" binding:"required"`
	UserID float64 `json:"user_id"`
}

type UpdatePost struct {
	Title string `form:"title" json:"title" binding:"required,min=3"`
	Body  string `form:"body" json:"body" binding:"required"`
}
