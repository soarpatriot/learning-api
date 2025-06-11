package handlers

import (
	"learning-api/helpers"
	"learning-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var newDouyinClientFunc = helpers.NewDouyinClient

// TokenRequest represents the expected request body for /token
// Only a 'code' param is required

type TokenRequest struct {
	Code string `json:"code" binding:"required"`
}

// RefreshTokenRequest represents the expected request body for /refresh-token
type RefreshTokenRequest struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// PostToken handles POST /token
func PostToken(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code param is required"})
		return
	}

	client := newDouyinClientFunc()
	result, err := client.Jscode2session(req.Code, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Here you would typically return the token to the client
	c.JSON(http.StatusOK, result)

}

// PostRefreshToken handles POST /refresh-token
func PostRefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access_token and refresh_token params are required"})
		return
	}

	// Skip auth middleware validation
	// Check if tokens exist in the token table
	token, err := models.FindToken(req.AccessToken, req.RefreshToken)
	if err != nil || token == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid tokens"})
		return
	}

	// Check if refresh token is expired
	expirationTime := token.CreatedAt.Add(time.Duration(token.RefreshTokenExpiresIn) * time.Second)
	if expirationTime.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	// Generate new token
	newToken, err := token.RefreshToNewToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		return
	}

	// Return new token
	c.JSON(http.StatusOK, newToken)
}
