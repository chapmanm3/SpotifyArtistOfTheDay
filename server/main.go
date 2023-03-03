package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type ExplicitContentBody struct {
  FilterEnabled bool `json:"filter_enabled"`
  FilterLocked bool `json:"filter_locked"`
}

type ExternalUrlsBody struct {
  Spotify string `json:"spotify"`
}

type FollowersBody struct {
  Href string `json:"href"`
  Total int `json:"total"`
}

type ImagesBody struct {
  Url string `json:"url"`
  Height int `json:"height"`
  Width int `json:width`
}

type UserProfileResponse struct {
  Country string `json:"contry"`
  DisplayName string `json:"display_name"`
  Email string `json:"email"`
  ExplicitContent ExplicitContentBody `json:"explicit_content"`
  ExternalUrls ExternalUrlsBody `json:"external_urls"`
  Followers FollowersBody `json:"followers"`
  Href string `json:href`
  Images []ImagesBody `json:images`
  Product string `json:product`
  Type string `json:type`
  Uri string `json:uri`
}

func main() {
	router := gin.Default()

	router.GET("/login", authUser)
	router.GET("/callback", authCallback)

	router.Run()
}

func authUser(c *gin.Context) {

	fmt.Println("In Auth")

	clientID := os.Getenv("SAD_CLIENT_ID")
	fmt.Printf("clientID: %v", clientID)
	//clientSecret := os.Getenv("SAD_CLIENT_SECRET")
	redirectUri := "http://localhost:8080/callback"
	fmt.Printf("redirect URI: %v", redirectUri)
	scope := "user-read-private user-read-email user-top-read"

	c.Redirect(301, fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=%v&client_id=%v&scope=%v&redirect_uri=%v", "code", clientID, scope, redirectUri))

}

func authCallback(c *gin.Context) {
	code := c.Query("code")

	clientID := os.Getenv("SAD_CLIENT_ID")
	clientSecret := os.Getenv("SAD_CLIENT_SECRET")

	fmt.Printf("ClientID: %v", clientID)
	fmt.Println()

	fmt.Printf("ClientSecret: %v", clientSecret)
	fmt.Println()

	fmt.Printf("Code: %v", code)
	fmt.Println()

	if code != "" {
		token, err := getAuthToken(code, clientID, clientSecret)

		if token != "" {
      userProfile, err := getUserProfile(token)
      if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
      }
      c.IndentedJSON(http.StatusOK, gin.H{"auth_code": token, "user Profile": *userProfile})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		}
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "No Code Recieved"})
	}
}

func getAuthToken(code string, clientID string, clientSecret string) (string, error) {

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", "http://localhost:8080/callback")

	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		"https://accounts.spotify.com/api/token",
		strings.NewReader(form.Encode()),
	)
	authCode := "Basic " + b64.URLEncoding.EncodeToString([]byte(clientID+":"+clientSecret))
	req.Header.Add("Authorization", authCode)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)

	if res != nil {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", errors.New("Shit broke ouch")
		}
		var jsonBody AuthTokenResponse
		json.Unmarshal(body, &jsonBody)
		return string(jsonBody.AccessToken), nil
	}
	return "", err
}

func getUserProfile( authToken string ) ( *UserProfileResponse, error ) {
  userProfile := &UserProfileResponse{}

  fmt.Println("Getting user profile")
  
  client := &http.Client{}
  req, err := http.NewRequest(
    "GET",
    "https://api.spotify.com/v1/me",
    nil,
  )
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", authToken))

  if err != nil {
    fmt.Println(err)
    return nil, err
  }

  res, err := client.Do(req)

  if err != nil {
    fmt.Println(err)
    return nil, err
  }

  body, err := io.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  fmt.Println(body)

  json.Unmarshal(body, userProfile)

  fmt.Println(*userProfile)
  
  return userProfile, nil
}
