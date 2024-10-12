package parser

import (
	"json-parser/internal/ast"
	"json-parser/internal/lexer"
	"reflect"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ast.Node
		wantErr  bool
	}{
		{
			name:     "Parse null",
			input:    "null",
			expected: ast.Null{},
			wantErr:  false,
		},
		{
			name:     "Parse true",
			input:    "true",
			expected: ast.Boolean(true),
			wantErr:  false,
		},
		{
			name:     "Parse false",
			input:    "false",
			expected: ast.Boolean(false),
			wantErr:  false,
		},
		{
			name:     "Parse number",
			input:    "123.45",
			expected: ast.Number(123.45),
			wantErr:  false,
		},
		{
			name:     "Parse string",
			input:    `"hello"`,
			expected: ast.String("hello"),
			wantErr:  false,
		},
		{
			name:  "Parse simple object",
			input: `{"key": "value"}`,
			expected: ast.Object{
				"key": ast.String("value"),
			},
			wantErr: false,
		},
		{
			name:     "Parse simple array",
			input:    `[1, 2, 3]`,
			expected: ast.Array{ast.Number(1), ast.Number(2), ast.Number(3)},
			wantErr:  false,
		},
		{
			name:    "Parse invalid input",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:  "Parse nested object",
			input: `{"outer": {"inner": 42}}`,
			expected: ast.Object{
				"outer": ast.Object{
					"inner": ast.Number(42),
				},
			},
			wantErr: false,
		},
		{
			name:     "Parse nested array",
			input:    `[1, [2, 3], 4]`,
			expected: ast.Array{ast.Number(1), ast.Array{ast.Number(2), ast.Number(3)}, ast.Number(4)},
			wantErr:  false,
		},
		{
			name:     "Parse empty object",
			input:    `{}`,
			expected: ast.Object{},
			wantErr:  false,
		},
		{
			name:     "Parse empty array",
			input:    `[]`,
			expected: ast.Array{},
			wantErr:  false,
		},
		{
			name: "Parse complex object",
			input: `{
				"string": "value",
				"number": 42,
				"bool": true,
				"null": null,
				"array": [1, "two", false],
				"object": {"nested": "ok"}
			}`,
			expected: ast.Object{
				"string": ast.String("value"),
				"number": ast.Number(42),
				"bool":   ast.Boolean(true),
				"null":   ast.Null{},
				"array":  ast.Array{ast.Number(1), ast.String("two"), ast.Boolean(false)},
				"object": ast.Object{"nested": ast.String("ok")},
			},
			wantErr: false,
		},
		{
			name:     "Parse special numbers",
			input:    `[1e10, -3.14, 0.123]`,
			expected: ast.Array{ast.Number(1e10), ast.Number(-3.14), ast.Number(0.123)},
			wantErr:  false,
		},
		{
			name:     "Parse string with escapes",
			input:    `"Hello,\n\t\"World\"!"`,
			expected: ast.String("Hello,\n\t\"World\"!"),
			wantErr:  false,
		},
		{
			name:    "Error: Unclosed object",
			input:   `{"key": "value"`,
			wantErr: true,
		},
		{
			name:    "Error: Unclosed array",
			input:   `[1, 2, 3`,
			wantErr: true,
		},
		{
			name:    "Error: Missing colon in object",
			input:   `{"key" "value"}`,
			wantErr: true,
		},
		{
			name:    "Error: Trailing comma",
			input:   `[1, 2, 3,]`,
			wantErr: true,
		},
		{
			name:    "Error: Invalid Unicode escape",
			input:   `"\u123G"`,
			wantErr: true,
		},
		{
			name:    "Error: Unexpected end of input in string",
			input:   `"unclosed string`,
			wantErr: true,
		},
		{
			name:    "Error: Invalid number format",
			input:   `123.`,
			wantErr: true,
		},
		{
			name:     "Parse very large number",
			input:    `1e308`,
			expected: ast.Number(1e308),
			wantErr:  false,
		},
		{
			name:    "Error: Number too large",
			input:   `1e309`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := NewParser(l)
			got, err := p.Parse()

			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Parser.Parse() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParser_parseObject(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ast.Object
		wantErr  bool
	}{
		{
			name:     "Parse empty object",
			input:    "{}",
			expected: ast.Object{},
			wantErr:  false,
		},
		{
			name:  "Parse simple object",
			input: `{"key": "value"}`,
			expected: ast.Object{
				"key": ast.String("value"),
			},
			wantErr: false,
		},
		{
			name: "Parse complex object",
			input: `{
				"string": "value",
				"number": 42,
				"bool": true,
				"null": null,
				"array": [1, 2, 3],
				"object": {"nested": "ok"}
			}`,
			expected: ast.Object{
				"string": ast.String("value"),
				"number": ast.Number(42),
				"bool":   ast.Boolean(true),
				"null":   ast.Null{},
				"array":  ast.Array{ast.Number(1), ast.Number(2), ast.Number(3)},
				"object": ast.Object{"nested": ast.String("ok")},
			},
			wantErr: false,
		},
		{
			name:    "Error: Missing colon",
			input:   `{"key" "value"}`,
			wantErr: true,
		},
		{
			name:    "Error: Missing value",
			input:   `{"key":}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := NewParser(l)
			p.nextToken() // Consume the first token
			got, err := p.parseObject()

			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Parser.parseObject() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParser_parseArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ast.Array
		wantErr  bool
	}{
		{
			name:     "Parse empty array",
			input:    "[]",
			expected: ast.Array{},
			wantErr:  false,
		},
		{
			name:     "Parse simple array",
			input:    "[1, 2, 3]",
			expected: ast.Array{ast.Number(1), ast.Number(2), ast.Number(3)},
			wantErr:  false,
		},
		{
			name:     "Parse mixed type array",
			input:    `[1, "two", true, null, {"key": "value"}, [1, 2]]`,
			expected: ast.Array{ast.Number(1), ast.String("two"), ast.Boolean(true), ast.Null{}, ast.Object{"key": ast.String("value")}, ast.Array{ast.Number(1), ast.Number(2)}},
			wantErr:  false,
		},
		{
			name:    "Error: Missing comma",
			input:   "[1 2]",
			wantErr: true,
		},
		{
			name:    "Error: Trailing comma",
			input:   "[1, 2,]",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := NewParser(l)
			p.nextToken() // Consume the first token
			got, err := p.parseArray()

			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Parser.parseArray() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParser_parseNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ast.Number
		wantErr  bool
	}{
		{
			name:     "Parse integer",
			input:    "42",
			expected: ast.Number(42),
			wantErr:  false,
		},
		{
			name:     "Parse float",
			input:    "3.14",
			expected: ast.Number(3.14),
			wantErr:  false,
		},
		{
			name:     "Parse negative number",
			input:    "-123.45",
			expected: ast.Number(-123.45),
			wantErr:  false,
		},
		{
			name:     "Parse exponential notation",
			input:    "1e10",
			expected: ast.Number(1e10),
			wantErr:  false,
		},
		{
			name:    "Error: Invalid number",
			input:   "12.34.56",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := NewParser(l)
			p.nextToken() // Consume the first token
			got, err := p.parseNumber()

			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.expected {
				t.Errorf("Parser.parseNumber() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// 在文件末尾添加以下基准测试函数

func BenchmarkParser_Parse_Simple(b *testing.B) {
	input := `{"key": "value"}`
	for i := 0; i < b.N; i++ {
		l := lexer.NewLexer(input)
		p := NewParser(l)
		_, err := p.Parse()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParser_Parse_Complex(b *testing.B) {
	input := `{
		"string": "value",
		"number": 42,
		"bool": true,
		"null": null,
		"array": [1, "two", false],
		"object": {"nested": "ok"}
	}`
	for i := 0; i < b.N; i++ {
		l := lexer.NewLexer(input)
		p := NewParser(l)
		_, err := p.Parse()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParser_Parse_LargeArray(b *testing.B) {
	input := "[" + strings.Repeat("1,", 9999) + "1]"
	for i := 0; i < b.N; i++ {
		l := lexer.NewLexer(input)
		p := NewParser(l)
		_, err := p.Parse()
		if err != nil {
			b.Fatal(err)
		}
	}
}
