package rest

import (
	"errors"
	"net/http"
	"time"

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

	if err = a.Models.Permission.AddForUser(user.ID, data.PermissionMovieRead); err != nil {
		a.HandleError(w, r, err)
		return
	}

	token, err := a.Models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		a.HandleError(w, r, err)
		return
	}

	var output struct {
		User  *data.User `json:"user"`
		Token string     `json:"token"`
	}
	output.User = user
	output.Token = token.Plaintext

	if err = ujson.Write(w, http.StatusCreated, output, nil); err != nil {
		a.HandleError(w, r, err)
	}
}

func (a *Application) activateUserHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	if err := ujson.Read(w, r, &input); err != nil {
		a.HandleBadRequest(w, r, "error while decoding input", err)
		return
	}

	v := validator.New()

	if data.ValidateToken(v, input.TokenPlaintext); !v.Valid() {
		a.HandleValidationErrors(w, r, v.Errors)
		return
	}

	user, err := a.Models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired token")
			a.HandleValidationErrors(w, r, v.Errors)
		default:
			a.HandleError(w, r, err)
		}
		return
	}

	user.Activated = true

	err = a.Models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			a.HandleConflict(w, r, "unable to complete the request due to a conflict", err)
		default:
			a.HandleError(w, r, err)
		}
		return
	}

	err = a.Models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		a.HandleError(w, r, err)
		return
	}

	if err = ujson.Write(w, http.StatusOK, user, nil); err != nil {
		a.HandleError(w, r, err)
	}
}
