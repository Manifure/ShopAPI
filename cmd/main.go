package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Address struct {
	ID      int    `json:"id"`
	Country string `gorm:"column:country" json:"country"`
	City    string `gorm:"column:city" json:"city"`
	Street  string `gorm:"column:street" json:"street"`
}

type Client struct {
	ID               int       `json:"id"`
	Name             string    `gorm:"column:client_name" json:"client_name"`
	Surname          string    `gorm:"column:client_surname" json:"client_surname"`
	Birthday         string    `json:"birthday"`
	Gender           string    `json:"gender"`
	RegistrationDate time.Time `gorm:"column:registration_date" json:"registration_date"`
	AddressID        int       `gorm:"column:address_id" json:"address_id"`
}

var db *gorm.DB

func (Address) TableName() string {
	return "address"
}

func initDB() {
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

func addClientHandler(c *gin.Context) {
	var client Client

	if err := c.BindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error bind json": err.Error()})
		return
	}
	client.RegistrationDate = time.Now()

	if err := db.Create(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"client": client})
}

func deleteClientHandler(c *gin.Context) {
	clientID := c.Query("id")

	if err := db.Delete(&Client{}, clientID).Error; err != nil {
		return
	}
}

func getClientsByFullName(c *gin.Context) {
	clientName := c.Query("client_name")
	clientSurname := c.Query("client_surname")

	var client []Client

	if err := db.Where("client_name = ? AND client_surname = ?", clientName, clientSurname).Find(&client).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"client": client})

}

func addAddressHandler(c *gin.Context) {
	var address Address

	if err := c.BindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid input": err.Error()})
	}

	if err := db.Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"address": address})
}

func updateClientAddressHandler(c *gin.Context) {

	clientIdStr := c.Query("address_id")
	if clientIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"invalid input": clientIdStr})
		return
	}
	clientAddressID, err := strconv.Atoi(clientIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid input": clientIdStr})
	}

	var client Client

	if err := db.First(&client, clientAddressID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	var address Address

	if err := db.First(&address, "id = ?", client.AddressID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	if err := c.BindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid input": err.Error()})
		return
	}

	if err := db.Save(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error3": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"client": client})
}

func main() {
	initDB()

	router := gin.Default()

	router.POST("/client", addClientHandler)
	router.DELETE("/client", deleteClientHandler)
	router.GET("/client", getClientsByFullName)
	router.PUT("/client", updateClientAddressHandler)
	router.POST("/address", addAddressHandler)

	router.Run("localhost:8080")
}
