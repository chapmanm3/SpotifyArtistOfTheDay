package database

import (
	"SpotifyArtistOfTheDay/types"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openDB() gorm.DB {
  db, err := gorm.Open(postgres.Open(os.Getenv("SAD_DB_GORM_STRING")), &gorm.Config{})
  if err != nil {
    panic("Failed to Connect to DB")
  }
  return db
}

func seedDB() {
  db := openDB()
  db.AutoMigrate(&types.UserInfo{})

  
}


