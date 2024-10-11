package lexer

import (
	"unicode"
	"unicode/utf8"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenString
	TokenNumber
	TokenTrue
	TokenFalse
	TokenNull
	TokenLeftBrace
	TokenRightBrace
	TokenLeftBracket
	TokenRightBracket
	TokenComma
	TokenColon
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
	line         int
	column       int
}

// New 创建一个新的词法分析器实例
func New(input string) *Lexer {
	return &Lexer{
		input: input,
		// ... 其他初始化
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		var width int
		l.ch, width = utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.position = l.readPosition
		l.readPosition += width
	}

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// Add more methods for tokenization...

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.position >= len(l.input) {
		return Token{Type: TokenEOF}
	}

	switch l.input[l.position] {
	case '{':
		l.position++
		return Token{Type: TokenLeftBrace, Value: "{"}
	case '}':
		l.position++
		return Token{Type: TokenRightBrace, Value: "}"}
	case '[':
		l.position++
		return Token{Type: TokenLeftBracket, Value: "["}
	case ']':
		l.position++
		return Token{Type: TokenRightBracket, Value: "]"}
	case ',':
		l.position++
		return Token{Type: TokenComma, Value: ","}
	case ':':
		l.position++
		return Token{Type: TokenColon, Value: ":"}
	case '"':
		return l.readString()
	}

	if unicode.IsDigit(rune(l.input[l.position])) || l.input[l.position] == '-' {
		return l.readNumber()
	}

	return l.readKeyword()
}

func (l *Lexer) skipWhitespace() {
	for l.position < len(l.input) && unicode.IsSpace(rune(l.input[l.position])) {
		l.position++
	}
}

func (l *Lexer) readString() Token {
	l.position++ // Skip the opening quote
	start := l.position
	for l.position < len(l.input) && l.input[l.position] != '"' {
		if l.input[l.position] == '\\' && l.position+1 < len(l.input) {
			l.position++ // Skip the escape character
		}
		l.position++
	}
	if l.position >= len(l.input) {
		return Token{Type: TokenString, Value: l.input[start:]}
	}
	value := l.input[start:l.position]
	l.position++ // Skip the closing quote
	return Token{Type: TokenString, Value: value}
}

func (l *Lexer) readNumber() Token {
	start := l.position
	for l.position < len(l.input) && (unicode.IsDigit(rune(l.input[l.position])) || l.input[l.position] == '.' || l.input[l.position] == 'e' || l.input[l.position] == 'E' || l.input[l.position] == '+' || l.input[l.position] == '-') {
		l.position++
	}
	return Token{Type: TokenNumber, Value: l.input[start:l.position]}
}

func (l *Lexer) readKeyword() Token {
	start := l.position
	for l.position < len(l.input) && unicode.IsLetter(rune(l.input[l.position])) {
		l.position++
	}
	value := l.input[start:l.position]
	switch value {
	case "true":
		return Token{Type: TokenTrue, Value: value}
	case "false":
		return Token{Type: TokenFalse, Value: value}
	case "null":
		return Token{Type: TokenNull, Value: value}
	default:
		// This should not happen in valid JSON, but we'll return it as a string for now
		return Token{Type: TokenString, Value: value}
	}
}

// Implement readString, readNumber, and readKeyword methods here
