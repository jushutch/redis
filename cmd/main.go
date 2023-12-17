package main

import (
	"log/slog"
	"os"

	"github.com/jushutch/redis/server"
	"github.com/namsral/flag"
)

func main() {
	// Parse configuration from command line/environment variables
	var conf server.Config
	var debug bool
	flag.StringVar(&conf.Host, "host", "localhost", "host to run the Redis server on")
	flag.StringVar(&conf.Port, "port", "6379", "port to run the Redis server on")
	flag.BoolVar(&debug, "debug", false, "set logging level to debug")
	flag.Parse()

	// Create logger
	logOptions := &slog.HandlerOptions{Level: slog.LevelWarn}
	if debug {
		logOptions.Level = slog.LevelDebug
		logOptions.AddSource = true
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, logOptions))

	// Run service
	if err := server.New(logger).Run(conf); err != nil {
		logger.Error("failed to run server", "error", err)
		os.Exit(1)
	}
}
