package web

import (
	"context"
	"net/http"
)

func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {
	return nil
}
