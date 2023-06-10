package handlers

import (
	"SpotifyArtistOfTheDay/database"
	"SpotifyArtistOfTheDay/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GetUsersTopArtists(authToken string) {

	usersTopArtists, err := getUsersTopArtistsQuery(authToken, 0)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
  writeArtistsToDB(usersTopArtists.Items)

	totalResults := usersTopArtists.Total
	for currentOffset := 50; currentOffset <= totalResults; currentOffset += 50 {
    queryResults, err := getUsersTopArtistsQuery(authToken, currentOffset)
    if err != nil {
      fmt.Println(err)
    }
    writeArtistsToDB(queryResults.Items)
	}
}

func writeArtistsToDB(artists []types.ArtistObject) {
	for i := 0; i < len(artists); i++ {
		database.SetArtistInfo(&artists[i])
	}
}

func getUsersTopArtistsQuery(authToken string, offset int) (*types.UsersTopArtistsResponse, error) {
	usersTopArtists := &types.UsersTopArtistsResponse{}

	fmt.Println("Getting Users Top Artists")

	params := url.Values{}
	params.Add("time_range", "long_term")
	params.Add("limit", "50")
	params.Add("offset", fmt.Sprintf("%v", offset))

	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://api.spotify.com/v1/me/top/artists?"+params.Encode(),
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

	json.Unmarshal(body, usersTopArtists)

	fmt.Println(*usersTopArtists)

	return usersTopArtists, nil
}
