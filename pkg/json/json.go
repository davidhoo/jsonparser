package json

import (
	"json-parser/internal/ast"
	"json-parser/internal/parser"
)

// Parse 解析 JSON 字符串并返回 ast.Node
func Parse(input string) (ast.Node, error) {
	p := parser.NewParser(input)
	return p.Parse()
}

// 如果需要，可以添加其他导出的函数或方法
