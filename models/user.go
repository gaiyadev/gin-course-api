package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string     `form:"name" json:"name" binding:"required,min=3"`
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
	Name  string `json:"name"`
	Email string `json:"email"`
	ID    uint   `json:"id"`
	jwt.StandardClaims
}

type ChangePassword struct {
	CurrentPassword   string `form:"current_password" json:"current_password" binding:"required,min=6"`
	NewPassword       string `form:"new_password" json:"new_password" binding:"required,min=6"`
	ConfirmedPassword string `form:"confirmed_password" json:"confirmed_password" binding:"required,min=6"`
}
