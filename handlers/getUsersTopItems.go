package handlers

import (
	"SpotifyArtistofTheDay/main/database"
	"SpotifyArtistofTheDay/main/types"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h DBHandlerService) GetUsersTopArtists(c *gin.Context) {
	authCode, err := c.Request.Cookie("auth_code")
	if err != nil {
		fmt.Printf("No auth_code Cookie found")
		c.IndentedJSON(http.StatusForbidden, "Not Authorized")
	}
	authToken := authCode.Value

  usersTopItems, err := getUsersTopItems(h.DB, authToken)

  if err != nil {
    c.IndentedJSON(http.StatusInternalServerError, gin.H{})
  }

	c.IndentedJSON(http.StatusOK, gin.H{"artist": *getRandomArtists(usersTopItems)})
}

func getUsersTopItems(db *gorm.DB, authToken string) ([]*types.ArtistInfo, error) {
	userId, err := database.GetUserID(db, authToken)
	if err != nil {
		fmt.Printf("User ID Not Found")
		return nil, err
	}

	items, err := database.GetUsersTopArtists(db, uint(*userId))
	if err != nil {
		fmt.Printf("User ID Not Found")
		return nil, err
	}
	if len(items) > 0 {
		return items, nil
	}

	itemsQuery, err := getUsersTopArtistsQuery(authToken, 0)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var itemsQueryTotal = itemsQuery.Items[:]
  fmt.Println("Writting Artists to DB")
	writeArtistsToDB(db, itemsQueryTotal)

	x := mapArtistResponseToArtistInfo(itemsQueryTotal)
	writeArtistsToUser(db, x, *userId)

	return x, nil
}

//Generics baby!
func getRandomArtists[K any] (artists []*K) *K {
	randIndex := rand.Intn(len(artists))
	fmt.Printf("Random Artist: %+v \n", artists[randIndex])
	return artists[randIndex]
}

func transformArtistObject(artist *types.ArtistObject) *types.ArtistInfo {
  var artistInfo types.ArtistInfo
	artistInfo = types.ArtistInfo{
		SpotifyUrl: artist.ExternalUrls.Spotify,
		SpotifyId:  artist.Id,
		Image:      artist.Images[0].Url,
		Name:       artist.Name,
		Uri:        artist.Uri,
	}
  return &artistInfo
}

func writeArtistsToDB(db *gorm.DB, artists []types.ArtistObject) {
	for i := 0; i < len(artists); i++ {
		database.SetArtistInfo(db, &artists[i])
	}
}

func writeArtistsToUser(db *gorm.DB, artists []*types.ArtistInfo, userId int) {
	for i := 0; i < len(artists); i++ {
		database.SetUsersTopArtists(db, userId, artists)
	}
}

func mapArtistResponseToArtistInfo(artists []types.ArtistObject) []*types.ArtistInfo {
	x := make([]*types.ArtistInfo, 0)
  if len(artists) == 0 {
    fmt.Println(fmt.Errorf("Empty Array Passed"))
  }
  for _, value := range artists {
    x = append(x, transformArtistObject(&value))
	}
	return x
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

	json.Unmarshal(body, usersTopArtists)

	return usersTopArtists, nil
}
