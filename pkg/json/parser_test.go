package json

import (
	"errors"
	"json-parser/internal/ast"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ast.Node
		wantErr  bool
	}{
		{"Parse simple object", `{"key": "value"}`, ast.Object{"key": ast.String("value")}, false},
		{"Parse simple array", `[1, 2, 3]`, ast.Array{ast.Number(1), ast.Number(2), ast.Number(3)}, false},
		{"Parse nested object", `{"outer": {"inner": 42}}`, ast.Object{"outer": ast.Object{"inner": ast.Number(42)}}, false},
		{"Parse nested array", `[1, [2, 3], 4]`, ast.Array{ast.Number(1), ast.Array{ast.Number(2), ast.Number(3)}, ast.Number(4)}, false},
		{"Parse empty object", `{}`, ast.Object{}, false},
		{"Parse empty array", `[]`, ast.Array{}, false},
		{"Parse null", `null`, ast.Null{}, false},
		{"Parse true", `true`, ast.Boolean(true), false},
		{"Parse false", `false`, ast.Boolean(false), false},
		{"Parse string with escapes", `"Hello,\n\t\"World\"!"`, ast.String("Hello,\n\t\"World\"!"), false},
		{"Parse number", `123.45`, ast.Number(123.45), false},
		{"Parse negative number", `-42`, ast.Number(-42), false},
		{"Parse exponential number", `1.23e-4`, ast.Number(1.23e-4), false},
		{"Parse invalid JSON", `{"key": "value",}`, nil, true},
		{"Parse unclosed object", `{"key": "value"`, nil, true},
		{"Parse unclosed array", `[1, 2, 3`, nil, true},
		{"Parse invalid number", `123.`, nil, true},
		{"Parse invalid escape", `"\u123G"`, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Parse() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStringify(t *testing.T) {
	tests := []struct {
		name     string
		input    ast.Node
		expected string
	}{
		{"Stringify simple object", ast.Object{"key": ast.String("value")}, `{"key":"value"}`},
		{"Stringify simple array", ast.Array{ast.Number(1), ast.Number(2), ast.Number(3)}, `[1,2,3]`},
		{"Stringify nested object", ast.Object{"outer": ast.Object{"inner": ast.Number(42)}}, `{"outer":{"inner":42}}`},
		{"Stringify nested array", ast.Array{ast.Number(1), ast.Array{ast.Number(2), ast.Number(3)}, ast.Number(4)}, `[1,[2,3],4]`},
		{"Stringify null", ast.Null{}, `null`},
		{"Stringify boolean", ast.Boolean(true), `true`},
		{"Stringify string", ast.String("Hello, World!"), `"Hello, World!"`},
		{"Stringify number", ast.Number(123.45), `123.45`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Stringify(tt.input); got != tt.expected {
				t.Errorf("Stringify() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid simple object", `{"key": "value"}`, true},
		{"Valid simple array", `[1, 2, 3]`, true},
		{"Valid complex JSON", `{"key1": [1, 2, {"key2": null}], "key3": true}`, true},
		{"Invalid JSON (trailing comma)", `{"key": "value",}`, false},
		{"Invalid JSON (unclosed object)", `{"key": "value"`, false},
		{"Invalid JSON (unclosed array)", `[1, 2, 3`, false},
		{"Invalid JSON (invalid number)", `{"key": 123.}`, false},
		{"Invalid JSON (invalid escape)", `{"key": "\u123G"}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateJSON(tt.input); got != tt.expected {
				t.Errorf("ValidateJSON() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPrettyPrint(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:  "Pretty print simple object",
			input: `{"key":"value"}`,
			expected: `{
  "key": "value"
}`,
			wantErr: false,
		},
		{
			name:  "Pretty print nested object",
			input: `{"outer":{"inner":42}}`,
			expected: `{
  "outer": {
    "inner": 42
  }
}`,
			wantErr: false,
		},
		{
			name:  "Pretty print array",
			input: `[1,2,3]`,
			expected: `[
  1,
  2,
  3
]`,
			wantErr: false,
		},
		{
			name:    "Pretty print invalid JSON",
			input:   `{"key": "value",}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrettyPrint(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrettyPrint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("PrettyPrint() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetValueByPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		path     []string
		expected ast.Node
		wantErr  error
	}{
		{"Get value from simple object", `{"key": "value"}`, []string{"key"}, ast.String("value"), nil},
		{"Get value from nested object", `{"outer": {"inner": 42}}`, []string{"outer", "inner"}, ast.Number(42), nil},
		{"Get value from array", `[1, 2, 3]`, []string{"1"}, ast.Number(2), nil},
		{"Get value from non-existent path", `{"key": "value"}`, []string{"nonexistent"}, nil, &KeyNotFoundError{Key: "nonexistent"}},
		{"Get value from array with invalid index", `[1, 2, 3]`, []string{"3"}, nil, &IndexOutOfRangeError{Index: 3, Length: 3}},
		{"Invalid path", `{"key": "value"}`, []string{"key", "invalid"}, nil, &InvalidPathError{Path: []string{"key", "invalid"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetValueByPath(tt.input, tt.path)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetValueByPath() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetValueByPath() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("GetValueByPath() unexpected error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GetValueByPath() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParse_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Empty input", "", true},
		{"Only whitespace", "   \t\n", true},
		{"Very large number", "1e1000", true},
		{"Very small number", "1e-1000", true},
		{"Number with many decimal places", "0." + strings.Repeat("1", 1000), false},
		{"Deeply nested object", "{" + strings.Repeat("\"a\":{", 100) + "\"b\":1" + strings.Repeat("}", 100), false},
		{"Deeply nested array", "[" + strings.Repeat("[", 100) + "1" + strings.Repeat("]", 100), false},
		{"Unicode escape", `"\u0041\u0042\u0043"`, false},
		{"Invalid Unicode escape", `"\u123G"`, true},
		{"Incomplete escape", `"\`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	input := `{"key1": "value1", "key2": 42, "key3": [1, 2, 3], "key4": {"nested": true}}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStringify(b *testing.B) {
	input := ast.Object{
		"key1": ast.String("value1"),
		"key2": ast.Number(42),
		"key3": ast.Array{ast.Number(1), ast.Number(2), ast.Number(3)},
		"key4": ast.Object{"nested": ast.Boolean(true)},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Stringify(input)
	}
}

func BenchmarkPrettyPrint(b *testing.B) {
	input := `{"key1":"value1","key2":42,"key3":[1,2,3],"key4":{"nested":true}}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := PrettyPrint(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}
