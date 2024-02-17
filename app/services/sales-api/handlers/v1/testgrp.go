package testgrp

import (
	"context"
	"github.com/ktruedat/ultimateService/foundation/web"
	"go.uber.org/zap"
	"net/http"
)

type Handlers struct {
	Log *zap.SugaredLogger
}

func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	return web.Respond(ctx, w, status, http.StatusOK)
}
