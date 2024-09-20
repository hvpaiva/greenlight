package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/cmd/api/app"
	"github.com/hvpaiva/greenlight/internal/data"
	"github.com/hvpaiva/greenlight/pkg/ujson"
	"github.com/hvpaiva/greenlight/pkg/validator"
)

func (h *Handler) CreteAuthTokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ujson.Read(w, r, &input); err != nil {
		h.App.HandleBadRequest(w, r, "error while decoding input", err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePassword(v, input.Password)

	if !v.Valid() {
		h.App.HandleValidationErrors(w, r, v.Errors)
		return
	}

	user, err := h.Models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.App.HandleError(w, r, app.NewErrorWithCause("invalid credentials", http.StatusUnauthorized, err))
		default:
			h.App.HandleError(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		h.App.HandleError(w, r, err)
		return
	}

	if !match {
		h.App.HandleError(w, r, app.NewError("invalid credentials", http.StatusUnauthorized))
		return
	}

	token, err := h.Models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		h.App.HandleError(w, r, err)
		return
	}

	var output struct {
		AccessToken *data.Token `json:"access_token"`
	}

	output.AccessToken = token

	if err = ujson.Write(w, http.StatusCreated, output, nil); err != nil {
		h.App.HandleError(w, r, err)
	}

}
