package database

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"SpotifyArtistOfTheDay/types"
)

type UserInfo struct {
	user_id int
	email   string
}

func openConnection() *pgxpool.Pool {

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("SAD_DB_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
	}
	//defer dbpool.Close()

	return dbpool
}

func GetUserInfo(userId string) (UserInfo, error) {
	dbpool := openConnection()
	defer dbpool.Close()

	var result UserInfo
	err := dbpool.QueryRow(context.Background(), fmt.Sprintf("select * from UserInfo where user_id = '%v'", userId)).Scan(&result.user_id, &result.email)

	fmt.Printf("Query Results: %v", result)

	return result, err
}

func SetUserInfo(user types.UserProfileResponse) {
	dbpool := openConnection()
	defer dbpool.Close()

	fmt.Printf("SetUserInfo User Value: %v \n", user)
	fmt.Printf("Country: %v \n", user.Country)
	fmt.Printf("DisplayName: %v \n", user.DisplayName)
	fmt.Printf("Email: %v \n", user.DisplayName)
	fmt.Printf("ExplicitContent Filter: %v \n", user.ExplicitContent.FilterEnabled)
	fmt.Printf("Followers: %v \n", user.Followers.Total)
	fmt.Printf("Image Url: %v \n", user.Images[0].Url)
	fmt.Printf("Image Url Length: %v \n", strings.Count(user.Images[0].Url, ""))
	fmt.Printf("Product: %v \n", user.Product)
	fmt.Printf("Type: %v \n", user.Type)
	fmt.Printf("Uri: %v \n", user.Uri)

	queryString := fmt.Sprintf(
		"INSERT INTO UserInfo (country, display_name, email, explicit_content, followers, image, product, type, uri) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.Country,
		user.DisplayName,
		user.Email,
		user.ExplicitContent.FilterEnabled,
		user.Followers.Total,
		user.Images[0].Url,
		user.Product,
		user.Type,
		user.Uri,
	)

	results, err := dbpool.Exec(context.Background(), queryString)

	if err != nil {
		fmt.Fprintf(os.Stderr, "SetUserInfo failure: %v", err)
	}

	fmt.Printf("SetUserInfo Results: %v", results)
}

func GetUserIdFromAuthToken(authToken string) string {
	dbpool := openConnection()
	defer dbpool.Close()

	var returnID string
	err := dbpool.QueryRow(
		context.Background(),
		fmt.Sprintf("select user_id from AuthInfo where access_token = '%v'", authToken),
	).Scan(&returnID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Idk man something fucked up: %v", err)
	}
	return returnID
}

func getUserIdFromEmail(email string) string {
	dbpool := openConnection()
	defer dbpool.Close()

	fmt.Printf("\nEmail: %v \n", email)

	var returnId string
	err := dbpool.QueryRow(
		context.Background(),
		fmt.Sprintf("select user_id from UserInfo where email = '%v'", email),
	).Scan(&returnId)

	if err != nil {
		fmt.Fprintf(os.Stderr, "getUserIdFromEmail failed: %v", err)
	}

	return returnId
}

func SetAuthInfo(info types.AuthTokenResponse, email string) {
	dbpool := openConnection()
	defer dbpool.Close()

	userId := getUserIdFromEmail(email)

	result, err := dbpool.Exec(
		context.Background(),
		"INSERT INTO AuthInfo (user_id, access_token, token_type, scope, expires_in, refresh_token) VALUES (?, ?, ?, ?, ?, ?)",
		userId,
		info.AccessToken,
		info.TokenType,
		info.Scope,
		info.ExpiresIn,
		info.RefreshToken,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "SetAuthInfo: %v", err)
	}
	fmt.Printf("SetAuthInfo Result: %v", result)
}
