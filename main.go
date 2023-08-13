package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"SpotifyArtistofTheDay/main/database"
  "SpotifyArtistofTheDay/main/handlers"
)

func main() {
  DB := database.InitDB()
  h := handlers.New(DB)

  config := cors.DefaultConfig()
  config.AllowOrigins = []string{"http://localhost:5173"}
  config.AllowCredentials = true

	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./dist", true)))
	router.Use(cors.New(config))

	api := router.Group("/api")
	{
    api.GET("/healthCheck", h.GetHealthCheck)
		api.GET("/login", h.AuthUser)
		api.GET("/callback", h.AuthCallback)
		api.GET("/userInfo", h.GetUserInfo)
    api.GET("/artist/today", h.GetUsersTopArtists)
	}

	router.Run()
}



