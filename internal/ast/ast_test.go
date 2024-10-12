package ast

import (
	"testing"
)

func TestASTNodes(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected string
	}{
		{"String", String("test"), `"test"`},
		{"Number (integer)", Number(42), "42"},
		{"Number (float)", Number(3.14), "3.14"},
		{"Boolean (true)", Boolean(true), "true"},
		{"Boolean (false)", Boolean(false), "false"},
		{"Null", Null{}, "null"},
		{"Object", Object{"key": String("value")}, `{"key":"value"}`},
		{"Array", Array{Number(1), String("two"), Boolean(true)}, `[1,"two",true]`},
		{"Empty Object", Object{}, `{}`},
		{"Empty Array", Array{}, `[]`},
		{"Nested Object", Object{"outer": Object{"inner": Number(42)}}, `{"outer":{"inner":42}}`},
		{"Nested Array", Array{Array{Number(1), Number(2)}, Number(3)}, `[[1,2],3]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.node.String()
			if got != tt.expected {
				t.Errorf("%s string representation mismatch: got %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestNodeString(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected string
	}{
		{"String", String("test"), `"test"`},
		{"Number (integer)", Number(42), "42"},
		{"Number (float)", Number(3.14), "3.14"},
		{"Boolean (true)", Boolean(true), "true"},
		{"Boolean (false)", Boolean(false), "false"},
		{"Null", Null{}, "null"},
		{"Object", Object{"key": String("value")}, `{"key":"value"}`},
		{"Array", Array{Number(1), String("two"), Boolean(true)}, `[1,"two",true]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.node.String()
			if result != tt.expected {
				t.Errorf("%s string representation mismatch: got %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}
