package routes

import (
	"errors"
	"gin-course/config"
	"gin-course/custom"
	"gin-course/database"
	"gin-course/helpers"
	"gin-course/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// SignUp user
func SignUp(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]custom.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = custom.ErrorMsg{Field: fe.Field(), Message: custom.GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors":     out,
				"statusCode": http.StatusBadRequest,
				//"e":          ve.Error(),
			})
		}
		return
	}

	//Hashing password
	hash, _ := HashPassword(user.Password)

	// Creating new user struct
	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hash,
	}

	// Checking if user already exist
	found := database.DB.Create(&newUser).Error
	if found != nil {
		c.JSON(http.StatusConflict, gin.H{
			"statusCode": http.StatusConflict,
			"status":     "Failed",
			"message":    "Account already exist",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"status":     "Success",
		"data": gin.H{
			"id":    newUser.ID,
			"email": newUser.Email,
			"name":  newUser.Name,
		},
		"message": "Account created successfully",
	})
	return
}

func SignIn(c *gin.Context) {
	var user models.User
	var login models.Login

	if err := c.ShouldBindJSON(&login); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]custom.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = custom.ErrorMsg{Field: fe.Field(), Message: custom.GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors":     out,
				"statusCode": http.StatusBadRequest,
				"status":     "Failed",
			})
		}
		return
	}

	user.Password = login.Password
	user.Email = login.Email

	err := database.DB.Where("email = ?", user.Email).Take(&user).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"status":     "Failed",
			"message":    "Email or Password is invalid",
		})
		return
	}

	// Verifying user password against the hash in the database
	match := VerifyPasswordHash(login.Password, user.Password)
	if !match {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"status":     "Fail",
			"message":    "Email or Password is invalid!",
		})
		return
	}

	// here, we have kept it as 60 minutes
	expirationTime := time.Now().Add(60 * time.Minute)
	var JwtSecret = []byte(config.Config("JWT_SECRET"))

	claims := &models.Claims{
		Name:  user.Name,
		Email: user.Email,
		ID:    user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, JwtErr := token.SignedString(JwtSecret)

	if JwtErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"status":     "failed",
		})
		return
	}

	//c.SetCookie("accessToken", accessToken, 60*60*24, "/", "http://localhost:8080", true, true)

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"status":     "Success",
		"data": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
		"accessToken": accessToken,
		"message":     "Sign in successfully",
	})
	return
}

func ChangePassword(c *gin.Context) {
	var user models.User
	var changePassword models.ChangePassword

	//validate
	if err := c.ShouldBindJSON(&changePassword); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]custom.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = custom.ErrorMsg{Field: fe.Field(), Message: custom.GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors":     out,
				"statusCode": http.StatusBadRequest,
			})
		}
		return
	}

	tokenString := c.GetHeader("Authorization")
	claims := helpers.ExtractClaims(tokenString)
	var userId = claims["id"].(float64)

	err := database.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"status":     "Successful",
			"message":    "Post not found",
		})
		return
	}

	match := VerifyPasswordHash(changePassword.CurrentPassword, user.Password)
	if !match {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"status":     "Failed",
			"message":    "Current password is not correct",
		})
		return
	}

	if changePassword.NewPassword != changePassword.ConfirmedPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"status":     "Failed",
			"message":    "Password mismatch",
		})
		return
	}

	hash, _ := HashPassword(changePassword.NewPassword)

	user.Password = hash
	result := database.DB.Save(&user).Error

	if result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"status":     "Failed",
			"message":    "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"status":     "Success",
		"message":    "Updated successfully",
		"data": gin.H{
			"name": user.Name,
			"id":   user.ID,
		},
	})
	return
}