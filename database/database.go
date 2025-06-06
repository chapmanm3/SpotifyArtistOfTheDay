package database

import (
	"SpotifyArtistofTheDay/main/types"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	db.AutoMigrate(&types.Waitlist{})

	return db
}

func SetUsersTopArtists(db *gorm.DB, userId int, artists []types.ArtistInfo) {
	user := &types.UserInfo{
		Model: gorm.Model{ID: uint(userId)},
	}
	//db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user)
	//db.Update("Artists", artists).Where(user)
	//db.Save(&user)
	db.Model(&user).Association("Artists").Append(artists)
}

func GetUsersTopArtists(db *gorm.DB, userId uint) ([]types.ArtistInfo, error) {
	user := types.UserInfo{Model: gorm.Model{ID: uint(userId)}}
	var artists []types.ArtistInfo
	err := db.Model(&user).Association("Artists").Find(&artists)
	if err != nil {
		return nil, err
	}
	return artists, nil
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

	db.Debug().Omit("Artists", "CurrentArtist").Create(&userInfoInsert)
	db.Save(&userInfoInsert)
}

func GetUserInfo(db *gorm.DB, authToken string) (*types.UserInfo, error) {
	var userInfo types.UserInfo

	if authToken == "" {
		return nil, fmt.Errorf("No Auth Token Passed to GetUserInfo")
	}

	authInfo, authErr := GetAuthInfo(db, authToken)

	if authErr != nil {
		fmt.Println(authErr)
		return nil, authErr
	}

	err := db.Joins("AuthInfo").Find(&userInfo, "access_token = ?", authInfo.AccessToken).Error

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &userInfo, nil
}

func GetAuthInfo(db *gorm.DB, authToken string) (*types.AuthInfo, error) {
	var authInfo types.AuthInfo

	err := db.Find(&authInfo, "access_token = ?", authToken).Error

	if err != nil {
		return nil, err
	}

	if checkTokenExpired(&authInfo) {
		return nil, fmt.Errorf("Auth Token Expired")
	}

	return &authInfo, nil
}

func checkTokenExpired(authToken *types.AuthInfo) bool {
	expireInDir, err := time.ParseDuration(fmt.Sprintf("%ds", authToken.ExpiresIn))

	if err != nil {
		fmt.Printf("Unable to format Expires In Duration")
		return true
	}

	tokenExpireTime := authToken.CreatedAt.Add(expireInDir)

	if tokenExpireTime.After(time.Now()) {
		return false
	}

	return true
}

func GetUserID(db *gorm.DB, authToken string) (*int, error) {
	var authRecord types.AuthInfo

	err := db.Where(types.AuthInfo{AccessToken: authToken}).First(&authRecord).Error

	if err != nil {
		fmt.Printf("Auth Token does not exist")
		return nil, err
	}

	return &authRecord.UserInfoID, nil
}

func SetArtistInfo(db *gorm.DB, artist *types.ArtistObject) {
	artistInfoInsert := &types.ArtistInfo{
		SpotifyUrl: artist.ExternalUrls.Spotify,
		SpotifyId:  artist.Id,
		Image:      artist.Images[0].Url,
		Name:       artist.Name,
		Uri:        artist.Uri,
	}
	db.Save(&artistInfoInsert)
}

func SetUsersCurrentArtist(db *gorm.DB, userId uint, artist *types.ArtistInfo) {
	user := &types.UserInfo{
		Model:           gorm.Model{ID: uint(userId)},
		CurrentArtistID: artist.ID,
	}
	db.Model(&user).Update("CurrentArtistID", artist.ID)
}

func GetUsersCurrentArtist(db *gorm.DB, userId uint) (*types.ArtistInfo, error) {
	user := types.UserInfo{
		Model: gorm.Model{ID: userId},
	}
	db.Find(&user)

	artist := types.ArtistInfo{
		Model: gorm.Model{ID: user.CurrentArtistID},
	}
	db.Find(&artist)

	return &artist, nil
}

func AddToWaitlist(db *gorm.DB, email string) {
	entry := &types.Waitlist{
		Email: email,
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&entry)
}
