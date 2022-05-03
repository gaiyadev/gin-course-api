package routes

import (
	"errors"
	"gin-course/custom"
	"gin-course/database"
	"gin-course/helpers"
	"gin-course/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)


func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}

		limit, _ := strconv.Atoi(q.Get("limit"))
		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}



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
		Title:      post.Title,
		Body:       post.Body,
		UserID:     userId,
		CategoryID: post.CategoryID,
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
			"id":       newPost.ID,
			"title":    newPost.Title,
			"body":     newPost.Body,
			"category": newPost.CategoryID,
		},
	})
	return
}

func FetchPosts(c *gin.Context) {
	var posts []models.Post
	page := c.Query("page")
	limit := c.Query("limit")

	err := database.DB.Order("id desc, title").Preload("User").Preload("Category").Scopes(Paginate(c.Request)).Find(&posts).Error

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
		"meta": gin.H{
			"page":  page,
			"limit": limit,
		},
	})

}

func FetchUserPosts(c *gin.Context) {
	var posts []models.Post

	tokenString := c.GetHeader("Authorization")
	claims := helpers.ExtractClaims(tokenString)
	var userId = claims["id"].(float64)
	page := c.Query("page")
	limit := c.Query("limit")

	err := database.DB .Order("id desc, title").Where("user_id = ?", userId).Preload("Category").Scopes(Paginate(c.Request)).Find(&posts).Error

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
		"meta": gin.H{
			"page":  page,
			"limit": limit,
		},
	})
	return
}

// FetchPost by id
func FetchPost(c *gin.Context) {
	id := c.Param("postId")
	var post []models.Post

	err := database.DB.Preload("User").Preload("Category").Preload("Category").First(&post, id).Error

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

	err := database.DB.Where("user_id = ?", userId).First(&post, id).Error
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
