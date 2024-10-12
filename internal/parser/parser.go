package parser

import (
	"fmt"
	"json-parser/internal/ast"
	"json-parser/internal/lexer"
	"strconv"
)

// Parser represents a JSON parser.
type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
}

// NewParser creates a new Parser instance.
func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances the parser to the next token.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Parse parses the input and returns the root AST node.
func (p *Parser) Parse() (ast.Node, error) {
	var result ast.Node
	var err error

	switch p.curToken.Type {
	case lexer.TokenLeftBrace:
		result, err = p.parseObject()
	case lexer.TokenLeftBracket:
		result, err = p.parseArray()
	case lexer.TokenString:
		result = ast.String(p.curToken.Value)
	case lexer.TokenNumber:
		result, err = p.parseNumber()
	case lexer.TokenTrue:
		result = ast.Boolean(true)
	case lexer.TokenFalse:
		result = ast.Boolean(false)
	case lexer.TokenNull:
		result = ast.Null{}
	case lexer.TokenEOF:
		return nil, newParseError(p.l, fmt.Sprintf("unexpected end of input"))
	default:
		return nil, newParseError(p.l, fmt.Sprintf("unexpected token: %v", p.curToken))
	}

	if err != nil {
		return nil, err
	}

	// 检查是否还有多余的标记
	p.nextToken()
	if p.curToken.Type != lexer.TokenEOF {
		return nil, newParseError(p.l, fmt.Sprintf("unexpected token at end of input: %v", p.curToken))
	}

	return result, nil
}

// parseObject parses a JSON object.
func (p *Parser) parseObject() (ast.Object, error) {
	obj := make(ast.Object)
	p.nextToken() // consume '{'

	for p.curToken.Type != lexer.TokenRightBrace {
		if p.curToken.Type != lexer.TokenString {
			return nil, newParseError(p.l, fmt.Sprintf("expected string key, got %v", p.curToken))
		}
		key := p.curToken.Value
		p.nextToken()

		if p.curToken.Type != lexer.TokenColon {
			return nil, newParseError(p.l, fmt.Sprintf("expected ':', got %v", p.curToken))
		}
		p.nextToken()

		value, err := p.Parse()
		if err != nil {
			return nil, err
		}
		obj[key] = value

		p.nextToken()
		if p.curToken.Type == lexer.TokenComma {
			p.nextToken()
		} else if p.curToken.Type != lexer.TokenRightBrace {
			return nil, newParseError(p.l, fmt.Sprintf("expected ',' or '}', got %v", p.curToken))
		}
	}

	return obj, nil
}

// parseArray parses a JSON array.
func (p *Parser) parseArray() (ast.Array, error) {
	arr := make(ast.Array, 0)
	p.nextToken() // consume '['

	for p.curToken.Type != lexer.TokenRightBracket {
		value, err := p.Parse()
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)

		p.nextToken()
		if p.curToken.Type == lexer.TokenComma {
			p.nextToken()
		} else if p.curToken.Type != lexer.TokenRightBracket {
			return nil, newParseError(p.l, fmt.Sprintf("expected ',' or ']', got %v", p.curToken))
		}
	}

	return arr, nil
}

// parseNumber parses a JSON number.
func (p *Parser) parseNumber() (ast.Number, error) {
	n, err := strconv.ParseFloat(p.curToken.Value, 64)
	if err != nil {
		return ast.Number(0), newParseError(p.l, fmt.Sprintf("invalid number: %s", p.curToken.Value))
	}
	return ast.Number(n), nil
}

// Add more parsing methods as needed...
