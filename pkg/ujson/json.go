package ujson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"
)

func Write(w http.ResponseWriter, status int, data any, headers http.Header) error {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	res = append(res, '\n')

	maps.Insert(w.Header(), maps.All(headers))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(res); err != nil {
		return err
	}

	return nil
}

func Read(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		var (
			syntaxError           *json.SyntaxError
			unmarshalTypeError    *json.UnmarshalTypeError
			invalidUnmarshalError *json.InvalidUnmarshalError
			maxBytesError         *http.MaxBytesError
		)

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("request body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("request body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("request body contains incorrect JSON type (at position %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("request body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			return fmt.Errorf("request body contains unknown JSON field %s", strings.TrimPrefix(err.Error(), "json: unknown field "))
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("request body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must only contain a single JSON object")
	}

	return nil
}
