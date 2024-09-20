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

func (h *Handler) registerUserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := ujson.Read(w, r, &input)
	if err != nil {
		return erro.BadRequest.WithMessage(err.Error())
	}

	user := &data.User{
		Name:      input.Name,
		Email:     data.Email(input.Email),
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		return erro.ThrowInternalServer("setting password", err)
	}

	v := validator.New()

	if user.Validate(v); !v.Valid() {
		return erro.NewValidationErr("user validation", v.Errors)
	}

	err = h.Models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "user with this email address already exists")
			return erro.NewValidationErr("user insert", v.Errors)
		default:
			return erro.ThrowInternalServer("user insert", err)
		}
	}

	if err = h.Models.Permission.AddForUser(user.ID, data.PermissionMovieRead); err != nil {
		return erro.Throw(erro.InternalServer, erro.Cause("add permission", err))
	}

	token, err := h.Models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		return erro.ThrowInternalServer("create token", err)
	}

	var output struct {
		User  *data.User `json:"user"`
		Token string     `json:"token"`
	}
	output.User = user
	output.Token = token.Plaintext

	if err = ujson.Write(w, http.StatusCreated, output, nil); err != nil {
		return erro.Throw(erro.InternalServer, erro.Cause("output response", err))
	}

	return nil
}

func (h *Handler) activateUserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	if err := ujson.Read(w, r, &input); err != nil {
		return erro.BadRequest.WithMessage(err.Error())
	}

	v := validator.New()

	if data.ValidateToken(v, input.TokenPlaintext); !v.Valid() {
		return erro.NewValidationErr("token validation", v.Errors)
	}

	user, err := h.Models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired token")
			return erro.NewValidationErr("get user by token", v.Errors)
		default:
			return erro.ThrowInternalServer("get user by token", err)
		}
	}

	user.Activated = true

	err = h.Models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			return erro.Throw(erro.Conflict, erro.Cause("update user", err))
		default:
			return erro.ThrowInternalServer("update user", err)
		}
	}

	err = h.Models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		return erro.ThrowInternalServer("delete user tokens", err)
	}

	if err = ujson.Write(w, http.StatusOK, user, nil); err != nil {
		return erro.ThrowInternalServer("output response", err)
	}

	return nil
}
