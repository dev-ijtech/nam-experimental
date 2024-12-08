package namhttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dev-ijtech/nam-experimental"
)

type Config struct {
	Addr string
	Port int
}

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func decodeValid[T nam.Validator](r *http.Request) (T, nam.ProblemSet, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nam.ProblemSet{}, fmt.Errorf("decode json: %w", err)
	}

	if problems := v.Valid(); len(problems.Set) == 1 {
		return v, problems, fmt.Errorf("invalid %T: %d problem", v, len(problems.Set))
	} else if len(problems.Set) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems.Set))
	}

	return v, nam.ProblemSet{}, nil
}
