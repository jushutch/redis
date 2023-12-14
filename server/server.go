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
	Ping() serializer.RESPType
	Echo(message serializer.BulkString) serializer.RESPType
	Set(key, value serializer.BulkString) serializer.RESPType
	Get(key serializer.BulkString) serializer.RESPType
}

type Server struct {
	manager Manager
}

func New() *Server {
	return &Server{
		manager: manager.NewManager(),
	}
}

func (s *Server) Run(conf Config, logger *slog.Logger) error {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", conf.Host+":"+conf.Port)
	if err != nil {
		return fmt.Errorf("error listening: %w", err)
	}
	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Listening on " + conf.Host + ":" + conf.Port)
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
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	input := string(buf[:n])
	fmt.Printf("\n%q\n", input)
	respType, x := serializer.Serialize(input)
	commands, ok := respType.(serializer.Array)
	if !ok {
		return
	}

	fmt.Printf("\nRESP Type: %q\nBytes Read: %d\n", respType.Deserialize(), x)
	response := s.handleCommands(commands.Elements)
	if response == nil {
		return
	}
	// Send a response back to person contacting us.
	conn.Write([]byte(response.Deserialize()))
}

func (s *Server) handleCommands(commands []serializer.RESPType) serializer.RESPType {
	commandName, ok := commands[0].(serializer.BulkString)
	if !ok {
		return nil
	}
	switch commandName.Value {
	case "PING":
		return s.manager.Ping()
	case "ECHO":
		message, ok := commands[1].(serializer.BulkString)
		if !ok {
			return nil
		}
		return s.manager.Echo(message)
	case "SET":
		key, ok := commands[1].(serializer.BulkString)
		if !ok {
			return nil
		}
		value, ok := commands[2].(serializer.BulkString)
		if !ok {
			return nil
		}
		return s.manager.Set(key, value)
	case "GET":
		key, ok := commands[1].(serializer.BulkString)
		if !ok {
			return nil
		}
		return s.manager.Get(key)
	}
	return nil
}
