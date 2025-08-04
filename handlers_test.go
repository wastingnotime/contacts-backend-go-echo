package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var mockDB *gorm.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestCreateContact(t *testing.T) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/contacts", strings.NewReader(`{"firstName":"John","lastName":"Doe","phoneNumber":"1234567890"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handler{mockDB}

	if assert.NoError(t, h.CreateContact(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.NotEmpty(t, rec.Header().Get(echo.HeaderLocation))

		id := strings.Replace(rec.Header().Get(echo.HeaderLocation), "/contacts/", "", -1)

		var co contact
		mockDB.First(&co)
		assert.Equal(t, id, co.ID)
		assert.Equal(t, "John", co.FirstName)
		assert.Equal(t, "Doe", co.LastName)
		assert.Equal(t, "1234567890", co.PhoneNumber)

		mockDB.Delete(&co)
	}
}

func TestCreateContactValidation(t *testing.T) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/contacts", strings.NewReader(`{"lastName":"Doe","phoneNumber":"1234567890"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handler{mockDB}

	if assert.NoError(t, h.CreateContact(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateContact(t *testing.T) {
	sample := contact{ID: uuid.New().String(), FirstName: "John", LastName: "Doe", PhoneNumber: "1234567890"}
	mockDB.Create(sample)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.PUT, "/contacts/"+sample.ID, strings.NewReader(`{"firstName":"John1","lastName":"Doe1","phoneNumber":"12345678901"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handler{mockDB}

	if assert.NoError(t, h.UpdateContact(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)

		var co contact
		mockDB.Where(&contact{ID: sample.ID}).First(&co)
		assert.Equal(t, "John1", co.FirstName)
		assert.Equal(t, "Doe1", co.LastName)
		assert.Equal(t, "12345678901", co.PhoneNumber)

		mockDB.Delete(&co)
	}
}

func TestUpdateContactValidation(t *testing.T) {
	sample := contact{ID: uuid.New().String(), FirstName: "John", LastName: "Doe", PhoneNumber: "1234567890"}
	mockDB.Create(sample)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.PUT, "/contacts/"+sample.ID, strings.NewReader(`{"lastName":"Doe1","phoneNumber":"12345678901"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handler{mockDB}

	if assert.NoError(t, h.UpdateContact(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var co contact
		mockDB.Where(&contact{ID: sample.ID}).First(&co)
		assert.Equal(t, "John", co.FirstName)
		assert.Equal(t, "Doe", co.LastName)
		assert.Equal(t, "1234567890", co.PhoneNumber)

		mockDB.Delete(&co)
	}
}

func TestDeleteContact(t *testing.T) {
	sample := contact{ID: uuid.New().String(), FirstName: "John", LastName: "Doe", PhoneNumber: "1234567890"}
	mockDB.Create(sample)

	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/contacts/"+sample.ID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handler{mockDB}

	if assert.NoError(t, h.DeleteContact(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)

		var co contact
		result := mockDB.Where(&contact{ID: sample.ID}).First(&co)
		assert.Error(t, result.Error)
		assert.True(t, errors.Is(result.Error, gorm.ErrRecordNotFound))
	}
}

func TestGetContact(t *testing.T) {
	sample := contact{ID: uuid.New().String(), FirstName: "John", LastName: "Doe", PhoneNumber: "1234567890"}
	mockDB.Create(sample)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/contacts/"+sample.ID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handler{mockDB}

	if assert.NoError(t, h.GetContact(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		rsp := rec.Result()
		defer rsp.Body.Close()
		b, _ := io.ReadAll(rsp.Body)
		assert.Contains(t, string(b), sample.ID)
		assert.Contains(t, string(b), "John")
		assert.Contains(t, string(b), "Doe")
		assert.Contains(t, string(b), "1234567890")

		var co contact
		mockDB.Where(&contact{ID: sample.ID}).First(&co)
		mockDB.Delete(&co)
	}
}

func TestGetContacts(t *testing.T) {
	sample := contact{ID: uuid.New().String(), FirstName: "John", LastName: "Doe", PhoneNumber: "1234567890"}
	mockDB.Create(sample)
	sample1 := contact{ID: uuid.New().String(), FirstName: "John1", LastName: "Doe1", PhoneNumber: "12345678901"}
	mockDB.Create(sample1)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/contacts", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handler{mockDB}

	if assert.NoError(t, h.GetContacts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		rsp := rec.Result()
		defer rsp.Body.Close()
		b, _ := io.ReadAll(rsp.Body)
		assert.Contains(t, string(b), sample.ID)
		assert.Contains(t, string(b), "John")
		assert.Contains(t, string(b), "Doe")
		assert.Contains(t, string(b), "1234567890")
		assert.Contains(t, string(b), sample1.ID)
		assert.Contains(t, string(b), "John1")
		assert.Contains(t, string(b), "Doe1")
		assert.Contains(t, string(b), "12345678901")

		var co contact
		mockDB.Where(&contact{ID: sample.ID}).First(&co)
		mockDB.Delete(&co)
		mockDB.Where(&contact{ID: sample1.ID}).First(&co)
		mockDB.Delete(&co)
	}
}

func setup() {
	//database --------

	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(&contact{})

	mockDB = db
}

func teardown() {
	mockDB.Close()
}
