package ast

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestASTNodes(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected interface{}
	}{
		{"String", String("test"), "test"},
		{"Number (integer)", Number(42), float64(42)},
		{"Number (float)", Number(3.14), 3.14},
		{"Boolean (true)", Boolean(true), true},
		{"Boolean (false)", Boolean(false), false},
		{"Null", Null{}, nil},
		{"Object", Object{"key": String("value")}, map[string]interface{}{"key": "value"}},
		{"Array", Array{Number(1), String("two"), Boolean(true)}, []interface{}{float64(1), "two", true}},
		{"Empty Object", Object{}, map[string]interface{}{}},
		{"Empty Array", Array{}, []interface{}{}},
		{"Nested Object", Object{"outer": Object{"inner": Number(42)}}, map[string]interface{}{"outer": map[string]interface{}{"inner": float64(42)}}},
		{"Nested Array", Array{Array{Number(1), Number(2)}, Number(3)}, []interface{}{[]interface{}{float64(1), float64(2)}, float64(3)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonBytes, err := json.Marshal(tt.node)
			if err != nil {
				t.Fatalf("%s JSON marshaling error: %v", tt.name, err)
			}

			var unmarshaled interface{}
			err = json.Unmarshal(jsonBytes, &unmarshaled)
			if err != nil {
				t.Fatalf("%s JSON unmarshaling error: %v", tt.name, err)
			}

			if !reflect.DeepEqual(unmarshaled, tt.expected) {
				t.Errorf("%s JSON mismatch: got %v, want %v", tt.name, unmarshaled, tt.expected)
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
