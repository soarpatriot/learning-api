package middlewares_test

import (
	"learning-api/config"
	"learning-api/middlewares"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"learning-api/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	models.SetDB(models.InitTestDB()) // Initialize test database

	router := gin.Default()
	router.Use(middlewares.AuthMiddleware())
	router.GET("/topics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/topics", nil)
	user := models.User{OpenID: "test_openid", UnionID: "test_unionid", SessionKey: "test_sessionkey"}
	token, _ := models.GenerateToken()
	user.Tokens = []models.Token{*token}
	models.GetDB().Create(&user)
	models.GetDB().Save(&token)

	req.Header.Set("Authorization", "Bearer "+token.AccessToken) // Use generated valid token

	router.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "Access granted")
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	router := gin.Default()
	router.Use(middlewares.AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {

	models.SetDB(models.InitTestDB()) // Initialize test database
	router := gin.Default()
	router.Use(middlewares.AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	//set client secret environment variable

	user := models.User{OpenID: "test_openid", UnionID: "test_unionid", SessionKey: "test_sessionkey"}
	expiredToken, _ := generateTestToken()
	user.Tokens = []models.Token{*expiredToken}
	models.GetDB().Create(&user)
	models.GetDB().Save(&expiredToken)

	// Generate expired token
	req.Header.Set("Authorization", "Bearer "+expiredToken.AccessToken) // Use generated expired token

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "expired")
}

func generateTestToken() (*models.Token, error) {
	config := config.LoadConfig()
	const accessTokenExpiresIn = 3600      // 1 hour (seconds)
	const refreshTokenExpiresIn = 31536000 // 1 year (seconds)
	secret := config.ClientSecret
	token := &models.Token{
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		AccessTokenExpiresIn:  accessTokenExpiresIn,
		RefreshTokenExpiresIn: refreshTokenExpiresIn,
	}
	accessToken, err := models.GenJWTToken(secret, -time.Duration(accessTokenExpiresIn)*time.Second)
	if err != nil {
		return nil, err
	}
	token.AccessToken = accessToken
	refreshToken, err := models.GenJWTToken(secret, -time.Duration(refreshTokenExpiresIn)*time.Second)
	if err != nil {
		return nil, err
	}
	token.RefreshToken = refreshToken
	return token, nil
}
