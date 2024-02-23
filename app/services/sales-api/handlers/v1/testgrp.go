package testgrp

import (
	"context"
	"errors"
	"github.com/ktruedat/ultimateService/foundation/web"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
)

type Handlers struct {
	Log *zap.SugaredLogger
}

func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return errors.New("untrusted error")
	}
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	return web.Respond(ctx, w, status, http.StatusOK)
}
