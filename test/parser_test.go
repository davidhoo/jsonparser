package test

import (
	"json-parser/pkg/json"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{"name": "Alice", "age": 30}`, `{"name": "Alice", "age": 30}`},
		{`[1, 2, 3]`, `[1, 2, 3]`},
		{`true`, `true`},
		{`false`, `false`},
		{`null`, `null`},
		{`42`, `42`},
		{`"Hello, World!"`, `"Hello, World!"`},
	}

	for _, tt := range tests {
		result, err := json.Parse(tt.input)
		if err != nil {
			t.Errorf("Parse(%q) returned error: %v", tt.input, err)
			continue
		}

		// 将结果转换为字符串进行比较
		resultStr := result.String()
		if resultStr != tt.expected {
			t.Errorf("Parse(%q) = %q, want %q", tt.input, resultStr, tt.expected)
		}
	}
}

func TestParseError(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`{"name": "Alice", "age": }`},
		{`[1, 2, 3`},
		{`{"key": "value"`},
		{`truee`},
	}

	for _, tt := range tests {
		_, err := json.Parse(tt.input)
		if err == nil {
			t.Errorf("Parse(%q) did not return an error, want error", tt.input)
		}
	}
}
