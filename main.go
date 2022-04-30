package main

import (
	"fmt"
	"gin-course/database"
	"gin-course/middleware"
	"gin-course/routes"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	// API Routes
	{
		v1.GET("/", routes.Index)
		v1.POST("/accounts/signup", routes.SignUp)
		v1.POST("/accounts/signin", routes.SignIn)
		v1.GET("/posts", routes.FetchPosts)
		v1.GET("/posts/:postId", routes.FetchPost)
	}

	secured := v1.Group("/").Use(middleware.AuthMiddleware())
	{
		secured.POST("posts/", routes.CreatePost)
		secured.DELETE("posts/:postId", routes.DeletePost)
		secured.PUT("posts/:postId", routes.UpdatePost)
		secured.GET("posts/user", routes.FetchUserPosts)

	}
	// DB Connection
	database.DBConnection()
	database.AutoMigrate()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	// listen and serve on 0.0.0.0:8080
	err := router.Run(":" + port)
	if err != nil {
		fmt.Print("Server not starting...")
	}
}
