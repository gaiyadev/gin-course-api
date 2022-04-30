package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	//ID       uint   `gorm:"primaryKey"json:"id"`
	Name     string `form:"name" json:"name" binding:"required,min=3"`
	Email    string `form:"email" json:"email" binding:"required,email" gorm:"index:idx_email,unique"`
	Password string `form:"password" json:"password" binding:"required,min=6"`
	Posts    []Post `json:"posts"`
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
