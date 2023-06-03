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

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"

	"SpotifyArtistOfTheDay/database"
	"SpotifyArtistOfTheDay/types"
)

func main() {
	database.SeedDB()

  config := cors.DefaultConfig()
  config.AllowOrigins = []string{"http://localhost:5173"}
  config.AllowCredentials = true

	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./dist", true)))
	router.Use(cors.New(config))

	api := router.Group("/api")
	{
		api.GET("/login", authUser)
		api.GET("/callback", authCallback)
		api.GET("/userInfo", getUserInfo)
	}

	//database.GetUserInfo()

	router.Run()
}

func authUser(c *gin.Context) {

	fmt.Println("In Auth")

	clientID := os.Getenv("SAD_CLIENT_ID")
	fmt.Printf("clientID: %v", clientID)
	//clientSecret := os.Getenv("SAD_CLIENT_SECRET")
	redirectUri := "http://localhost:8080/api/callback"
	fmt.Printf("redirect URI: %v", redirectUri)
	scope := "user-read-private user-read-email user-top-read"

	c.Redirect(301, fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=%v&client_id=%v&scope=%v&redirect_uri=%v", "code", clientID, scope, redirectUri))
}

func authCallback(c *gin.Context) {
	code := c.Query("code")

	clientID := os.Getenv("SAD_CLIENT_ID")
	clientSecret := os.Getenv("SAD_CLIENT_SECRET")

	if code != "" {
		tokenObj, err := getAuthTokenResponse(code, clientID, clientSecret)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		userProfile, err := getSpotifyUserProfile(tokenObj.AccessToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't resolve Spotify Profile'")
		}
		database.SetUserInfo(userProfile, tokenObj)

		c.SetCookie("auth_code", tokenObj.AccessToken, 10, "/", c.Request.URL.Hostname(), false, true)
		c.Redirect(301, "http://localhost:5173/")
		//c.IndentedJSON(http.StatusOK, gin.H{"auth_code": token, "user Profile": *userProfile})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "No Code Given"})
	}
}

func getAuthTokenResponse(code string, clientID string, clientSecret string) (*types.AuthTokenResponse, error) {

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", "http://localhost:8080/api/callback")

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

func getUserInfo(c *gin.Context) {
	fmt.Println("Get User Info")
	//authCode := c.Request.Header.Get("auth_code")
	authCode, err := c.Request.Cookie("auth_code")

	if err != nil {
		fmt.Printf("No auth_code Cookie found")
		c.IndentedJSON(http.StatusForbidden, "Not Authorized")
	}

	user, err := database.GetUserInfo(authCode.Value)

	if err != nil {
		if err == pgx.ErrNoRows {
			spotifyUser, err := getSpotifyUserProfile(authCode.Value)
			if err != nil {
				fmt.Fprintf(os.Stderr, "getSpotifyUserProfile broke with err: %v", err)
			}
			c.IndentedJSON(http.StatusOK, gin.H{"spotify_user_info": spotifyUser})
		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{"user_info": user})
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
