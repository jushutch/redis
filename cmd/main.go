package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jushutch/redis/server"
	"github.com/namsral/flag"
)

func main() {
	var conf server.Config
	flag.StringVar(&conf.Host, "host", "localhost", "host to run the Redis server on")
	flag.StringVar(&conf.Port, "port", "6379", "port to run the Redis server on")
	debug := flag.Bool("debug", false, "set logging level to debug")
	flag.Parse()
	logOptions := &slog.HandlerOptions{Level: slog.LevelInfo}
	if *debug {
		logOptions.Level = slog.LevelDebug
	}
	handler := slog.NewJSONHandler(os.Stdout, logOptions)
	logger := slog.New(handler)
	if err := server.New().Run(conf, logger); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
