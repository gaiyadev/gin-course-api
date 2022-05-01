package routes

import (
	"errors"
	"gin-course/custom"
	"gin-course/database"
	"gin-course/helpers"
	"gin-course/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
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

	err := database.DB.Where("name = ?", category.Name).First(&category).Error
	if err == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"status":     "Failed",
			"message":    "Already exist",
		})
		return
	}
	// Ok
	newCategory := models.Category{
		Name:   category.Name,
		UserID: userId,
	}

	ok := database.DB.Create(&newCategory).Error
	if ok != nil {
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
		"message":    "Post created successfully",
		"data": gin.H{
			"id":   newCategory.ID,
			"name": newCategory.Name,
		},
	})
	return
}

func FetchCategories(c *gin.Context) {
	var categories []models.Category

	err := database.DB.Order("id desc, name").Preload("User").Find(&categories).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"status":     "Failed",
			"message":    "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Fetched successfully",
		"statusCode": http.StatusOK,
		"status":     "Success",
		"data":       categories,
	})

}

func FetchUserCategories(c *gin.Context) {
	var categories []models.Category
	tokenString := c.GetHeader("Authorization")
	claims := helpers.ExtractClaims(tokenString)
	var userId = claims["id"].(float64)

	err := database.DB.Order("id desc, name").Where("user_id = ?", userId).Find(&categories).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"status":     "Failed",
			"message":    "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Fetched successfully",
		"statusCode": http.StatusOK,
		"status":     "Success",
		"data":       categories,
	})
	return
}

// FetchCategory by id
func FetchCategory(c *gin.Context) {
	id := c.Param("categoryId")
	var category []models.Category

	err := database.DB.Preload("User").First(&category, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"status":     "Failed",
			"message":    "Post not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Fetched successfully",
		"data":       category,
		"statusCode": http.StatusOK,
		"status":     "Success",
	})
	return
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("categoryId")

	tokenString := c.GetHeader("Authorization")
	claims := helpers.ExtractClaims(tokenString)
	var userId = claims["id"].(float64)

	var category []models.Category

	err := database.DB.Where("user_id = ?", userId).First(&category, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"status":     "Success",
			"message":    "Post not found",
		})
		return
	}

	deletePost := database.DB.Delete(&category, id).Error
	if deletePost != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"status":     "Failed",
			"message":    "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Deleted successfully",
		"data":       category,
		"statusCode": http.StatusOK,
		"status":     "Success",
	})
	return
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("categoryId")
	var category models.Category
	var updateCategory models.UpdateCategory

	//validate
	if err := c.ShouldBindJSON(&updateCategory); err != nil {
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

	err := database.DB.Where("id = ?", id).First(&category).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"status":     "Successful",
			"message":    "Post not found",
		})
		return
	}

	category.Name = updateCategory.Name
	result := database.DB.Save(&category).Error

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
		"message":    "Post updated successfully",
		"data": gin.H{
			"name": category.Name,
			"id":   category.ID,
		},
	})
	return
}
