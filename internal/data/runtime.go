package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Runtime int

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

func (r *Runtime) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte(strconv.Quote("0 min")), nil
	}

	jsonValue := fmt.Sprintf("%d min", *r)

	return []byte(strconv.Quote(jsonValue)), nil
}

func (r *Runtime) UnmarshalJSON(bytes []byte) error {
	unquoted, err := strconv.Unquote(string(bytes))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquoted, " ")

	if len(parts) > 1 && parts[1] != "min" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(i)

	return nil
}

func (r *Runtime) String() string {
	if r == nil {
		return "0 min"
	}

	return fmt.Sprintf("%d min", *r)
}
