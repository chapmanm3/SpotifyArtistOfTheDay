package database

import (
	"context"
	"fmt"
	"os"

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

func GetUserInfo(userId string) UserInfo {
	dbpool := openConnection()
	defer dbpool.Close()

	var result UserInfo;
	err := dbpool.QueryRow(context.Background(), fmt.Sprintf("select * from userInfo where user_id = '%v'", userId)).Scan(&result.user_id, &result.email)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Query Row failed: %v \n", err)
	}

	fmt.Printf("Query Results: %v", result)

	return result
}

func SetAuthInfo(info types.AuthTokenResponse) {
	dbpool := openConnection()
	defer dbpool.Close()

	result, err := dbpool.Exec(
		context.Background(),
		"INSERT INTO authInfo (access_token, token_type, scope, expires_in, refresh_token) VALUES (?, ?, ?, ?, ?)",
		info.AccessToken,
		info.TokenType,
		info.Scope,
		info.ExpiresIn,
		info.RefreshToken,
	)

	if err != nil {
		fmt.Errorf("SetAuthInfo: %v", err)
	}
	fmt.Printf("SetAuthInfo Result: %v", result)
}
