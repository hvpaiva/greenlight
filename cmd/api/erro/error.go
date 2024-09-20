package erro

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hvpaiva/greenlight/cmd/api/app"
	"github.com/hvpaiva/greenlight/pkg/ujson"
)

var (
	Forbidden           = New("request forbidden", http.StatusForbidden)
	Unauthorized        = New("request unauthorized", http.StatusUnauthorized)
	NotFound            = New("resource not found", http.StatusNotFound)
	BadRequest          = New("bad request", http.StatusBadRequest)
	Conflict            = New("conflict", http.StatusConflict)
	InternalServer      = New("internal server error", http.StatusInternalServerError)
	MethodNotAllowed    = New("method not allowed", http.StatusMethodNotAllowed)
	UnprocessableEntity = New("validation failed", http.StatusUnprocessableEntity)
	TooManyRequests     = New("too many request", http.StatusTooManyRequests)
)

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type InternalErr struct {
	Err   Error
	Cause string
}

type ValidationErr struct {
	Err  Error
	errs map[string]string
}

type output struct {
	Error Error `json:"error"`
}

func (e Error) Error() string {
	return e.Message
}

func New(message string, status int) Error {
	return Error{
		Message: message,
		Status:  status,
	}
}

func (e Error) WithMessage(message string) Error {
	e.Message = message
	return e
}

func Throw(err Error, cause string) InternalErr {
	return InternalErr{
		Err:   err,
		Cause: cause,
	}
}

func ThrowInternalServer(when string, err error) InternalErr {
	return Throw(InternalServer, Cause(when, err))
}

func NewValidationErr(when string, errs map[string]string) ValidationErr {
	return ValidationErr{
		Err:  UnprocessableEntity.WithMessage(fmt.Sprintf("%s error", when)),
		errs: errs,
	}
}

func (e InternalErr) Error() string {
	return e.Cause
}

func (e ValidationErr) Error() string {
	return e.Err.Message
}

func Cause(when string, err error) string {
	return fmt.Errorf("%s error: %w", when, err).Error()
}

func Handle(a *app.Application, w http.ResponseWriter, r *http.Request, err error) {
	var internal InternalErr
	if errors.As(err, &internal) {
		write(w, internal.Err.Status, output{Error: internal.Err})
		log(a, r, internal.Cause, slog.Any("error", internal.Err))
		return
	}

	var otherErr Error
	if errors.As(err, &otherErr) {
		write(w, otherErr.Status, output{Error: otherErr})
		log(a, r, otherErr.Message, slog.Any("error", otherErr))
		return
	}

	var validationErr ValidationErr
	if errors.As(err, &validationErr) {
		HandleMultiple(a, w, r, validationErr.Err, validationErr.errs)
	}

	write(w, http.StatusInternalServerError, output{Error: InternalServer})
	log(a, r, err.Error(), slog.Any("error", err))
}

func HandleMultiple(a *app.Application, w http.ResponseWriter, r *http.Request, err Error, errs map[string]string) {
	log(a, r, err.Message, slog.Any("errors", errs))
	write(w, err.Status, map[string]any{"message": err.Message, "errors": errs})
}

func write(w http.ResponseWriter, status int, data any) {
	err := ujson.Write(w, status, data, nil)
	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
	}
}

func log(a *app.Application, r *http.Request, message string, args any) {
	a.Logger.Error(message, slog.String("method", r.Method), slog.String("url", r.URL.String()), args)
}
