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
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": "SALES-API",
	}
	log, err := config.Build()
	if err != nil {
		fmt.Println("Error constructin logger:", err)
		os.Exit(1)
	}

	defer func(log *zap.Logger) {
		err := log.Sync()
		if err != nil {
			fmt.Println("Error syncing logger:", err)
			os.Exit(1)
		}
	}(log)

	// Perform the startup and shutdown sequence
	if err := run(log.Sugar()); err != nil {
		fmt.Println("Error running the app:", err)
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	return nil
}
