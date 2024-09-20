package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/internal/data"
	"github.com/hvpaiva/greenlight/pkg/filters"
	"github.com/hvpaiva/greenlight/pkg/query"
	"github.com/hvpaiva/greenlight/pkg/ujson"
	"github.com/hvpaiva/greenlight/pkg/validator"
)

func (h *Handler) getMovieHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := parseId(params)
	if err != nil {
		h.App.HandleNotFound(w, r, "movie not found", err)
		return
	}

	movie, err := h.Models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.App.HandleNotFound(w, r, "error while getting movie", err)
		default:
			h.App.HandleError(w, r, err)
		}
		return
	}

	if err = ujson.Write(w, http.StatusOK, movie, nil); err != nil {
		h.App.HandleError(w, r, err)
	}
}

func (h *Handler) createMovieHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	if err := ujson.Read(w, r, &input); err != nil {
		h.App.HandleBadRequest(w, r, "error while decoding input", err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if movie.Validate(v); !v.Valid() {
		h.App.HandleValidationErrors(w, r, v.Errors)
		return
	}

	err := h.Models.Movies.Insert(movie)
	if err != nil {
		h.App.HandleError(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	if err = ujson.Write(w, http.StatusCreated, movie, headers); err != nil {
		h.App.HandleError(w, r, err)
	}
}

func (h *Handler) updateMovieHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := parseId(params)
	if err != nil {
		h.App.HandleNotFound(w, r, "movie not found", err)
		return
	}

	movie, err := h.Models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.App.HandleNotFound(w, r, "error while getting movie", err)
		default:
			h.App.HandleError(w, r, err)
		}
		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.Itoa(int(movie.Version)) != r.Header.Get("X-Expected-Version") {
			h.App.HandleConflict(w, r,
				"the expected version does not match the current version of the movie",
				errors.New("version mismatch"),
			)
			return
		}
	}

	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	if err = ujson.Read(w, r, &input); err != nil {
		h.App.HandleBadRequest(w, r, "error while decoding input", err)
		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	v := validator.New()

	if movie.Validate(v); !v.Valid() {
		h.App.HandleValidationErrors(w, r, v.Errors)
		return
	}

	if err = h.Models.Movies.Update(movie); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			h.App.HandleConflict(w, r, "error while updating movie due to a conflict, please try again", err)
		default:
			h.App.HandleError(w, r, err)
		}
		return
	}

	if err = ujson.Write(w, http.StatusOK, movie, nil); err != nil {
		h.App.HandleError(w, r, err)
	}

}

func (h *Handler) patchMovieHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := parseId(params)
	if err != nil {
		h.App.HandleNotFound(w, r, "movie not found", err)
		return
	}

	movie, err := h.Models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.App.HandleNotFound(w, r, "error while getting movie", err)
		default:
			h.App.HandleError(w, r, err)
		}
		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.Itoa(int(movie.Version)) != r.Header.Get("X-Expected-Version") {
			h.App.HandleConflict(w, r,
				"the expected version does not match the current version of the movie",
				errors.New("version mismatch"),
			)
			return
		}
	}

	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	if err = ujson.Read(w, r, &input); err != nil {
		h.App.HandleBadRequest(w, r, "error while decoding input", err)
		return
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
		h.App.HandleValidationErrors(w, r, v.Errors)
		return
	}

	if err = h.Models.Movies.Update(movie); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			h.App.HandleConflict(w, r, "error while updating movie due to a conflict, please try again", err)
		default:
			h.App.HandleError(w, r, err)
		}
		return
	}

	if err = ujson.Write(w, http.StatusOK, movie, nil); err != nil {
		h.App.HandleError(w, r, err)
	}

}

func (h *Handler) deleteMovieHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := parseId(params)
	if err != nil {
		h.App.HandleNotFound(w, r, "movie not found", err)
		return
	}

	err = h.Models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.App.HandleNotFound(w, r, "error while deleting movie", err)
		default:
			h.App.HandleError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) showMoviesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
		h.App.HandleValidationErrors(w, r, v.Errors)
		return
	}

	movies, metadata, err := h.Models.Movies.GetAll(input.Title, input.Genres, input.Filter)
	if err != nil {
		h.App.HandleError(w, r, err)
		return
	}

	var output struct {
		Metadata filters.Metadata `json:"metadata"`
		Movies   []*data.Movie    `json:"movies"`
	}
	output.Movies = movies
	output.Metadata = metadata

	if err = ujson.Write(w, http.StatusOK, output, nil); err != nil {
		h.App.HandleError(w, r, err)
	}
}

func parseId(param httprouter.Params) (int64, error) {
	id, err := strconv.ParseInt(param.ByName("id"), 10, 64)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("error while parsing id from params: %s", err.Error()))
	}

	if id < 1 {
		return 0, errors.New(fmt.Sprintf("id provided is invalid, it should be greater than 0"))
	}

	return id, nil
}
