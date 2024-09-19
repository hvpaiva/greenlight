package rest

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/internal/data"
	"github.com/hvpaiva/greenlight/pkg/ujson"
	"github.com/hvpaiva/greenlight/pkg/validator"
)

func (a *Application) registerUserHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := ujson.Read(w, r, &input)
	if err != nil {
		a.HandleBadRequest(w, r, "error while decoding input", err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     data.Email(input.Email),
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		a.HandleError(w, r, err)
		return
	}

	v := validator.New()

	if user.Validate(v); !v.Valid() {
		a.HandleValidationErrors(w, r, v.Errors)
		return
	}

	err = a.Models.Users.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			a.HandleValidationErrors(w, r, v.Errors)
		default:
			a.HandleError(w, r, err)
		}
		return
	}

	if err = ujson.Write(w, http.StatusCreated, user, nil); err != nil {
		a.HandleError(w, r, err)
	}
}
