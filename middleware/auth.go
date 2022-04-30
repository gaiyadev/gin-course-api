package middleware

import (
	"errors"
	"gin-course/config"
	"gin-course/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var JwtSecret = []byte(config.Config("JWT_SECRET"))

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&models.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		},
	)
	if err != nil {
		return nil
	}
	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		err = errors.New("couldn't parse claims")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message":    "No authorized token provided",
				"statusCode": http.StatusUnauthorized,
				"status":     "failed",
			})
			return
		}
		splitToken := strings.Split(tokenString, "Bearer ")
		tokenString = splitToken[1]
		err := ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message":    "Unauthorized access",
				"statusCode": http.StatusUnauthorized,
				"status":     "failed",
			})
			return
		}
		c.Next()
	}
}
