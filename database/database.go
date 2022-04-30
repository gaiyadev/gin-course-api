package database

import (
	"fmt"
	"gin-course/config"
	"gin-course/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB
var dsn = config.Config("DATABASE_URL")

func DBConnection() {
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database")
	}
	fmt.Println("Rest API Connected to database successfully")
	ok := database.AutoMigrate(&models.Post{}, &models.User{})
	if ok != nil {
		fmt.Println("Database migration not successful")
	}
	fmt.Println("Database migration successful")
	DB = database
}
