package parser

import (
	"fmt"
	"json-parser/internal/ast"
	"json-parser/internal/lexer"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string
}

// NewParser 创建一个新的解析器实例
func NewParser(input string) *Parser {
	return &Parser{l: lexer.New(input)} // 使用 New 而不是 NewLexer
}

func (p *Parser) Parse() (ast.Node, error) {
	token := p.l.NextToken()
	switch token.Type {
	case lexer.TokenLeftBrace:
		return p.parseObject()
	case lexer.TokenLeftBracket:
		return p.parseArray()
	case lexer.TokenString:
		return ast.String(token.Value), nil
	case lexer.TokenNumber:
		return ast.Number(parseFloat(token.Value)), nil
	case lexer.TokenTrue:
		return ast.Boolean(true), nil
	case lexer.TokenFalse:
		return ast.Boolean(false), nil
	case lexer.TokenNull:
		return ast.Null{}, nil
	default:
		return nil, fmt.Errorf("unexpected token: %v", token)
	}
}

func (p *Parser) parseObject() (ast.Object, error) {
	obj := make(ast.Object)
	for {
		keyToken := p.l.NextToken()
		if keyToken.Type == lexer.TokenRightBrace {
			break
		}
		if keyToken.Type != lexer.TokenString {
			return nil, fmt.Errorf("expected string key, got %v", keyToken)
		}

		colonToken := p.l.NextToken()
		if colonToken.Type != lexer.TokenColon {
			return nil, fmt.Errorf("expected colon, got %v", colonToken)
		}

		value, err := p.Parse()
		if err != nil {
			return nil, err
		}

		obj[keyToken.Value] = value

		commaToken := p.l.NextToken()
		if commaToken.Type == lexer.TokenRightBrace {
			break
		}
		if commaToken.Type != lexer.TokenComma {
			return nil, fmt.Errorf("expected comma or closing brace, got %v", commaToken)
		}
	}
	return obj, nil
}

func (p *Parser) parseArray() (ast.Array, error) {
	arr := make(ast.Array, 0)
	for {
		value, err := p.Parse()
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)

		commaToken := p.l.NextToken()
		if commaToken.Type == lexer.TokenRightBracket {
			break
		}
		if commaToken.Type != lexer.TokenComma {
			return nil, fmt.Errorf("expected comma or closing bracket, got %v", commaToken)
		}
	}
	return arr, nil
}

func parseFloat(s string) float64 {
	// This is a simplified version. You might want to use strconv.ParseFloat in a real implementation.
	var result float64
	fmt.Sscanf(s, "%f", &result)
	return result
}

// Add more parsing methods as needed...
