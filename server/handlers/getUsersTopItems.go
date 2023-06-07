package handlers

import (
	"SpotifyArtistOfTheDay/types"
	"fmt"
	"net/http"
)

func GetUsersTopArtists (authToken string) {
  usersTopArtists := &types.UsersTopArtistsResponse{}

  fmt.Println("Getting Users Top Artists")

  client := &http.Client{}
  req, err := http.NewRequest(
    "GET",
    "https://api.spotify.com/v1/me/top/artists",
    nil,
  )
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", authToken))
}
