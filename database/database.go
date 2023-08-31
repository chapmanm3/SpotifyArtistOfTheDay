package database

import (
	"SpotifyArtistofTheDay/main/types"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("SAD_DB_GORM_STRING")), &gorm.Config{})
	if err != nil {
		panic("Failed to Connect to DB")
	}
	return db
}

func InitDB() *gorm.DB {
	db := openDB()

	db.AutoMigrate(&types.UserInfo{})
	db.AutoMigrate(&types.AuthInfo{})
	db.AutoMigrate(&types.ArtistInfo{})

  return db
}

func SetUserInfo(db *gorm.DB, userResponse *types.UserProfileResponse, authResponse *types.AuthTokenResponse) {
	userInfoInsert := types.UserInfo{
		Country:         userResponse.Country,
		DisplayName:     userResponse.DisplayName,
		Email:           userResponse.Email,
		ExplicitContent: userResponse.ExplicitContent.FilterEnabled,
		Followers:       userResponse.Followers.Total,
		ImageUrl:        userResponse.Images[0].Url,
		Uri:             userResponse.Uri,
		AuthInfo: types.AuthInfo{
			AccessToken:  authResponse.AccessToken,
			TokenType:    authResponse.TokenType,
			Scope:        authResponse.Scope,
			ExpiresIn:    authResponse.ExpiresIn,
			RefreshToken: authResponse.RefreshToken,
		},
	}

	db.Create(&userInfoInsert)
	db.Save(&userInfoInsert)
}

func  GetUserInfo(db *gorm.DB, authToken string) (*types.UserInfo, error) {
	var userInfo types.UserInfo
	fmt.Printf("authToken from GetUserInfo: %s", authToken)

	if authToken == "" {
		return nil, fmt.Errorf("No Auth Token Passed to GetUserInfo")
	}

	err := db.Joins("AuthInfo").Find(&userInfo, "access_token = ?", authToken).Error

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &userInfo, nil
}

func  SetArtistInfo(db *gorm.DB, artist *types.ArtistObject) {
	artistInfoInsert := types.ArtistInfo{
		SpotifyUrl: artist.ExternalUrls.Spotify,
		SpotifyId:  artist.Id,
		Image:      artist.Images[0].Url,
		Name:       artist.Name,
		Uri:        artist.Uri,
	}

	db.Create(&artistInfoInsert)
	db.Save(&artistInfoInsert)
}
