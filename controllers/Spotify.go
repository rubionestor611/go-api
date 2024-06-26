package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

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

type SpotifyTrack struct {
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
}

type CurrentTrack struct {
	Item SpotifyTrack `json:"item"`
}

type SpotifyUser struct {
	Username  string `json:"display_name"`
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
	URI    string `json:"uri"`
	HREF   string `json:"href"`
	Images []struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"images"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

type RecentlyPlayed struct {
	Items []struct {
		Track struct {
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
		} `json:"track"`
	} `json:"items"`
}

func GetAccessToken(key string) (*TokenResponse, error) {
	godotenv.Load(".env")

	credentials := SpotifyCredentials{
		ClientID:     os.Getenv("CLIENTID"),
		ClientSecret: os.Getenv("CLIENTSECRET"),
	}

	// Encode client ID and client secret in base64
	auth := base64.StdEncoding.EncodeToString([]byte(credentials.ClientID + ":" + credentials.ClientSecret))
	refresh := os.Getenv(key)
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

func GetTopTracks() (*SpotifyTopTracksResponse, error) {
	token, err := GetAccessToken("REFRESH")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/tracks?time_range=short_term", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

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

func GetNowListening() (*CurrentTrack, error) {
	token, err := GetAccessToken("REFRESH")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, err
	}

	defer res.Body.Close()

	var nowPlaying CurrentTrack
	err = json.NewDecoder(res.Body).Decode(&nowPlaying)
	if err != nil {
		return nil, err
	}

	fmt.Println(nowPlaying)
	if len(nowPlaying.Item.Name) == 0 {
		return nil, fmt.Errorf("No now playing item able to be retrieved")
	}

	return &nowPlaying, nil
}

func GetSpotifyProfile() (*SpotifyUser, error) {
	token, err := GetAccessToken("REFRESH")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/users/rubiones2001", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var info SpotifyUser
	err = json.NewDecoder(res.Body).Decode(&info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func GetSpotifyLastPlayed() (*RecentlyPlayed, error) {
	token, err := GetAccessToken("REFRESH_RECENT")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/recently-played?limit=1", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	fmt.Println("Authorization", "Bearer "+token.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var recent RecentlyPlayed
	err = json.NewDecoder(res.Body).Decode(&recent)
	if err != nil {
		return nil, err
	}

	log.Println(res.Body)

	return &recent, nil
}
