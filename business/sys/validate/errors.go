package validate

import "errors"

// ErrInvalidID occurs when an ID is not in a valid form.
var ErrInvalidID = errors.New("not a proper ID value")

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Error  string `json:"error"`
	Fields string `json:"fields,omitempty"`
}
