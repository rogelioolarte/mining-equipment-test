package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofrs/uuid/v5"
)

func DecodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	// 1. Security limit: Maximum 1MB (adjust according to your needs)
	// This protects your memory from DoS attacks
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1 MB
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)

	// 2. Optional: Strict (Fails if fields not in the DTO are sent)
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		return err
	}

	// 3. Optional but recommended: Ensure there is no "garbage" after the JSON
	// Example: {"nametest"} {malicious_code}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must only contain a single JSON object")
	}

	return nil
}

// queryString returns the query param value or nil when absent/empty.
func queryString(r *http.Request, key string) *string {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	return &v
}

// queryUUID returns a *uuid.UUID for the given query param, or nil when absent.
// Writes a 400 and returns nil, false if the value is present but invalid.
func queryUUID(w http.ResponseWriter, r *http.Request, key string) (*uuid.UUID, bool) {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil, true
	}
	id, err := uuid.FromString(v)
	if err != nil {
		BadRequest(w, "invalid query param: "+key)
		return nil, false
	}
	return &id, true
}

// queryBool returns a *bool for the given query param, or nil when absent.
func queryBool(r *http.Request, key string) *bool {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return nil
	}
	return &b
}

// queryDate parses a date query param in "2006-01-02" format.
// Returns nil when absent.
func queryDate(r *http.Request, key string) *time.Time {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return nil
	}
	return &t
}

// queryPage extracts page/page_size query params with sensible defaults.
func queryPage(r *http.Request) (page int, pageSize int) {
	page = 1
	pageSize = 20
	if v := r.URL.Query().Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := r.URL.Query().Get("page_size"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			pageSize = n
		}
	}
	return
}
