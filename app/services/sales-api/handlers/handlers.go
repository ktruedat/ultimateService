// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"expvar"
	testgrp "github.com/ktruedat/ultimateService/app/services/sales-api/handlers/v1"
	"github.com/ktruedat/ultimateService/business/web/mid"
	"github.com/ktruedat/ultimateService/foundation/web"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
	"os"
)

// DebugStandardLibraryMux registers all the debug routes from the standard library
// into a new mux bypassing in the use of DefaultServerMux. Using the
// DefaultServerMux would be a security risk since a dependency could inject a
// handler into our service without us knowing it
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library endpoints
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

func APIMux(cfg APIMuxConfig) *web.App {
	// Construct the web.App which holds all routes.
	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
	)

	// Load the routes for the different versions of the API.
	v1(app, cfg)

	return app
}

func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"
	tgh := testgrp.Handlers{
		Log: cfg.Log,
	}
	app.Handle(http.MethodGet, "v1", "/test", tgh.Test)

}
