package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Validator interface {
	Validate(ctx context.Context) map[string]string
}

func decodeAndValidate[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}
	if problems := v.Validate(r.Context()); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}
	return v, nil, nil
}
