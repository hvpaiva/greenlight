package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/cmd/api/erro"
	"github.com/hvpaiva/greenlight/internal/data"
	"github.com/hvpaiva/greenlight/pkg/filters"
	"github.com/hvpaiva/greenlight/pkg/query"
	"github.com/hvpaiva/greenlight/pkg/ujson"
	"github.com/hvpaiva/greenlight/pkg/validator"
)

func (h *Handler) getMovieHandler(w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r)
	if err != nil {
		return erro.Throw(erro.BadRequest.WithMessage("invalid id"), erro.Cause("parsing id", err))
	}

	movie, err := h.Models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return erro.NotFound.WithMessage("the movie you are looking for does not exist")
		default:
			return erro.ThrowInternalServer("get movie", err)
		}
	}

	if err = ujson.Write(w, http.StatusOK, movie, nil); err != nil {
		return erro.ThrowInternalServer("output response", err)
	}

	return nil
}

func (h *Handler) createMovieHandler(w http.ResponseWriter, r *http.Request) error {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	if err := ujson.Read(w, r, &input); err != nil {
		return erro.BadRequest.WithMessage(err.Error())
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if movie.Validate(v); !v.Valid() {
		return erro.NewValidationErr("movie validation", v.Errors)
	}

	err := h.Models.Movies.Insert(movie)
	if err != nil {
		return erro.ThrowInternalServer("insert movie", err)
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	if err = ujson.Write(w, http.StatusCreated, movie, headers); err != nil {
		return erro.ThrowInternalServer("output response", err)
	}

	return nil
}

func (h *Handler) updateMovieHandler(w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r)
	if err != nil {
		return erro.Throw(erro.BadRequest.WithMessage("invalid id"), erro.Cause("parsing id", err))
	}

	movie, err := h.Models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return erro.NotFound.WithMessage("the movie you are looking for does not exist")
		default:
			return erro.ThrowInternalServer("get movie", err)
		}
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.Itoa(int(movie.Version)) != r.Header.Get("X-Expected-Version") {
			return erro.Conflict.WithMessage("the expected version does not match the current version of the movie")
		}
	}

	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	if err = ujson.Read(w, r, &input); err != nil {
		return erro.BadRequest.WithMessage(err.Error())
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	v := validator.New()

	if movie.Validate(v); !v.Valid() {
		return erro.NewValidationErr("movie validation", v.Errors)
	}

	if err = h.Models.Movies.Update(movie); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			return erro.Conflict.WithMessage("error while updating movie due to a conflict, please try again")
		default:
			return erro.ThrowInternalServer("update movie", err)
		}
	}

	if err = ujson.Write(w, http.StatusOK, movie, nil); err != nil {
		return erro.ThrowInternalServer("output response", err)
	}

	return nil
}

func (h *Handler) patchMovieHandler(w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r)
	if err != nil {
		return erro.Throw(erro.BadRequest.WithMessage("invalid id"), erro.Cause("parsing id", err))
	}

	movie, err := h.Models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return erro.NotFound.WithMessage("the movie you are looking for does not exist")
		default:
			return erro.ThrowInternalServer("get movie", err)
		}
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.Itoa(int(movie.Version)) != r.Header.Get("X-Expected-Version") {
			return erro.Conflict.WithMessage("the expected version does not match the current version of the movie")
		}
	}

	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	if err = ujson.Read(w, r, &input); err != nil {
		return erro.BadRequest.WithMessage(err.Error())
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}

	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	v := validator.New()

	if movie.Validate(v); !v.Valid() {
		return erro.NewValidationErr("movie validation", v.Errors)
	}

	if err = h.Models.Movies.Update(movie); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			return erro.Conflict.WithMessage("error while updating movie due to a conflict, please try again")
		default:
			return erro.ThrowInternalServer("update movie", err)
		}
	}

	if err = ujson.Write(w, http.StatusOK, movie, nil); err != nil {
		return erro.ThrowInternalServer("output response", err)
	}

	return nil
}

func (h *Handler) deleteMovieHandler(w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r)
	if err != nil {
		return erro.Throw(erro.BadRequest.WithMessage("invalid id"), erro.Cause("parsing id", err))
	}

	err = h.Models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return erro.NotFound.WithMessage("the movie you are looking for does not exist")
		default:
			return erro.ThrowInternalServer("get movie", err)
		}
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) showMoviesHandler(w http.ResponseWriter, r *http.Request) error {
	var input struct {
		Title  string
		Genres []string
		filters.Filter
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = query.ReadString(qs, "title", "")
	input.Genres = query.ReadCSV(qs, "genres", []string{})

	input.Page = query.ReadInt(qs, "page", 1, v)
	input.PageSize = query.ReadInt(qs, "page_size", 20, v)
	input.Sort = query.ReadString(qs, "sort", "id")
	input.SortSafeList = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if input.Filter.Validate(v); !v.Valid() {
		return erro.NewValidationErr("filter validation", v.Errors)
	}

	movies, metadata, err := h.Models.Movies.GetAll(input.Title, input.Genres, input.Filter)
	if err != nil {
		return erro.ThrowInternalServer("get all movies", err)
	}

	var output struct {
		Metadata filters.Metadata `json:"metadata"`
		Movies   []*data.Movie    `json:"movies"`
	}
	output.Movies = movies
	output.Metadata = metadata

	if err = ujson.Write(w, http.StatusOK, output, nil); err != nil {
		return erro.ThrowInternalServer("output response", err)
	}

	return nil
}

func parseId(r *http.Request) (int64, error) {
	param := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(param.ByName("id"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error while parsing id from params: %s", err.Error())
	}

	if id < 1 {
		return 0, fmt.Errorf("id provided is invalid, it should be greater than 0")
	}

	return id, nil
}
