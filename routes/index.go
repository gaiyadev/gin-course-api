package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"status":     "Success",
		"message":    "Welcome to Gin!",
		"developer":  "gaiyadev",
	})

}
