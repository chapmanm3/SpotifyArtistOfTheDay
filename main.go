package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"

	"SpotifyArtistofTheDay/main/database"
	"SpotifyArtistofTheDay/main/handlers"
)

func main() {

	//Load Env File
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Unable to load .env file")
	}

	//Init newrelic
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("SpotifyArtistOfTheDay_BE"),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)

	DB := database.InitDB()
	h := handlers.New(DB)

	config := cors.DefaultConfig()
	config.AllowOrigins = strings.Split(os.Getenv("ALLOW_ORIGINS"), ",")
	//config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowCredentials = true

	router := gin.Default()

  router.Use(nrgin.Middleware(app))
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
