package main

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/labstack/echo/v4"
	"net/http"
)

type handler struct {
	db *gorm.DB
}

func (h *handler) CreateContact(c echo.Context) error {
	payload := new(contact)
	if err := c.Bind(payload); err != nil {
		return err
	}

	payload.ID = uuid.New().String()

	h.db.Create(payload)

	c.Response().Header().Set(echo.HeaderLocation, "/contacts/"+payload.ID)

	return c.NoContent(http.StatusCreated)
}

func (h *handler) GetContacts(c echo.Context) error {
	var contacts []contact
	h.db.Find(&contacts)
	return c.JSON(http.StatusOK, contacts)
}

func (h *handler) GetContact(c echo.Context) error {
	id := c.Param("id")

	var co contact
	h.db.Where(&contact{ID: id}).First(&co)
	if co == (contact{}) {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, co)
}

func (h *handler) UpdateContact(c echo.Context) error {
	id := c.Param("id")

	payload := new(contact)
	if err := c.Bind(payload); err != nil {
		return err
	}

	var co contact
	h.db.Where(&contact{ID: id}).First(&co)
	if co == (contact{}) {
		return c.NoContent(http.StatusNotFound)
	}

	co.FirstName = payload.FirstName
	co.LastName = payload.LastName
	co.PhoneNumber = payload.PhoneNumber

	h.db.Save(&co)

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) DeleteContact(c echo.Context) error {
	id := c.Param("id")

	var co contact
	h.db.Where(&contact{ID: id}).First(&co)
	if co == (contact{}) {
		return c.NoContent(http.StatusNotFound)
	}

	h.db.Delete(&co)

	return c.NoContent(http.StatusNoContent)
}
