package main

import (
	"ShopAPI/internal"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

//func (internal.Address) TableName() string {
//	return "address"
//}

func initDB() {
	var err error
	dsn := "host=localhost user=postgres password=1 dbname=postgres port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&internal.Client{}, &internal.Address{})
	if err != nil {
		log.Fatalf("Failed to migrate DB: %v", err)
	}

}

func main() {
	router := gin.Default()

	internal.InitDB()

	router.POST("/client", internal.AddClientHandler)
	router.DELETE("/client", internal.DeleteClientHandler)
	router.GET("/client", internal.GetClientsByFullName)
	router.PUT("/client", internal.UpdateClientAddressHandler)
	router.POST("/address", internal.AddAddressHandler)

	router.Run("localhost:1488")
}
