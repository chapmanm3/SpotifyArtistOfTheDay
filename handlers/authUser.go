package handlers

import (
  "github.com/gin-gonic/gin"
  "fmt"
  "os"
)

func (h DBHandlerService) AuthUser(c *gin.Context) {

	fmt.Println("In Auth")

	clientID := os.Getenv("SAD_CLIENT_ID")
	fmt.Printf("clientID: %v", clientID)
	//clientSecret := os.Getenv("SAD_CLIENT_SECRET")
	redirectUri := "http://localhost:8080/api/callback"
	fmt.Printf("redirect URI: %v", redirectUri)
	scope := "user-read-private user-read-email user-top-read"

	c.Redirect(301, fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=%v&client_id=%v&scope=%v&redirect_uri=%v", "code", clientID, scope, redirectUri))
}
