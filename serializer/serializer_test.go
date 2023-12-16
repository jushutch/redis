package serializer

import "testing"

func TestSerialize(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"+OK\r\n"},
		{"+hello world\r\n"},
		{"-Error message\r\n"},
		{":+10000\r\n"},
		{":-10000\r\n"},
		{":10000\r\n"},
		{":0\r\n"},
		{"$-1\r\n"},
		{"$5\r\nhello\r\n"},
		{"$0\r\n\r\n"},
		{"*1\r\n$4\r\nping\r\n"},
		{"*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n"},
		{"*2\r\n$3\r\nget\r\n$3\r\nkey\r\n"},
		{"*-1\r\n"},
		{"*0\r\n"},
		{"*3\r\n$5\r\nhello\r\n$-1\r\n$5\r\nworld\r\n"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output, length := Serialize(test.input)
			if output == nil {
				t.Fatalf("tests[%q] - serialization output was unexpectedly nil", test.input)
			}
			if length != len(test.input) {
				t.Fatalf("tests[%q] - expected length to match input: got %d, wanted %d", test.input, length, len(test.input))
			}
			if output.Deserialize() != test.input {
				t.Fatalf("tests[%q] - got %s, wanted %s", test.input, output.Deserialize(), test.input)
			}
		})
	}
}
