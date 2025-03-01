package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type contact struct {
	ID          string `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
}

func main() {
	//environment --------
	godotenv.Load()
	environment := os.Getenv("ENVIRONMENT")
	dbLocation := os.Getenv("DB_LOCATION")

	//database --------
	db, err := gorm.Open("sqlite3", dbLocation)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&contact{})

	//api --------
	e := echo.New()

	if environment == "development" {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	e.HideBanner = true

	h := &handler{db}

	e.POST("/contacts", h.CreateContact)
	e.GET("/contacts", h.GetContacts)
	e.GET("/contacts/:id", h.GetContact)
	e.PUT("/contacts/:id", h.UpdateContact)
	e.DELETE("/contacts/:id", h.DeleteContact)

	e.Logger.Fatal(e.Start(":8010"))
}
