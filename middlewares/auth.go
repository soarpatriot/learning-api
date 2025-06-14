package middlewares

import (
	"fmt"
	"learning-api/config"
	"learning-api/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/token" {
			c.Next()
			return
		}

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing", "code": "4970"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if !validateToken(tokenString, c) {
			c.Abort()
			return
		}

		user, _ := loadUserFromToken(tokenString)

		c.Set("currentUser", user)
		c.Next()
	}
}

func validateToken(tokenString string, c *gin.Context) bool {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.LoadConfig().ClientSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		if strings.Contains(err.Error(), "token is expired") {
			fmt.Println("err", err, "  error()", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "expired token", "code": "4980"})
			return false
		}
		fmt.Println("err", err, "  error()", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "code": "4981"})
		return false
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	// exp := int64(claims["exp"].(float64))
	// if exp < time.Now().Unix() {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Expired token"})
	// 	return false
	// }

	if !ok || claims["exp"] == nil || claims["iat"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims", "code": "4982"})
		return false
	}

	return true
}

func loadUserFromToken(tokenString string) (models.User, error) {
	var token models.Token
	if err := models.GetDB().Where("access_token = ?", tokenString).First(&token).Error; err != nil {
		return models.User{}, err
	}

	var user models.User
	if err := models.GetDB().Where("id = ?", token.UserID).First(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}
