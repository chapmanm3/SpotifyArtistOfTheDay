package handlers

import (
	"SpotifyArtistofTheDay/main/database"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

func (h DBHandlerService) GetUserInfo(c *gin.Context) {
	fmt.Println("Get User Info")
	//authCode := c.Request.Header.Get("auth_code")
	authCode, err := c.Request.Cookie("auth_code")

	if err != nil {
		fmt.Printf("No auth_code Cookie found")
		c.IndentedJSON(http.StatusForbidden, "Not Authorized")
	}

	user, err := database.GetUserInfo(h.DB, authCode.Value)

	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Fprintf(os.Stderr, "No User Found for that Auth ID")
		}
    fmt.Fprintf(os.Stderr, "GetUserInfo Query failed")
	}

	c.IndentedJSON(http.StatusOK, gin.H{"user_info": user})
}
