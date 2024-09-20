package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/cmd/api/erro"
	"github.com/hvpaiva/greenlight/internal/data"
	"github.com/hvpaiva/greenlight/pkg/ujson"
	"github.com/hvpaiva/greenlight/pkg/validator"
)

func (h *Handler) CreteAuthTokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ujson.Read(w, r, &input); err != nil {
		return erro.BadRequest.WithMessage(err.Error())
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePassword(v, input.Password)

	if !v.Valid() {
		return erro.NewValidationErr("auth validation", v.Errors)
	}

	user, err := h.Models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return erro.Throw(erro.Unauthorized, erro.Cause("get user by email", err))
		default:
			return erro.ThrowInternalServer("get user by email", err)
		}
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		return erro.ThrowInternalServer("matching password", err)
	}

	if !match {
		return erro.Throw(erro.Unauthorized, erro.Cause("invalid credential", err))
	}

	token, err := h.Models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		return erro.ThrowInternalServer("create token", err)
	}

	var output struct {
		AccessToken *data.Token `json:"access_token"`
	}

	output.AccessToken = token

	if err = ujson.Write(w, http.StatusCreated, output, nil); err != nil {
		return erro.ThrowInternalServer("output response", err)
	}

	return nil
}
