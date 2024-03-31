package serializer

import (
	"fmt"
	"strconv"
)

type SimpleString string

// +OK\r\n
func serializeSimpleString(input string) (SimpleString, int) {
	var str string
	for i := 1; input[i] != '\r'; i++ {
		str += string(input[i])
	}
	return SimpleString(str), len(str) + 3
}

func (s SimpleString) Deserialize() string {
	return fmt.Sprintf("%c%s%s", SIMPLE_STRING, string(s), TERMINATOR)
}

type SimpleError string

// -Error message\r\n
func serializeSimpleError(input string) (SimpleError, int) {
	var err string
	for i := 1; input[i] != '\r'; i++ {
		err += string(input[i])
	}
	return SimpleError(err), len(err) + 3
}

func (s SimpleError) Deserialize() string {
	return fmt.Sprintf("%c%s%s", SIMPLE_ERROR, string(s), TERMINATOR)
}

type Integer struct {
	Prefix bool
	Value  int64
}

// :+1000\r\n
func serializeInteger(input string) (Integer, int) {
	var integer Integer
	var valueStr string
	for i := 1; input[i] != '\r'; i++ {
		if input[i] == '+' {
			integer.Prefix = true
		}
		valueStr += string(input[i])
	}
	if value, err := strconv.ParseInt(valueStr, 10, 0); err == nil {
		integer.Value = value
	}
	return integer, len(valueStr) + 3
}

func (i Integer) Deserialize() string {
	if i.Prefix {
		return fmt.Sprintf("%c%+d%s", INTEGER, i.Value, TERMINATOR)
	}
	return fmt.Sprintf("%c%d%s", INTEGER, i.Value, TERMINATOR)
}

type BulkString struct {
	Length int64
	Value  string
}

// $5\r\nhello\r\n
func serializeBulkString(input string) (BulkString, int) {
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

func (s BulkString) Deserialize() string {
	if s.Length == -1 {
		return fmt.Sprintf("%c%d%s", BULK_STRING, s.Length, TERMINATOR)
	}
	return fmt.Sprintf("%c%d%s%s%s", BULK_STRING, s.Length, TERMINATOR, s.Value, TERMINATOR)
}

type Array struct {
	Length   int64
	Elements []RESPType
}

// *1\r\n$4\r\nping\r\n
func serializeArray(input string) (Array, int) {
	var array Array
	var numElementsStr string
	var i int
	for i = 1; input[i] != '\r'; i++ {
		numElementsStr += string(input[i])
	}
	i += 2
	numElements, _ := strconv.ParseInt(numElementsStr, 10, 0)
	array.Length = numElements
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

func (a Array) Deserialize() string {
	var elementsStr string
	for _, element := range a.Elements {
		elementsStr += element.Deserialize()
	}
	return fmt.Sprintf("%c%d%s%s", ARRAY, a.Length, TERMINATOR, elementsStr)
}
