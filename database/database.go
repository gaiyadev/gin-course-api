package database

import (
	"fmt"
	"gin-course/config"
	"gin-course/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dsn = config.Config("DATABASE_URL")

func DBConnection() {
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database")
	}
	fmt.Println("Rest API Connected to database successfully")
	DB = database
}

func AutoMigrate() {
	err := DB.AutoMigrate(&models.Post{}, &models.User{})
	if err != nil {
		fmt.Println("Database migration not successful")
	}
	fmt.Println("Database migration successful")
}
