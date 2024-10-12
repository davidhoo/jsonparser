package lexer

// TokenType represents the type of a token.
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

// Token represents a lexical token.
type Token struct {
	Type  TokenType
	Value string
}

// Lexer represents a lexical analyzer.
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

// NewLexer creates a new Lexer instance.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position in the input string.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// NextToken returns the next token in the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	var tok Token
	switch l.ch {
	case 0:
		tok = Token{Type: TokenEOF, Value: ""}
	case '{':
		tok = Token{Type: TokenLeftBrace, Value: string(l.ch)}
	case '}':
		tok = Token{Type: TokenRightBrace, Value: string(l.ch)}
	case '[':
		tok = Token{Type: TokenLeftBracket, Value: string(l.ch)}
	case ']':
		tok = Token{Type: TokenRightBracket, Value: string(l.ch)}
	case ',':
		tok = Token{Type: TokenComma, Value: string(l.ch)}
	case ':':
		tok = Token{Type: TokenColon, Value: string(l.ch)}
	case '"':
		return l.readString()
	default:
		if isDigit(l.ch) || l.ch == '-' {
			return l.readNumber()
		} else if isLetter(l.ch) {
			return l.readKeyword()
		} else {
			tok = Token{Type: TokenString, Value: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

// skipWhitespace skips any whitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readString reads a string token.
func (l *Lexer) readString() Token {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	str := l.input[position:l.position]
	l.readChar() // consume closing quote
	return Token{Type: TokenString, Value: str}
}

// readNumber reads a number token.
func (l *Lexer) readNumber() Token {
	position := l.position
	for isDigit(l.ch) || l.ch == '.' || l.ch == 'e' || l.ch == 'E' || l.ch == '+' || l.ch == '-' {
		l.readChar()
	}
	return Token{Type: TokenNumber, Value: l.input[position:l.position]}
}

// readKeyword reads a keyword token (true, false, or null).
func (l *Lexer) readKeyword() Token {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	keyword := l.input[position:l.position]
	switch keyword {
	case "true":
		return Token{Type: TokenTrue, Value: keyword}
	case "false":
		return Token{Type: TokenFalse, Value: keyword}
	case "null":
		return Token{Type: TokenNull, Value: keyword}
	default:
		return Token{Type: TokenString, Value: keyword}
	}
}

// isDigit checks if a character is a digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isLetter checks if a character is a letter.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// PeekChar returns the next character without advancing the position.
func (l *Lexer) PeekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// GetLine returns the current line number.
func (l *Lexer) GetLine() int {
	return l.line
}

// GetColumn returns the current column number.
func (l *Lexer) GetColumn() int {
	return l.column
}
