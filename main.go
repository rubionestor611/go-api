package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// SpotifyCredentials represents the client ID and client secret
type SpotifyCredentials struct {
	ClientID     string
	ClientSecret string
}

// SpotifyTokenResponse represents the response from Spotify when obtaining an access token
type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// SpotifyTopTracksResponse represents the response from Spotify when retrieving top tracks
type SpotifyTopTracksResponse struct {
	Items []struct {
		Name string `json:"name"`
	} `json:"items"`
}

type AlbumsRet struct {
	Res struct {
		Items []struct {
			Name string `json:"name"`
		} `json:"items"`
	} `json:"albums"`
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
}

func getAccessToken(credentials SpotifyCredentials) (string, error) {
	// Encode client ID and client secret in base64
	auth := base64.StdEncoding.EncodeToString([]byte(credentials.ClientID + ":" + credentials.ClientSecret))
	fmt.Println(auth)
	// Prepare request body
	payload := strings.NewReader("grant_type=client_credentials")
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Parse response
	var tokenResponse SpotifyTokenResponse
	err = json.NewDecoder(res.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func getTopTracks(accessToken string) ([]struct{ Name string }, error) {
	// req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/tracks", nil)
	// if err != nil {
	// 	return nil, err
	// }
	// req.Header.Add("Authorization", "Bearer "+accessToken)

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/search?q=ziggy+stardust&type=album", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var albumsRet AlbumsRet
	json.NewDecoder(res.Body).Decode(&albumsRet)
	fmt.Println(albumsRet)
	return nil, nil

	//var tracksResponse SpotifyTopTracksResponse
	//err = json.NewDecoder(res.Body).Decode(&tracksResponse)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return []struct{ Name string }(tracksResponse.Items), nil
}
