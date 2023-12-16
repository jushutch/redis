package serializer

// Define type prefixes
const (
	TERMINATOR    = "\r\n"
	SIMPLE_STRING = '+'
	SIMPLE_ERROR  = '-'
	INTEGER       = ':'
	BULK_STRING   = '$'
	ARRAY         = '*'
)

// RESPType represents a valid RESP data type
type RESPType interface {
	// Deserialize returns the string representation of a RESP data type
	Deserialize() string
}

// Serialize parses a raw request into its valid RESP data type
func Serialize(input string) (RESPType, int) {
	if len(input) == 0 {
		return nil, 0
	}
	switch input[0] {
	case SIMPLE_STRING:
		return serializeSimpleString(input)
	case SIMPLE_ERROR:
		return serializeSimpleError(input)
	case INTEGER:
		return serializeInteger(input)
	case BULK_STRING:
		return serializeBulkString(input)
	case ARRAY:
		return serializeArray(input)
	default:
		return nil, 0
	}
}
