package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=77788999 dbname=postgres port=5432 sslmode=disable"
var err error
DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
	log.Fatal("Failed to connect to database:", err)
}
}