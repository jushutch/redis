package server

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/jushutch/redis/manager"
	"github.com/jushutch/redis/serializer"
)

// Config contains configuration parameters for the server
type Config struct {
	Host string
	Port string
}

// Manager handles Redis commands
type Manager interface {
	HandleCommand(command serializer.Array) serializer.RESPType
}

// Server listens for RESP requests and sends RESP responses
type Server struct {
	manager Manager
	logger  *slog.Logger
}

// New creates a new server
func New(logger *slog.Logger) *Server {
	return &Server{
		manager: manager.New(logger),
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
		go s.handleRequest(conn)
	}
}

// handleRequest serializes the request into commands and executes them
func (s *Server) handleRequest(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		if errors.Is(err, io.EOF) {
			s.logger.Error("failed to read from connection", "error", err)
		}
		return
	}
	request := string(buffer[:bytesRead])
	s.logger.Info("received request", "request", request)

	// Execute all commands in request and send responses
	for i := 0; i < bytesRead; {
		commandType, bytes := serializer.Serialize(request[i:])
		i += bytes
		logger := s.logger.With(
			"bytes", bytes,
			"request", request,
			"command", request[i-bytes:i],
		)
		logger.Info("serialized command")
		command, ok := commandType.(serializer.Array)
		if !ok {
			logger.Error("command was not an array")
			continue
		}
		response := s.manager.HandleCommand(command)
		if response == nil {
			logger.Error("response was nil")
			continue
		}
		logger = logger.With("response", response.Deserialize())
		logger.Info("send response")
		bytesWritten, err := conn.Write([]byte(response.Deserialize()))
		if err != nil {
			s.logger.Error("failed to write response")
			continue
		}
		s.logger.Info("sent response", "bytes_written", bytesWritten)
	}
}
