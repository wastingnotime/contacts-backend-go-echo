package main

import (
	"net/http"
	"os"

	"github.com/google/uuid"
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
	godotenv.Load()

	// ENVIRONMENT=development
	environment := os.Getenv("ENVIRONMENT")

	contacts := []contact{
		{
			ID:          uuid.New().String(),
			FirstName:   "Albert",
			LastName:    "Einstein",
			PhoneNumber: "1111-1111",
		}, {
			ID:          uuid.New().String(),
			FirstName:   "Mary",
			LastName:    "Curie",
			PhoneNumber: "2222-1111",
		},
	}

	e := echo.New()

	if environment == "development" {
		e.Use(middleware.Logger())
	}

	e.Use(middleware.Recover())
	e.HideBanner = true

	e.POST("/contacts", func(c echo.Context) error {
		co := new(contact)
		if err := c.Bind(co); err != nil {
			return err
		}

		co.ID = uuid.New().String()

		contacts = append(contacts, *co)

		return c.NoContent(http.StatusCreated)
	})

	e.GET("/contacts", func(c echo.Context) error {
		return c.JSON(http.StatusOK, contacts)
	})

	e.GET("/contacts/:id", func(c echo.Context) error {
		id := c.Param("id")
		var co contact
		for _, ct := range contacts {
			if ct.ID == id {
				co = ct
				break
			}
		}
		if co == (contact{}) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusOK, co)
	})

	e.PUT("/contacts/:id", func(c echo.Context) error {
		id := c.Param("id")

		cp := new(contact)
		if err := c.Bind(cp); err != nil {
			return err
		}

		var co contact
		for _, ct := range contacts {
			if ct.ID == id {
				co = ct
				break
			}
		}
		if co == (contact{}) {
			return c.NoContent(http.StatusNotFound)
		}

		co.FirstName = cp.FirstName
		co.LastName = cp.LastName
		co.PhoneNumber = cp.PhoneNumber

		return c.NoContent(http.StatusNoContent)
	})

	e.DELETE("/contacts/:id", func(c echo.Context) error {
		id := c.Param("id")

		var co contact
		var ix int = -1
		for i, ct := range contacts {
			if ct.ID == id {
				co = ct
				ix = i
				break
			}
		}
		if co == (contact{}) {
			return c.NoContent(http.StatusNotFound)
		}

		contacts = remove(contacts, ix)

		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":8010"))
}

func remove(s []contact, i int) []contact {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
