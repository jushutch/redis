package serializer

// TypePrefix is the prefix character indicating a RESP type
type TypePrefix rune

// Define type prefixes
const (
	SIMPLE_STRING TypePrefix = '+'
	SIMPLE_ERROR  TypePrefix = '-'
	INTEGER       TypePrefix = ':'
	BULK_STRING   TypePrefix = '$'
	ARRAY         TypePrefix = '*'
)

// String used as RESP terminator
const TERMINATOR = "\r\n"

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
	switch TypePrefix(input[0]) {
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
