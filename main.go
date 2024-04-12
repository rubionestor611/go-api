package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type SpotifyCredentials struct {
	ClientID     string
	ClientSecret string
}

type SpotifyTopTracksResponse struct {
	Items []struct {
		Name    string `json:"name"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Album struct {
			Images []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
		} `json:"album"`
	} `json:"items"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type Track struct {
	Item struct {
		Name    string `json:"name"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
	} `json:"item"`
}

func main() {

	godotenv.Load(".env")

	credentials := SpotifyCredentials{
		ClientID:     os.Getenv("CLIENTID"),
		ClientSecret: os.Getenv("CLIENTSECRET"),
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://rubionestor611.github.io"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/nestor/spotify", func(c *gin.Context) {

		accessToken, err := getAccessToken(credentials)
		if err != nil {
			c.Status(500)
		}

		fmt.Println(accessToken)

		top, _ := getTopTracks(accessToken.AccessToken)
		current, _ := getNowListening(accessToken.AccessToken)

		c.JSON(http.StatusOK, gin.H{
			"topTracks": top.Items,
			"current":   current,
		})
	})

	router.Run(":8080")
}

func getAccessToken(credentials SpotifyCredentials) (*TokenResponse, error) {
	// Encode client ID and client secret in base64
	auth := base64.StdEncoding.EncodeToString([]byte(credentials.ClientID + ":" + credentials.ClientSecret))
	refresh := os.Getenv("REFRESH")
	// Prepare request body
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refresh)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResponse TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func getTopTracks(accessToken string) (*SpotifyTopTracksResponse, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/tracks", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var tracksResponse SpotifyTopTracksResponse
	err = json.NewDecoder(res.Body).Decode(&tracksResponse)
	if err != nil {
		return nil, err
	}

	return &tracksResponse, nil
}

func getNowListening(accessToken string) (*Track, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, nil
	}

	defer res.Body.Close()

	var nowPlaying Track
	err = json.NewDecoder(res.Body).Decode(&nowPlaying)
	if err != nil {
		return nil, err
	}

	fmt.Println("AAA", nowPlaying)

	return &nowPlaying, nil
}
