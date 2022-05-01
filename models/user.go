package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName  string     `form:"first_name" json:"first_name" binding:"required,min=3"`
	LastName   string     `form:"last_name" json:"last_name" binding:"required,min=3"`
	Username   string     `form:"username" json:"username" binding:"required,min=3"`
	Email      string     `form:"email" json:"email" binding:"required,email" gorm:"index:idx_email,unique"`
	Password   string     `form:"password" json:"password" binding:"required,min=6"`
	Posts      []Post     `json:"posts" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Categories []Category `json:"categories" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Login struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=6"`
}

type Claims struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	ID        uint   `json:"id"`
	jwt.StandardClaims
}

type ChangePassword struct {
	CurrentPassword   string `form:"current_password" json:"current_password" binding:"required,min=6"`
	NewPassword       string `form:"new_password" json:"new_password" binding:"required,min=6"`
	ConfirmedPassword string `form:"confirmed_password" json:"confirmed_password" binding:"required,min=6"`
}

type UpdateAccount struct {
	FirstName string `form:"first_name" json:"first_name" binding:"required,min=3"`
	LastName  string `form:"last_name" json:"last_name" binding:"required,min=3"`
	Username  string `form:"username" json:"username" binding:"required,min=3"`
}
