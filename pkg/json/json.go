// Package json provides functionality for parsing, manipulating, and serializing JSON data.
package json

import (
	"json-parser/internal/ast"
	"json-parser/internal/lexer"
	"json-parser/internal/parser"
	"json-parser/internal/serializer"
	"strconv"
	"strings"
)

// Parse takes a JSON string and returns an AST Node.
// It returns an error if the input is not valid JSON.
//
// Example:
//
//	node, err := json.Parse(`{"key": "value"}`)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Use the node...
func Parse(input string) (ast.Node, error) {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	return p.Parse()
}

// Stringify takes an AST Node and returns a JSON string.
// It does not return an error as all AST Nodes are considered valid JSON.
//
// Example:
//
//	obj := ast.Object{"key": ast.String("value")}
//	jsonStr := json.Stringify(obj)
//	fmt.Println(jsonStr) // Output: {"key":"value"}
func Stringify(node ast.Node) string {
	return serializer.Serialize(node)
}

// ValidateJSON checks if a given string is valid JSON.
// It returns true if the input is valid JSON, false otherwise.
//
// Example:
//
//	isValid := json.ValidateJSON(`{"key": "value"}`)
//	fmt.Println(isValid) // Output: true
func ValidateJSON(input string) bool {
	_, err := Parse(input)
	return err == nil
}

// PrettyPrint takes a JSON string and returns a formatted JSON string.
// It returns an error if the input is not valid JSON.
//
// Example:
//
//	prettyJSON, err := json.PrettyPrint(`{"key":"value"}`)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(prettyJSON)
//	// Output:
//	// {
//	//   "key": "value"
//	// }
func PrettyPrint(input string) (string, error) {
	node, err := Parse(input)
	if err != nil {
		return "", err
	}
	return formatJSON(node, 0), nil
}

// GetValueByPath takes a JSON string and a path, and returns the value at that path.
// The path is a slice of strings representing the keys to traverse.
// It returns an error if the input is not valid JSON or if the path is invalid.
//
// Example:
//
//	jsonStr := `{"outer": {"inner": 42}}`
//	value, err := json.GetValueByPath(jsonStr, []string{"outer", "inner"})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(value) // Output: 42
func GetValueByPath(input string, path []string) (ast.Node, error) {
	node, err := Parse(input)
	if err != nil {
		return nil, err
	}
	return traversePath(node, path)
}

// traversePath is an internal function that traverses an AST Node using the given path.
func traversePath(node ast.Node, path []string) (ast.Node, error) {
	for _, key := range path {
		switch n := node.(type) {
		case ast.Object:
			var ok bool
			node, ok = n[key]
			if !ok {
				return nil, &KeyNotFoundError{Key: key}
			}
		case ast.Array:
			index, err := strconv.Atoi(key)
			if err != nil {
				return nil, &InvalidPathError{Path: path}
			}
			if index < 0 || index >= len(n) {
				return nil, &IndexOutOfRangeError{Index: index, Length: len(n)}
			}
			node = n[index]
		default:
			return nil, &InvalidPathError{Path: path}
		}
	}
	return node, nil
}

// formatJSON is an internal function that formats an AST Node as a pretty-printed JSON string.
func formatJSON(node ast.Node, indent int) string {
	var sb strings.Builder
	formatJSONToBuilder(&sb, node, indent)
	return sb.String()
}

// formatJSONToBuilder is an internal function that writes a formatted JSON representation of an AST Node to a strings.Builder.
func formatJSONToBuilder(sb *strings.Builder, node ast.Node, indent int) {
	indentStr := strings.Repeat("  ", indent)
	switch n := node.(type) {
	case ast.Object:
		if len(n) == 0 {
			sb.WriteString("{}")
			return
		}
		sb.WriteString("{\n")
		first := true
		for k, v := range n {
			if !first {
				sb.WriteString(",\n")
			}
			first = false
			sb.WriteString(indentStr)
			sb.WriteString("  ")
			sb.WriteString(strconv.Quote(k))
			sb.WriteString(": ")
			formatJSONToBuilder(sb, v, indent+1)
		}
		sb.WriteString("\n")
		sb.WriteString(indentStr)
		sb.WriteString("}")
	case ast.Array:
		if len(n) == 0 {
			sb.WriteString("[]")
			return
		}
		sb.WriteString("[\n")
		first := true
		for _, v := range n {
			if !first {
				sb.WriteString(",\n")
			}
			first = false
			sb.WriteString(indentStr)
			sb.WriteString("  ")
			formatJSONToBuilder(sb, v, indent+1)
		}
		sb.WriteString("\n")
		sb.WriteString(indentStr)
		sb.WriteString("]")
	case ast.String:
		sb.WriteString(strconv.Quote(string(n)))
	case ast.Number:
		sb.WriteString(strconv.FormatFloat(float64(n), 'f', -1, 64))
	case ast.Boolean:
		sb.WriteString(strconv.FormatBool(bool(n)))
	case ast.Null:
		sb.WriteString("null")
	default:
		sb.WriteString("")
	}
}
