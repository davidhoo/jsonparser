package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type Node interface {
	String() string
	PrettyString(indent string) string
	Query(query string) (Node, error)
	Search(value string) []Node
}

type Object map[string]Node
type Array []Node
type String string
type Number float64
type Boolean bool
type Null struct{}

// Implement String() method for each type
func (o Object) String() string {
	pairs := make([]string, 0, len(o))
	for k, v := range o {
		pairs = append(pairs, fmt.Sprintf("%q: %v", k, v))
	}
	return "{" + strings.Join(pairs, ", ") + "}"
}

func (a Array) String() string {
	elements := make([]string, len(a))
	for i, v := range a {
		elements[i] = v.String()
	}
	return "[" + strings.Join(elements, ", ") + "]"
}

func (s String) String() string {
	return fmt.Sprintf("%q", string(s))
}

func (n Number) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

func (b Boolean) String() string {
	return strconv.FormatBool(bool(b))
}

func (Null) String() string {
	return "null"
}

// Implement Node interface for each type...

func (o Object) Query(query string) (Node, error) {
	parts := strings.Split(query, "/")
	if len(parts) == 0 {
		return o, nil
	}
	if parts[0] == "" {
		parts = parts[1:]
	}
	if len(parts) == 0 {
		return o, nil
	}

	key := parts[0]
	value, ok := o[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	if len(parts) == 1 {
		return value, nil
	}

	return value.Query(strings.Join(parts[1:], "/"))
}

func (a Array) Query(query string) (Node, error) {
	parts := strings.Split(query, "/")
	if len(parts) == 0 {
		return a, nil
	}
	if parts[0] == "" {
		parts = parts[1:]
	}
	if len(parts) == 0 {
		return a, nil
	}

	if parts[0] == "*" {
		results := make(Array, 0)
		for _, item := range a {
			result, err := item.Query(strings.Join(parts[1:], "/"))
			if err == nil {
				results = append(results, result)
			}
		}
		return results, nil
	}

	index, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid array index: %s", parts[0])
	}

	if index < 0 || index >= len(a) {
		return nil, fmt.Errorf("array index out of bounds: %d", index)
	}

	if len(parts) == 1 {
		return a[index], nil
	}

	return a[index].Query(strings.Join(parts[1:], "/"))
}

// Implement Search methods for each type...

func (o Object) Search(value string) []Node {
	return nil
}

func (a Array) Search(value string) []Node {
	return nil
}

func (s String) Search(value string) []Node {
	return nil
}

func (n Number) Search(value string) []Node {
	return nil
}

func (b Boolean) Search(value string) []Node {
	return nil
}

func (Null) Search(value string) []Node {
	return nil
}

func (s String) Query(query string) (Node, error) {
	if query == "" {
		return s, nil
	}
	return nil, fmt.Errorf("cannot query string")
}

func (n Number) Query(query string) (Node, error) {
	if query == "" {
		return n, nil
	}
	return nil, fmt.Errorf("cannot query number")
}

func (b Boolean) Query(query string) (Node, error) {
	if query == "" {
		return b, nil
	}
	return nil, fmt.Errorf("cannot query boolean")
}

func (Null) Query(query string) (Node, error) {
	if query == "" {
		return Null{}, nil
	}
	return nil, fmt.Errorf("cannot query null")
}

// New PrettyString methods for each type
func (o Object) PrettyString(indent string) string {
	if len(o) == 0 {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteString("{\n")
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}
	for i, k := range keys {
		v := o[k]
		sb.WriteString(indent + "  ")
		sb.WriteString(color.GreenString("%s", strconv.Quote(k)))
		sb.WriteString(color.WhiteString(": "))
		sb.WriteString(v.PrettyString(indent + "  "))
		if i < len(keys)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString(indent + "}")
	return sb.String()
}

func (a Array) PrettyString(indent string) string {
	if len(a) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, v := range a {
		sb.WriteString(indent + "  ")
		sb.WriteString(v.PrettyString(indent + "  "))
		if i < len(a)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString(indent + "]")
	return sb.String()
}

func (s String) PrettyString(indent string) string {
	return color.BlueString("%s", strconv.Quote(string(s)))
}

func (n Number) PrettyString(indent string) string {
	return color.CyanString("%v", float64(n))
}

func (b Boolean) PrettyString(indent string) string {
	return color.YellowString("%v", bool(b))
}

func (Null) PrettyString(indent string) string {
	return color.MagentaString("null")
}

// ... (rest of the file remains unchanged)
