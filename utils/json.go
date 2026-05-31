package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

// Write JSON into http.ResponseWriter.
func WriteJSON(w http.ResponseWriter, httpStatus int, data any) {
	const op = "pkg.utils.json.WriteJSON"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	err := json.NewEncoder(w).Encode(data)
	if errors.Is(err, &json.UnsupportedTypeError{}) {
		msg := fmt.Sprintf("%s:%s", op, "error happened while encoding json data")
		http.Error(w, msg, http.StatusInternalServerError)
	}
}

type message struct {
	status int
	text   string
}

// Write 'http.Status' and 'msg' into http.ResponseWriter in JSON format.
func WriteJSONError(w http.ResponseWriter, httpStatus int, errMsg string) {
	msg := message{httpStatus, errMsg}
	WriteJSON(w, httpStatus, msg)
}

// Write 'http.WriteJSONErrorBadRequest' into http.ResponseWriter in JSON format.
func WriteJSONErrorBadRequest(w http.ResponseWriter) {
	msg := message{http.StatusBadRequest, http.StatusText(http.StatusBadRequest)}
	WriteJSON(w, http.StatusBadRequest, msg)
}

// Write 'http.StatusInternalServerError' into http.ResponseWriter in JSON format.
func WriteJSONErrorInternalServer(w http.ResponseWriter) {
	msg := message{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)}
	WriteJSON(w, http.StatusInternalServerError, msg)
}

// Read JSON from 'http.Request' into 'target'.
func ReadJSON(r *http.Request, target any) error {
	const op = "pkg.utils.json.ReadJSON"
	slog.Debug("trying to read json", "op", op)
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			return errors.New("content-type should be 'application/json'")
		}
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&target)
	if err != nil {
		var (
			syntaxError        *json.SyntaxError
			unmarshalTypeError *json.UnmarshalTypeError
		)
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("request body contains badly-formed JSON (at position %d)",
				syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf(
				"request body contains an invalid value for the %q field (at position %d)",
				unmarshalTypeError.Field, unmarshalTypeError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("request body contains badly-formed JSON")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("request body contains unknown field %s", fieldName)
		case errors.Is(err, io.EOF):
			return fmt.Errorf("request body must not be empty")
		default:
			return err
		}
	}
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("request body must only contain a single JSON object")
	}
	return nil
}
