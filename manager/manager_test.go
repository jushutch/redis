package manager

import (
	"log/slog"
	"os"
	"testing"

	"github.com/jushutch/redis/serializer"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -destination=repo_mock.go -package=manager github.com/jushutch/redis/manager Repo

func TestManager(t *testing.T) {
	tests := []struct {
		name     string
		input    serializer.Array
		expected serializer.RESPType
		expect   func(*MockRepo)
	}{
		{
			name: "ping",
			input: serializer.Array{
				Length:   1,
				Elements: []serializer.RESPType{serializer.BulkString{Length: 4, Value: "PING"}},
			},
			expected: serializer.SimpleString("PONG"),
			expect:   func(_ *MockRepo) {},
		},
		{
			name: "echo",
			input: serializer.Array{
				Length: 2,
				Elements: []serializer.RESPType{
					serializer.BulkString{Length: 4, Value: "ECHO"},
					serializer.BulkString{Length: 12, Value: "Hello world!"},
				},
			},
			expected: serializer.BulkString{Length: 12, Value: "Hello world!"},
			expect:   func(_ *MockRepo) {},
		},
		{
			name: "set",
			input: serializer.Array{
				Length: 2,
				Elements: []serializer.RESPType{
					serializer.BulkString{Length: 3, Value: "SET"},
					serializer.BulkString{Length: 3, Value: "key"},
					serializer.BulkString{Length: 5, Value: "value"},
				},
			},
			expected: serializer.SimpleString("OK"),
			expect: func(m *MockRepo) {
				m.EXPECT().Set("key", "value").Return(nil).Times(1)
			},
		},
		{
			name: "get",
			input: serializer.Array{
				Length: 2,
				Elements: []serializer.RESPType{
					serializer.BulkString{Length: 3, Value: "GET"},
					serializer.BulkString{Length: 3, Value: "key"},
				},
			},
			expected: serializer.BulkString{Length: 5, Value: "value"},
			expect: func(m *MockRepo) {
				m.EXPECT().Get("key").Return("value", nil).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
			logger := slog.New(handler)
			mockRepo := NewMockRepo(ctrl)
			manager := New(mockRepo, logger)
			test.expect(mockRepo)
			actual := manager.HandleCommand(test.input)
			if actual != test.expected {
				t.Fatalf("failed to handle command %q. got %v, wanted %v",
					test.input.Deserialize(),
					actual,
					test.expected,
				)
			}
		})
	}
}
