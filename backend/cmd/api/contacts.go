package main

import (
	"errors"
	"fmt"
	"net/http"

	"backend/internal/models"
	"backend/internal/validator"
)

func (app *application) ListContactHandler(w http.ResponseWriter, r *http.Request) {
	contacts, err := app.models.Contacts.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, contacts, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createContactHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	contact := &models.Contact{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
	}

	v := validator.New()
	if ValidateContact(v, contact); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Contacts.Insert(contact)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/contacts/%d", contact.ID))

	err = app.writeJSON(w, http.StatusCreated, contact, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	app.writeWsContact("createContact", contact)
}

func (app *application) showContactHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	c, err := app.models.Contacts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, c, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateContactHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the movie ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the existing movie record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	c, err := app.models.Contacts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
	}

	// Decode the JSON as normal.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	changes := make(map[string]map[string]string)
	if input.FirstName != "" {
		changes["first_name"] = map[string]string{
			"from": c.FirstName,
			"to":   input.FirstName,
		}
		c.FirstName = input.FirstName
	}
	if input.LastName != "" {
		changes["last_name"] = map[string]string{
			"from": c.LastName,
			"to":   input.LastName,
		}
		c.LastName = input.LastName
	}
	if input.Email != "" {
		changes["email"] = map[string]string{
			"from": c.Email,
			"to":   input.Email,
		}
		c.Email = input.Email
	}

	if input.PhoneNumber != "" {
		changes["phone_number"] = map[string]string{
			"from": c.PhoneNumber,
			"to":   input.PhoneNumber,
		}
		c.PhoneNumber = input.PhoneNumber
	}

	v := validator.New()
	if ValidateContact(v, c); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	contact, err := app.models.Contacts.Update(c, changes)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, contact, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// websocket update
	app.writeWsContact("updateContact", contact)
}

func (app *application) deleteContactHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Contacts.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, map[string]string{"message": "The contact successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	contact := &models.Contact{ID: uint64(id)}
	app.writeWsContact("deleteContact", contact)
}

func (app *application) GetContactHistoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	history, err := app.models.Contacts.GetContactHistoryByID(uint64(id))
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, history, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func ValidateContact(v *validator.Validator, input *models.Contact) {
	v.Check(input.FirstName != "", "first_name", "must be provided")
	v.Check(len(input.FirstName) <= 500, "first_name", "must not be more than 500 bytes long")
	v.Check(input.LastName != "", "last_name", "must be provided")
	v.Check(len(input.LastName) <= 500, "last_name", "must not be more than 500 bytes long")
	v.Check(input.PhoneNumber != "", "phone_number", "must be provided")
	v.Check(input.Email != "", "email", "must be provided")
	v.Check(validator.Matches(input.Email, validator.EmailRX), "email", "must be a valid email address")
}
