package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/internal/data"
	"github.com/hvpaiva/greenlight/pkg/ujson"
	"github.com/hvpaiva/greenlight/pkg/validator"
)

func (h *Handler) registerUserHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := ujson.Read(w, r, &input)
	if err != nil {
		h.App.HandleBadRequest(w, r, "error while decoding input", err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     data.Email(input.Email),
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		h.App.HandleError(w, r, err)
		return
	}

	v := validator.New()

	if user.Validate(v); !v.Valid() {
		h.App.HandleValidationErrors(w, r, v.Errors)
		return
	}

	err = h.Models.Users.Insert(user)
	if err != nil {
		switch {
		// If we get h ErrDuplicateEmail error, use the v.AddError() method to manually
		// add h message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "h user with this email address already exists")
			h.App.HandleValidationErrors(w, r, v.Errors)
		default:
			h.App.HandleError(w, r, err)
		}
		return
	}

	if err = h.Models.Permission.AddForUser(user.ID, data.PermissionMovieRead); err != nil {
		h.App.HandleError(w, r, err)
		return
	}

	token, err := h.Models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		h.App.HandleError(w, r, err)
		return
	}

	var output struct {
		User  *data.User `json:"user"`
		Token string     `json:"token"`
	}
	output.User = user
	output.Token = token.Plaintext

	if err = ujson.Write(w, http.StatusCreated, output, nil); err != nil {
		h.App.HandleError(w, r, err)
	}
}

func (h *Handler) activateUserHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	if err := ujson.Read(w, r, &input); err != nil {
		h.App.HandleBadRequest(w, r, "error while decoding input", err)
		return
	}

	v := validator.New()

	if data.ValidateToken(v, input.TokenPlaintext); !v.Valid() {
		h.App.HandleValidationErrors(w, r, v.Errors)
		return
	}

	user, err := h.Models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired token")
			h.App.HandleValidationErrors(w, r, v.Errors)
		default:
			h.App.HandleError(w, r, err)
		}
		return
	}

	user.Activated = true

	err = h.Models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			h.App.HandleConflict(w, r, "unable to complete the request due to h conflict", err)
		default:
			h.App.HandleError(w, r, err)
		}
		return
	}

	err = h.Models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		h.App.HandleError(w, r, err)
		return
	}

	if err = ujson.Write(w, http.StatusOK, user, nil); err != nil {
		h.App.HandleError(w, r, err)
	}
}
