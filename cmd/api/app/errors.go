package app

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hvpaiva/greenlight/pkg/ujson"
)

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Implement error interface
func (e Error) Error() string {
	return e.Message
}

func NewError(message string, status int) Error {
	return Error{
		Message: message,
		Status:  status,
	}
}

func NewErrorWithCause(message string, status int, cause error) Error {
	return Error{
		Message: fmt.Sprintf("%s: %v", message, cause.Error()),
		Status:  status,
	}
}

func HandleError(w http.ResponseWriter, r *http.Request, a *Application, err error) {
	var httpErr Error
	if errors.As(err, &httpErr) {
		httpErr.Log(a, r)
		httpErr.Write(w)
	} else {
		e := NewError(err.Error(), http.StatusInternalServerError)
		e.Log(a, r)
		ErrInternalServerError.Write(w)
	}
}

func (a *Application) HandleError(w http.ResponseWriter, r *http.Request, err any) {
	if e, ok := err.(error); ok {
		HandleError(w, r, a, e)
		return
	}

	if s, ok := err.(string); ok {
		HandleError(w, r, a, errors.New(s))
		return
	}

	ErrInternalServerError.Log(a, r)
	Write(w, http.StatusInternalServerError, err)
}

func (a *Application) HandleErrors(w http.ResponseWriter, r *http.Request, message string, status int, errs map[string]string) {
	a.Log(r, message, slog.Any("errors", errs))
	Write(w, status, map[string]any{"message": message, "errors": errs})
}

func (e Error) Log(a *Application, r *http.Request) {
	a.Logger.Error(e.Message, slog.String("method", r.Method), slog.String("url", r.URL.String()))
}

func (a *Application) Log(r *http.Request, message string, args any) {
	a.Logger.Error(message, slog.String("method", r.Method), slog.String("url", r.URL.String()), args)
}

func (e Error) Write(w http.ResponseWriter) {
	if err := ujson.Write(w, e.Status, e, nil); err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
	}
}

func Write(w http.ResponseWriter, status int, data any) {
	err := ujson.Write(w, status, data, nil)
	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
	}
}

var (
	ErrInternalServerError = NewError("internal server error", http.StatusInternalServerError)
)

func NotFound(message string) Error {
	return NewError(message, http.StatusNotFound)
}

func (a *Application) NotFoundFunc(w http.ResponseWriter, r *http.Request) {
	err := NotFound(fmt.Sprintf("the requested resource could not be found: %s", r.URL.String()))
	a.HandleError(w, r, err)
}

func (a *Application) MethodNotAllowedFunc(w http.ResponseWriter, r *http.Request) {
	err := NewError(fmt.Sprintf(
		"the requested method %s is not allowed for the resource %s", r.Method, r.URL.String()),
		http.StatusMethodNotAllowed,
	)
	a.HandleError(w, r, err)
}

func (a *Application) HandleBadRequest(w http.ResponseWriter, r *http.Request, message string, err error) {
	a.HandleError(w, r, NewErrorWithCause(message, http.StatusBadRequest, err))
}

func (a *Application) HandleNotFound(w http.ResponseWriter, r *http.Request, message string, err error) {
	a.HandleError(w, r, NewErrorWithCause(message, http.StatusNotFound, err))
}

func (a *Application) HandleValidationErrors(w http.ResponseWriter, r *http.Request, err map[string]string) {
	a.HandleErrors(w, r, "error validating request data", http.StatusUnprocessableEntity, err)
}

func (a *Application) HandleConflict(w http.ResponseWriter, r *http.Request, message string, err error) {
	a.HandleError(w, r, NewErrorWithCause(message, http.StatusConflict, err))
}

func (a *Application) HandleUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	a.HandleError(w, r, NewError(message, http.StatusUnauthorized))
}

func (a *Application) HandleForbidden(w http.ResponseWriter, r *http.Request, message string) {
	a.HandleError(w, r, NewError(message, http.StatusForbidden))
}
