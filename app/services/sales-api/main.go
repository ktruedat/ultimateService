package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

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

// Entry level run function
func run(log *zap.SugaredLogger) error {

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
