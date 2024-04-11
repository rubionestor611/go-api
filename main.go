package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type SpotifyCredentials struct {
	ClientID     string
	ClientSecret string
}

type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type SpotifyTopTracksResponse struct {
	Items []struct {
		Name string `json:"name"`
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

	// Obtain access token
	accessToken, err := getAccessToken(credentials)
	if err != nil {
		fmt.Println("Error obtaining access token:", err)
		return
	}

	// Retrieve top tracks
	topTracks, err := getTopTracks(accessToken)
	if err != nil {
		fmt.Println("Error retrieving top tracks:", err)
		return
	}

	// Print top tracks
	fmt.Println("Top Tracks:")
	for i, track := range topTracks {
		fmt.Printf("%d. %s\n", i+1, track.Name)
	}

	nowPlaying, err := getNowListening(accessToken)
	if err != nil {
		fmt.Println("Error retrieving current track:", err)
		return
	}

	fmt.Println(nowPlaying)

	router := gin.Default()

	router.GET("/nestor/spotify", func(c *gin.Context) {
		accessToken, err := getAccessToken(credentials)
		if err != nil {
			return
		}

		top, _ := getTopTracks(accessToken)
		current, _ := getNowListening(accessToken)

		fmt.Println(top, current)

		c.JSON(http.StatusOK, gin.H{
			"top": "hat",
		})
	})

	router.Run(":8080")
}

func getAccessToken(credentials SpotifyCredentials) (string, error) {
	// Encode client ID and client secret in base64
	auth := base64.StdEncoding.EncodeToString([]byte(credentials.ClientID + ":" + credentials.ClientSecret))
	refresh := os.Getenv("REFRESH")
	// Prepare request body
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refresh)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResponse TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func getTopTracks(accessToken string) ([]struct{ Name string }, error) {
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

	return []struct{ Name string }(tracksResponse.Items), nil
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
		return nil, fmt.Errorf("%d was returned", res.StatusCode)
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
