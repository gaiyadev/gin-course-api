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

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
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

	// Ok
	newPost := models.Post{
		Title:  post.Title,
		Body:   post.Body,
		UserID: userId,
	}

	err := database.DB.Create(&newPost).Error
	if err != nil {
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
			"id":    newPost.ID,
			"title": newPost.Title,
			"body":  newPost.Body,
		},
	})
	return
}

func FetchPosts(c *gin.Context) {
	var posts []models.Post

	err := database.DB.Order("id desc, title").Find(&posts).Error

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
		"data":       posts,
	})

}

func FetchUserPosts(c *gin.Context) {
	var posts []models.Post
	tokenString := c.GetHeader("Authorization")
	claims := helpers.ExtractClaims(tokenString)
	var userId = claims["id"].(float64)

	err := database.DB.Order("id desc, title").Where("user_id = ?", userId).Find(&posts).Error

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
		"data":       posts,
	})
	return
}

// FetchPost by id
func FetchPost(c *gin.Context) {
	id := c.Param("postId")
	var post []models.Post

	err := database.DB.First(&post, id).Error

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
		"data":       post,
		"statusCode": http.StatusOK,
		"status":     "Success",
	})
	return
}

func DeletePost(c *gin.Context) {
	id := c.Param("postId")

	tokenString := c.GetHeader("Authorization")
	claims := helpers.ExtractClaims(tokenString)
	var userId = claims["id"].(float64)

	var post []models.Post

	err := database.DB.Where("user_id =", userId).First(&post, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"status":     "Success",
			"message":    "Post not found",
		})
		return
	}

	deletePost := database.DB.Delete(&post, id).Error
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
		"data":       post,
		"statusCode": http.StatusOK,
		"status":     "Success",
	})
	return
}

func UpdatePost(c *gin.Context) {
	id := c.Param("postId")
	var post models.Post
	var updatePost models.UpdatePost

	//validate
	if err := c.ShouldBindJSON(&updatePost); err != nil {
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

	err := database.DB.Where("id = ?", id).First(&post).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"status":     "Successful",
			"message":    "Post not found",
		})
		return
	}

	post.Title = updatePost.Title
	post.Body = updatePost.Body
	result := database.DB.Save(&post).Error

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
			"title": post.Title,
			"body":  post.Body,
			"id":    post.ID,
		},
	})
	return
}
