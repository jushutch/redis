package server

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/jushutch/redis/manager"
	"github.com/jushutch/redis/serializer"
)

type Config struct {
	Host string
	Port string
}

type Manager interface {
	HandleCommand(command serializer.Array) serializer.RESPType
}

type Server struct {
	manager Manager
	logger  *slog.Logger
}

func New(logger *slog.Logger) *Server {
	return &Server{
		manager: manager.NewManager(logger),
		logger:  logger.With("name", "redis.server"),
	}
}

func (s *Server) Run(conf Config) error {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", conf.Host+":"+conf.Port)
	if err != nil {
		return fmt.Errorf("error listening: %w", err)
	}
	// Close the listener when the application closes.
	defer l.Close()

	s.logger.Info("listening", "host", conf.Host, "port", conf.Port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("error accepting: %w", err)
		}
		// Handle connections in a new goroutine.
		go s.handleRequest(conn)
	}
}

func (s *Server) handleRequest(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		s.logger.Error("failed to read from connection", "error", err)
		return
	}
	request := string(buf[:n])
	s.logger.Info("received request", "request", request)
	requestType, _ := serializer.Serialize(request)
	command, ok := requestType.(serializer.Array)
	if !ok {
		s.logger.Error("request was not an array", "request", request)
		return
	}
	response := s.manager.HandleCommand(command)
	if response == nil {
		s.logger.Error("response was nil", "request", request)
		return
	}
	s.logger.Info("send response", "response", response.Deserialize())
	conn.Write([]byte(response.Deserialize()))
}
