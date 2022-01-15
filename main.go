package main

import (
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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

	e.POST("/contacts", func(c echo.Context) error {
		payload := new(contact)
		if err := c.Bind(payload); err != nil {
			return err
		}

		payload.ID = uuid.New().String()

		db.Create(payload)

		return c.NoContent(http.StatusCreated)
	})

	e.GET("/contacts", func(c echo.Context) error {
		var contacts []contact
		db.Find(&contacts)
		return c.JSON(http.StatusOK, contacts)
	})

	e.GET("/contacts/:id", func(c echo.Context) error {
		id := c.Param("id")

		var co contact
		db.Where(&contact{ID: id}).First(&co)
		if co == (contact{}) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusOK, co)
	})

	e.PUT("/contacts/:id", func(c echo.Context) error {
		id := c.Param("id")

		payload := new(contact)
		if err := c.Bind(payload); err != nil {
			return err
		}

		var co contact
		db.Where(&contact{ID: id}).First(&co)
		if co == (contact{}) {
			return c.NoContent(http.StatusNotFound)
		}

		co.FirstName = payload.FirstName
		co.LastName = payload.LastName
		co.PhoneNumber = payload.PhoneNumber

		db.Save(&co)

		return c.NoContent(http.StatusNoContent)
	})

	e.DELETE("/contacts/:id", func(c echo.Context) error {
		id := c.Param("id")

		var co contact
		db.Where(&contact{ID: id}).First(&co)
		if co == (contact{}) {
			return c.NoContent(http.StatusNotFound)
		}

		db.Delete(&co)

		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":8010"))
}
