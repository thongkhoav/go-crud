package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thongkhoav/go-crud/initializers"
	"github.com/thongkhoav/go-crud/models"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func Signup(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,gt=5,lt=20"`
	}
	errorBind := c.BindJSON(&requestBody)
	if errorBind != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorBind.Error()})
		return
	}

	// check if user already exists
	userFound := models.User{}
	userExist := initializers.DB.Where(&models.User{Email: requestBody.Email}).First(&userFound)
	if userExist.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	// hash password
	hash, errHash := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if errHash != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while hashing password"})
		return
	}

	user := models.User{Email: requestBody.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully!",
	})
}

func Login(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	errorBind := c.BindJSON(&requestBody)
	if errorBind != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorBind.Error()})
		return
	}

	user := models.User{}
	initializers.DB.Where(&models.User{Email: requestBody.Email}).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	errCompareHash := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if errCompareHash != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while generating token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600 * 30, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})
}

func TestAuthenication(c *gin.Context) {
	user, _ := c.Get("user")

	// return user without password
	type UserResponse struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
	}
	userResponse := UserResponse{ID: user.(models.User).ID, Email: user.(models.User).Email}

	c.JSON(http.StatusOK, gin.H{
		"user": userResponse,
	})
}
