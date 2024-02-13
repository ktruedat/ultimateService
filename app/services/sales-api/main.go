package main

import (
	"errors"
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

	// ================================================================
	// GOMAXPROCS

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// ===============================================================
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

	// ==========================================================
	// App starting

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}

	log.Infow("startup", "config", out)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

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
