package middlewares_test

import (
	"learning-api/middlewares"
	"net/http"
	"net/http/httptest"
	"testing"

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
	token := models.NewToken()
	err := token.GenTokenWithDate()
	if err != nil {
		t.Fatal(err)
	}
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
	token := models.NewToken()

	accessTokenExpiresIn := -3600      // Set to expired
	refreshTokenExpiresIn := -31536000 // Set to expired
	token.SetAccessTokenAndRefreshToken(accessTokenExpiresIn, refreshTokenExpiresIn)
	return token, nil
}
