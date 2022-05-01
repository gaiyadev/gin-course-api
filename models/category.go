package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name   string  `form:"name" json:"name" binding:"required,min=3" gorm:"unique"`
	UserID float64 `json:"user_id"`
	User   *User   `json:"user"`
	Post   *Post   `json:"post"`
}

type UpdateCategory struct {
	Name string `form:"name" json:"name" binding:"required,min=3" gorm:"unique"`
}
