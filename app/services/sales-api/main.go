package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/ktruedat/ultimateService/app/services/sales-api/handlers"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

/*
TODO: Need to figure out timeouts for http service
*/

var build = "develop"

func main() {
	//Construct the application logger.
	log, err := initLogger("SALES-API")
	if err != nil {
		fmt.Println("Error constructing the app logger:", err)
		os.Exit(1)

	}
	defer log.Sync()

	// Perform the startup and shutdown sequence
	if err := run(log); err != nil {
		fmt.Println("Error running the app:", err)
		os.Exit(1)
	}
}

// Entry level run function for entering the program
func run(log *zap.SugaredLogger) error {

	// ==================================================================
	// GOMAXPROCS
	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// =================================================================
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s,mask"`
		}
	}{
		Version: conf.Version{
			SVN:  build, // System version number
			Desc: "copyright information here",
		},
	}

	const prefix = "SALES"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// Start Debug Service
	log.Infow("startup", "status", "debug router started", "host", cfg.Web.DebugHost)
	// The Debug function returns a mux to listen and serve on for all the debug
	// related endpoints. This includes the standard library endpoints

	//Construct the mux
	debugMux := handlers.DebugStandardLibraryMux()

	// Start the service listening for debug requests.
	// Not concerned with shutting this down with load shedding.
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			log.Errorw("shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	expvar.NewString("build").Set(build)

	// ==========================================================
	// App starting

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}

	log.Infow("startup", "config", out)

	// ===========================================================
	// Start API Service

	log.Infow("startup", "status", "initializing API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	apiMux := handlers.APIMux(handlers.APIMuxConfig{
		Shutdown: shutdown,
		Log:      log,
	})

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for api requests.
	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// =======================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	// Using blocking select
	select {
	case err := <-serverErrors:
		// if an error occurred here, it must be some low level networking problem
		// that happened, so there is no point in trying to shed load here.
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)
		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func initLogger(service string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": service,
	}
	log, err := config.Build()
	if err != nil {
		fmt.Println("Error constructin logger:", err)
		os.Exit(1)
	}
	return log.Sugar(), nil
}
