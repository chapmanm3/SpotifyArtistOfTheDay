package handlers

import (
	"SpotifyArtistofTheDay/main/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h DBHandlerService) AddEmailToWaitlist(c *gin.Context) {

	var json struct {
		Email string `json:"email" binding:"required"`
	}

	if c.Bind(&json) == nil {
		database.AddToWaitlist(h.DB, json.Email)
		c.JSON(http.StatusOK, gin.H{"message": "email added successfully"})
	}
}
