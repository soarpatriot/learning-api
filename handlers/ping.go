package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler handles GET /v1/ping
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
