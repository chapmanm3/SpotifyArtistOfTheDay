package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

func openConnection() {
  
  connection, err := pgx.Connect(context.Background(), os.Getenv("SAD_DB_URL"))

  if err != nil {
    fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
  }
  defer connection.Close(context.Background())
}

func GetUserInfo() {

}
