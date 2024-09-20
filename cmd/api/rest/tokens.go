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

func (a *Application) creteAuthTokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ujson.Read(w, r, &input); err != nil {
		a.HandleBadRequest(w, r, "error while decoding input", err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePassword(v, input.Password)

	if !v.Valid() {
		a.HandleValidationErrors(w, r, v.Errors)
		return
	}

	user, err := a.Models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.HandleError(w, r, NewHttpErrorWithCause("invalid credentials", http.StatusUnauthorized, err))
		default:
			a.HandleError(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		a.HandleError(w, r, err)
		return
	}

	if !match {
		a.HandleError(w, r, NewHttpError("invalid credentials", http.StatusUnauthorized))
		return
	}

	token, err := a.Models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		a.HandleError(w, r, err)
		return
	}

	var output struct {
		AccessToken *data.Token `json:"access_token"`
	}

	output.AccessToken = token

	if err = ujson.Write(w, http.StatusCreated, output, nil); err != nil {
		a.HandleError(w, r, err)
	}

}
