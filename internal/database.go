package internal

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func InitDB() {
	var err error
	dsn := "host=localhost user=postgres password=1 dbname=postgres port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&Client{}, &Address{})
	if err != nil {
		log.Fatalf("Failed to migrate DB: %v", err)
	}

}
