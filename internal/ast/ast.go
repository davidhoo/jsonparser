package ast

import (
	"fmt"
	"strconv"
	"strings"
)

// Node represents a node in the JSON AST.
type Node interface {
	String() string
}

// String represents a JSON string value.
type String string

// String returns the string representation of the JSON string.
func (s String) String() string {
	return fmt.Sprintf("%q", string(s)) // 使用 %q 来返回带引号的字符串
}

// Number represents a JSON number value.
type Number float64

// String returns the string representation of the JSON number.
func (n Number) String() string {
	return fmt.Sprintf("%g", float64(n))
}

// Boolean represents a JSON boolean value.
type Boolean bool

// String returns the string representation of the JSON boolean.
func (b Boolean) String() string {
	return strconv.FormatBool(bool(b))
}

// Null represents a JSON null value.
type Null struct{}

// String returns the string representation of the JSON null value.
func (n Null) String() string {
	return "null"
}

// Object represents a JSON object.
type Object map[string]Node

// String returns the string representation of the JSON object.
func (o Object) String() string {
	pairs := make([]string, 0, len(o))
	for k, v := range o {
		pairs = append(pairs, fmt.Sprintf("%q:%s", k, v.String()))
	}
	return "{" + strings.Join(pairs, ",") + "}"
}

// Array represents a JSON array.
type Array []Node

// String returns the string representation of the JSON array.
func (a Array) String() string {
	elements := make([]string, len(a))
	for i, v := range a {
		elements[i] = v.String()
	}
	return "[" + strings.Join(elements, ",") + "]"
}
