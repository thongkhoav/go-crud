package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thongkhoav/go-crud/initializers"
	"github.com/thongkhoav/go-crud/models"
)

func PostCreate(c *gin.Context) {
	var requestBody struct {
		Title string `json:"title" binding:"required"`
		Body  string `json:"body" binding:"required"`
	}
	errorBind := c.BindJSON(&requestBody)
	if errorBind != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorBind.Error()})
		return
	}

	// get user from context
	user, _ := c.Get("user")
	userData := user.(models.User)


	post := models.Post{Title: requestBody.Title, Body: requestBody.Body, UserID:  userData.ID}
	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	post.User = userData
	c.JSON(200, gin.H{
		"post": post,
	})
}

// Get all posts
func PostIndex(c *gin.Context) {
	var posts []models.Post
	initializers.DB.Preload("User").Find(&posts)

	c.JSON(200, gin.H{
		"posts": posts,
	})
}

// Get one post
func PostShow(c *gin.Context) {
	var post models.Post
	postId := c.Param("id")
	// initializers.DB.First(&post, postId)
	initializers.DB.Preload("User").First(&post, postId)

	if post.ID == 0 {
		c.JSON(404, gin.H{
			"error": "Post not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostUpdate(c *gin.Context) {
	var requestBody struct {
		Title string `json:"title" binding:"required"`
		Body  string `json:"body" binding:"required"`
	}
	errorBind := c.BindJSON(&requestBody)
	if errorBind != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorBind.Error()})
		return
	}

	var post models.Post
	postId := c.Param("id")
	initializers.DB.First(&post, postId)

	initializers.DB.Model(&post).Updates(models.Post{
		Title: requestBody.Title,
		Body:  requestBody.Body,
	})

	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostDelete(c *gin.Context) {
	postId := c.Param("id")

	initializers.DB.Delete(&models.Post{}, postId)

	c.JSON(200, gin.H{
		"message": "Post deleted",
	})
}
