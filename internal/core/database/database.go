package database

import (
	"fmt"
	 "food-delivery-workshop/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var (
	DB  *gorm.DB
	err error
)

func ConnectDB() {

	host := "localhost"
	port := "5432"
	user := "postgres"
	password := "1234"
	dbname := "food_delivery"
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to the database %v", err)
	}

	DB = DB.Debug()
	err = DB.AutoMigrate(&models.User{}, &models.Cart{}, &models.CartItem{}, &models.Promotion{}, &models.Product{})
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	log.Println("Database connection established successfully")

}
