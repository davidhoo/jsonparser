package lexer

import (
	"testing"
)

func TestLexer(t *testing.T) {
	input := `{"name": "Alice", "age": 30}`
	lexer := NewLexer(input)

	expectedTokens := []struct {
		Type  TokenType
		Value string
	}{
		{TokenLeftBrace, "{"},
		{TokenString, "name"},
		{TokenColon, ":"},
		{TokenString, "Alice"},
		{TokenComma, ","},
		{TokenString, "age"},
		{TokenColon, ":"},
		{TokenNumber, "30"},
		{TokenRightBrace, "}"},
		{TokenEOF, ""},
	}

	for _, expected := range expectedTokens {
		token := lexer.NextToken()
		if token.Type != expected.Type {
			t.Errorf("Expected token type %v, got %v", expected.Type, token.Type)
		}
		if token.Value != expected.Value {
			t.Errorf("Expected token value %q, got %q", expected.Value, token.Value)
		}
	}
}

func TestLexer_NextToken(t *testing.T) {
	input := `{"key": 123, "array": [true, false, null]}`
	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{TokenLeftBrace, "{"},
		{TokenString, "key"},
		{TokenColon, ":"},
		{TokenNumber, "123"},
		{TokenComma, ","},
		{TokenString, "array"},
		{TokenColon, ":"},
		{TokenLeftBracket, "["},
		{TokenTrue, "true"},
		{TokenComma, ","},
		{TokenFalse, "false"},
		{TokenComma, ","},
		{TokenNull, "null"},
		{TokenRightBracket, "]"},
		{TokenRightBrace, "}"},
		{TokenEOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedValue, tok.Value)
		}
	}
}
