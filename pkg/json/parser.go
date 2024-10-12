package json

import (
	"json-parser/internal/ast"
	"json-parser/internal/lexer"
	"json-parser/internal/parser"
)

// ParseJSON 解析 JSON 字符串并返回 ast.Node
func ParseJSON(input string) (ast.Node, error) {
	lex := lexer.NewLexer(input)
	p := parser.NewParser(lex)
	return p.Parse()
}
