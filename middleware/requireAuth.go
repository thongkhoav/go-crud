package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thongkhoav/go-crud/initializers"
	"github.com/thongkhoav/go-crud/models"
)

// Next or Abort the request if the user is not authenticated
func RequireAuth(c *gin.Context) {
	// get token from cookie
	tokenString, err := c.Cookie("Authorization")
	if tokenString == "" || err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	// parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	// validate token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check expiration
		if time.Now().Unix() > int64(claims["exp"].(float64)) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		
		var user models.User
		// get user from database
		initializers.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// set user in context
		c.Set("user", user)

		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
	}

}
