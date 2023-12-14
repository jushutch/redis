package serializer

import (
	"fmt"
	"strconv"
)

const (
	TERMINATOR    = "\r\n"
	SIMPLE_STRING = '+'
	SIMPLE_ERROR  = '-'
	INTEGER       = ':'
	BULK_STRING   = '$'
	ARRAY         = '*'
)

type RESPType interface {
	Deserialize() string
}

type SimpleString string

func (s SimpleString) Deserialize() string {
	return fmt.Sprintf("%c%s%s", SIMPLE_STRING, string(s), TERMINATOR)
}

type SimpleError string

func (s SimpleError) Deserialize() string {
	return fmt.Sprintf("%c%s%s", SIMPLE_ERROR, string(s), TERMINATOR)
}

type Integer struct {
	prefix bool
	value  int64
}

func (i Integer) Deserialize() string {
	if i.prefix {
		return fmt.Sprintf("%c%+d%s", INTEGER, i.value, TERMINATOR)
	}
	return fmt.Sprintf("%c%d%s", INTEGER, i.value, TERMINATOR)
}

type BulkString struct {
	Length int64
	Value  string
}

func (s BulkString) Deserialize() string {
	if s.Length == -1 {
		return fmt.Sprintf("%c%d%s", BULK_STRING, s.Length, TERMINATOR)
	}
	return fmt.Sprintf("%c%d%s%s%s", BULK_STRING, s.Length, TERMINATOR, s.Value, TERMINATOR)
}

type Array struct {
	length   int64
	Elements []RESPType
}

func (a Array) Deserialize() string {
	var elementsStr string
	for _, element := range a.Elements {
		elementsStr += element.Deserialize()
	}
	return fmt.Sprintf("%c%d%s%s", ARRAY, a.length, TERMINATOR, elementsStr)
}

func Serialize(input string) (RESPType, int) {
	if len(input) == 0 {
		return nil, 0
	}
	firstByte := input[0]
	switch firstByte {
	case SIMPLE_STRING:
		return SerializeSimpleString(input)
	case SIMPLE_ERROR:
		return SerializeSimpleError(input)
	case INTEGER:
		return SerializeInteger(input)
	case BULK_STRING:
		return SerializeBulkString(input)
	case ARRAY:
		return SerializeArray(input)
	default:
		return nil, 0
	}
}

// +OK\r\n
func SerializeSimpleString(input string) (SimpleString, int) {
	var str string
	for i := 1; input[i] != '\r'; i++ {
		str += string(input[i])
	}
	return SimpleString(str), len(str) + 3
}

// -Error message\r\n
func SerializeSimpleError(input string) (SimpleError, int) {
	var err string
	for i := 1; input[i] != '\r'; i++ {
		err += string(input[i])
	}
	return SimpleError(err), len(err) + 3
}

// :+1000\r\n
func SerializeInteger(input string) (Integer, int) {
	var integer Integer
	var valueStr string
	for i := 1; input[i] != '\r'; i++ {
		if input[i] == '+' {
			integer.prefix = true
		}
		valueStr += string(input[i])
	}
	if value, err := strconv.ParseInt(valueStr, 10, 0); err == nil {
		integer.value = value
	}
	return integer, len(valueStr) + 3
}

// $5\r\nhello\r\n
func SerializeBulkString(input string) (BulkString, int) {
	var bulkString BulkString
	var lengthStr string
	var value string
	var i int
	for i = 1; input[i] != '\r'; i++ {
		lengthStr += string(input[i])
	}
	length, _ := strconv.ParseInt(lengthStr, 10, 0)
	bulkString.Length = length
	i += 2
	if length == -1 {
		return bulkString, i
	}
	for ; input[i] != '\r'; i++ {
		value += string(input[i])
	}
	bulkString.Value = value
	return bulkString, i + 2
}

// *1\r\n$4\r\nping\r\n
func SerializeArray(input string) (Array, int) {
	var array Array
	var numElementsStr string
	var i int
	for i = 1; input[i] != '\r'; i++ {
		numElementsStr += string(input[i])
	}
	i += 2
	numElements, _ := strconv.ParseInt(numElementsStr, 10, 0)
	array.length = numElements
	if numElements == -1 {
		return array, i
	}
	elements := make([]RESPType, 0, numElements)

	for k := 0; k < int(numElements); k++ {
		nextElement, length := Serialize(input[i:])
		elements = append(elements, nextElement)
		i += length
	}
	array.Elements = elements
	return array, i
}
