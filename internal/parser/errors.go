package parser

import (
	"fmt"
	"json-parser/internal/lexer"
)

type ParseError struct {
	Line   int
	Column int
	Msg    string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at line %d, column %d: %s", e.Line, e.Column, e.Msg)
}

func newParseError(l *lexer.Lexer, msg string) *ParseError {
	return &ParseError{
		Line:   l.GetLine(),
		Column: l.GetColumn(),
		Msg:    msg,
	}
}
