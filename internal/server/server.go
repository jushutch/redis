package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/gofrs/uuid"
	"github.com/jushutch/redis/internal/logging"
	"github.com/jushutch/redis/internal/manager"
	"github.com/jushutch/redis/internal/repo"
	"github.com/jushutch/redis/internal/serializer"
	"github.com/jushutch/redis/internal/tracing"
)

// Config contains configuration parameters for the server
type Config struct {
	Host string
	Port string
}

// Manager handles Redis commands
type Manager interface {
	HandleCommand(ctx context.Context, command serializer.Array) serializer.RESPType
}

// Server listens for RESP requests and sends RESP responses
type Server struct {
	manager Manager
	logger  *slog.Logger
}

// New creates a new server
func New(logger *slog.Logger) *Server {
	repo := repo.New(logger)
	return &Server{
		manager: manager.New(repo, logger),
		logger:  logger.With("name", "redis.server"),
	}
}

// Run starts the server and beings listening for requests
func (s *Server) Run(conf Config) error {
	l, err := net.Listen("tcp", conf.Host+":"+conf.Port)
	if err != nil {
		return fmt.Errorf("error listening: %w", err)
	}
	defer l.Close()
	s.logger.Info("listening", "host", conf.Host, "port", conf.Port)

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("error accepting: %w", err)
		}
		traceID := uuid.Must(uuid.NewV4()).String()
		ctx := context.Background()
		ctx = tracing.SetTraceID(ctx, traceID)
		go s.handleRequest(ctx, conn)
	}
}

// handleRequest serializes the request into commands and executes them
func (s *Server) handleRequest(ctx context.Context, conn net.Conn) {
	logger := s.logger.With(logging.FieldsFromContext(ctx)...)
	defer conn.Close()
	buffer := make([]byte, 1024)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			logger.Error("failed to read from connection", "error", err)
		}
		return
	}
	request := string(buffer[:bytesRead])
	logger.Info("received request", "request", request)
	// Execute all commands in request and send responses
	for i := 0; i < bytesRead; {
		commandType, bytes := serializer.Serialize(request[i:])
		i += bytes
		cmdLogger := logger.With(
			"bytes", bytes,
			"request", request,
			"command", request[i-bytes:i],
		)
		cmdLogger.Info("serialized command")
		command, ok := commandType.(serializer.Array)
		if !ok {
			cmdLogger.Error("command was not an array")
			continue
		}
		response := s.manager.HandleCommand(ctx, command)
		if response == nil {
			cmdLogger.Error("response was nil")
			continue
		}
		if val, ok := response.(serializer.SimpleError); ok {
			cmdLogger.Error(string(val))
		}
		cmdLogger = cmdLogger.With("response", response.Deserialize())
		cmdLogger.Info("send response")
		bytesWritten, err := conn.Write([]byte(response.Deserialize()))
		if err != nil {
			cmdLogger.Error("failed to write response")
			continue
		}
		cmdLogger.Info("sent response", "bytes_written", bytesWritten)
	}
}
