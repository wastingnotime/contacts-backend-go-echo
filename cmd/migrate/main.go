package main

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
)

type contact struct {
	ID          string
	FirstName   string
	LastName    string
	PhoneNumber string
}

func main() {
	godotenv.Load()
	dbLocation := os.Getenv("DB_LOCATION")

	db, err := gorm.Open("sqlite3", dbLocation)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&contact{})

	log.Println("migrations applied")
}
