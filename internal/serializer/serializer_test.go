package serializer

import (
	"json-parser/internal/ast"
	"testing"
)

func TestSerialize(t *testing.T) {
	tests := []struct {
		name     string
		input    ast.Node
		expected string
	}{
		{
			name:     "Serialize string",
			input:    ast.String("test"),
			expected: `"test"`,
		},
		{
			name:     "Serialize number",
			input:    ast.Number(42.5),
			expected: `42.5`,
		},
		{
			name:     "Serialize boolean",
			input:    ast.Boolean(true),
			expected: `true`,
		},
		{
			name:     "Serialize null",
			input:    ast.Null{},
			expected: `null`,
		},
		{
			name: "Serialize object",
			input: ast.Object{
				"key1": ast.String("value1"),
				"key2": ast.Number(42),
			},
			expected: `{"key1":"value1","key2":42}`,
		},
		{
			name: "Serialize array",
			input: ast.Array{
				ast.String("item1"),
				ast.Number(2),
				ast.Boolean(true),
			},
			expected: `["item1",2,true]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Serialize(tt.input)
			if result != tt.expected {
				t.Errorf("Serialize() = %v, want %v", result, tt.expected)
			}
		})
	}
}
