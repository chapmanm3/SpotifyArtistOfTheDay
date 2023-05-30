package database

import (
	"SpotifyArtistOfTheDay/types"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openDB() gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("SAD_DB_GORM_STRING")), &gorm.Config{})
	if err != nil {
		panic("Failed to Connect to DB")
	}
	return *db
}

func SeedDB() {
	db := openDB()

	db.AutoMigrate(&types.UserInfo{})
	db.AutoMigrate(&types.AuthInfo{})

	db.Create(&types.UserInfo{
		Country:         "US",
		DisplayName:     "Test",
		Email:           "testEmail",
		ExplicitContent: false,
		Followers:       1,
		ImageUrl:        "testURL",
		Uri:             "testUri",
		AuthInfo: types.AuthInfo{
			AccessToken:  "testToken",
			TokenType:    "testTokenType",
			Scope:        "testTokenScope",
			ExpiresIn:    15,
			RefreshToken: "testRefreshToken",
		},
	})
}

func SetUserInfo(userResponse *types.UserProfileResponse, authResponse *types.AuthTokenResponse) {
	db := openDB()

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

func GetUserInfo(authToken string) (*types.UserInfo, error) {
	db := openDB()
	var userInfo types.UserInfo
	fmt.Printf("authToken from GetUserInfo: %s", authToken)

	if authToken == "" {
		return nil, fmt.Errorf("No Auth Token Passed to GetUserInfo")
	}

	err := db.Where("AuthInfo.AccessToken = ?", authToken).First(&userInfo).Error

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &userInfo, nil
}
