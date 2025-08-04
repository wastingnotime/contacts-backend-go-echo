package main

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type contact struct {
	ID          string `json:"id"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
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
	e.Validator = &CustomValidator{validator: validator.New()}

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
