package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

type handler struct {
	db *gorm.DB
}

func (h *handler) CreateContact(c echo.Context) error {
	payload := new(contact)
	if err := c.Bind(payload); err != nil {
		return err
	}

	if err := c.Validate(payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	payload.ID = uuid.New().String()
	result := h.db.Create(payload)
	if result.Error != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Response().Header().Set(echo.HeaderLocation, "/contacts/"+payload.ID)

	return c.NoContent(http.StatusCreated)
}

func (h *handler) GetContacts(c echo.Context) error {
	var contacts []contact
	result := h.db.Find(&contacts)
	if result.Error != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, contacts)
}

func (h *handler) GetContact(c echo.Context) error {
	id := c.Param("id")

	var co contact
	result := h.db.Where(&contact{ID: id}).First(&co)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.NoContent(http.StatusNotFound)
	}
	if result.Error != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, co)
}

func (h *handler) UpdateContact(c echo.Context) error {
	id := c.Param("id")

	payload := new(contact)
	if err := c.Bind(payload); err != nil {
		return err
	}

	if err := c.Validate(payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var co contact
	result := h.db.Where(&contact{ID: id}).First(&co)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.NoContent(http.StatusNotFound)
	}
	if result.Error != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	co.FirstName = payload.FirstName
	co.LastName = payload.LastName
	co.PhoneNumber = payload.PhoneNumber

	result = h.db.Save(&co)
	if result.Error != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) DeleteContact(c echo.Context) error {
	id := c.Param("id")

	var co contact
	result := h.db.Where(&contact{ID: id}).First(&co)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.NoContent(http.StatusNotFound)
	}
	if result.Error != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	result = h.db.Delete(&co)
	if result.Error != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
