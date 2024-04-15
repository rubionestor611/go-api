package routes

import (
	"example/go-api/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Nestor(g *gin.RouterGroup) {
	g.GET("/spotify/profile", func(c *gin.Context) {
		user, err := controllers.GetSpotifyProfile()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to get Nestor's Spotify info"})
		}
		c.JSON(http.StatusOK, user)
	})
	g.GET("/spotify/currently-playing", func(c *gin.Context) {
		current, err := controllers.GetNowListening()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to get Nestor's currently playing Spotify Song"})
		}
		if current == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Nestor's currently not listening to Spotify music"})
		}
		c.JSON(http.StatusOK, current)
	})
	g.GET("/spotify/top-tracks", func(c *gin.Context) {
		top, err := controllers.GetTopTracks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to get Nestor's top Spotify tracks"})
		}
		c.JSON(http.StatusOK, top)
	})
}
