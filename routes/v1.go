package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func V1(g *gin.RouterGroup) {
	g.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
}
