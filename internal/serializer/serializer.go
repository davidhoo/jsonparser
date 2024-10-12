package serializer

import (
	"fmt"
	"json-parser/internal/ast"
	"strconv"
	"strings"
)

// Serialize converts an AST node to its JSON string representation.
// It handles all types of AST nodes: Object, Array, String, Number, Boolean, and Null.
func Serialize(node ast.Node) string {
	switch n := node.(type) {
	case ast.Object:
		return serializeObject(n)
	case ast.Array:
		return serializeArray(n)
	case ast.String:
		return serializeString(string(n))
	case ast.Number:
		return serializeNumber(float64(n))
	case ast.Boolean:
		return serializeBoolean(bool(n))
	case ast.Null:
		return "null"
	default:
		return "" // This should never happen if all AST node types are handled
	}
}

// serializeObject converts an AST Object to its JSON string representation.
func serializeObject(obj ast.Object) string {
	if len(obj) == 0 {
		return "{}"
	}

	var pairs []string
	for k, v := range obj {
		pairs = append(pairs, fmt.Sprintf("%s:%s", serializeString(k), Serialize(v)))
	}
	return "{" + strings.Join(pairs, ",") + "}"
}

// serializeArray converts an AST Array to its JSON string representation.
func serializeArray(arr ast.Array) string {
	if len(arr) == 0 {
		return "[]"
	}

	var elements []string
	for _, v := range arr {
		elements = append(elements, Serialize(v))
	}
	return "[" + strings.Join(elements, ",") + "]"
}

// serializeString converts a Go string to its JSON string representation.
// It handles escaping special characters as per JSON specification.
func serializeString(s string) string {
	return strconv.Quote(s)
}

// serializeNumber converts a Go float64 to its JSON number representation.
func serializeNumber(n float64) string {
	return strconv.FormatFloat(n, 'g', -1, 64)
}

// serializeBoolean converts a Go bool to its JSON boolean representation.
func serializeBoolean(b bool) string {
	return strconv.FormatBool(b)
}
