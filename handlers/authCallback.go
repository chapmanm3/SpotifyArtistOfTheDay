package handlers

import (
	"SpotifyArtistofTheDay/main/database"
	"SpotifyArtistofTheDay/main/types"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
  "io"
	"strings"
)

func (h DBHandlerService) AuthCallback(c *gin.Context) {
	code := c.Query("code")

	clientID := os.Getenv("SAD_CLIENT_ID")
	clientSecret := os.Getenv("SAD_CLIENT_SECRET")
  serviceHostName := os.Getenv("AUTH_RETURN_URL")

  fmt.Printf("Service Hostname is %s \n", serviceHostName)

  if serviceHostName == "" {
    fmt.Println("No Hostname Header Found")
  }

	if code != "" {
		tokenObj, err := getAuthTokenResponse(code, clientID, clientSecret)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		userProfile, err := getSpotifyUserProfile(tokenObj.AccessToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't resolve Spotify Profile'")
		}
		database.SetUserInfo(h.DB, userProfile, tokenObj)

		c.SetCookie("auth_code", tokenObj.AccessToken, tokenObj.ExpiresIn, "/", c.Request.URL.Hostname(), false, false)
		c.Redirect(301, serviceHostName)
		//c.IndentedJSON(http.StatusOK, gin.H{"auth_code": tokenObj.AccessToken, "user Profile": *userProfile})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "No Code Given"})
	}
}

func getAuthTokenResponse(code string, clientID string, clientSecret string) (*types.AuthTokenResponse, error) {

  serviceHostName := os.Getenv("SERVICE_URL")

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
  form.Add("redirect_uri", fmt.Sprintf("%s/api/callback", serviceHostName))

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
			return nil, errors.New("Shit broke ouch")
		}
		var jsonBody types.AuthTokenResponse
		json.Unmarshal(body, &jsonBody)
		return &jsonBody, nil
	}
	return nil, err
}

func getSpotifyUserProfile(authToken string) (*types.UserProfileResponse, error) {
	userProfile := &types.UserProfileResponse{}

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
