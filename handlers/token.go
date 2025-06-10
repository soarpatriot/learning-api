package handlers

import (
	"learning-api/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

var newDouyinClientFunc = helpers.NewDouyinClient

// TokenRequest represents the expected request body for /token
// Only a 'code' param is required

type TokenRequest struct {
	Code string `json:"code" binding:"required"`
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
	c.JSON(http.StatusOK, gin.H{"token": result})

}
