package helpers

import (
	"fmt"
	"gin-course/config"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

func ExtractClaims(accessToken string) map[string]interface{} {
	var JwtSecret = []byte(config.Config("JWT_SECRET"))
	splitToken := strings.Split(accessToken, "Bearer ")
	accessToken = splitToken[1]

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	// ... error handling
	if err != nil {
		fmt.Println("Error validating jwt")
	}
	return claims
}
