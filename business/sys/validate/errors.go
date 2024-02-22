package validate

import "errors"

var ErrInvalidID = errors.New("not a proper ID value")

type ErrorResponse struct {
	Error  string `json:"error"`
	Fields string `json:"fields,omitempty"`
}
