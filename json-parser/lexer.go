package jsonparser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// TokenType represents the type of a token.
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError
	TokenLeftBrace    // {
	TokenRightBrace   // }
	TokenLeftBracket  // [
	TokenRightBracket // ]
	TokenColon        // :
	TokenComma        // ,
	TokenString
	TokenNumber
	TokenTrue
	TokenFalse
	TokenNull
)

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string // raw string of the token
	Pos     int    // start position in input
}

// String returns a human-readable representation of the token.
func (t Token) String() string {
	switch t.Type {
	case TokenEOF:
		return "EOF"
	case TokenError:
		return fmt.Sprintf("ERROR(%s)", t.Literal)
	case TokenLeftBrace:
		return "{"
	case TokenRightBrace:
		return "}"
	case TokenLeftBracket:
		return "["
	case TokenRightBracket:
		return "]"
	case TokenColon:
		return ":"
	case TokenComma:
		return ","
	case TokenString:
		return fmt.Sprintf("STRING(%s)", t.Literal)
	case TokenNumber:
		return fmt.Sprintf("NUMBER(%s)", t.Literal)
	case TokenTrue:
		return "true"
	case TokenFalse:
		return "false"
	case TokenNull:
		return "null"
	default:
		return "UNKNOWN"
	}
}

// Lexer scans the input string and produces tokens.
type Lexer struct {
	input   string // the string being scanned
	pos     int    // current position in the input (points to current rune)
	readPos int    // current reading position (after current rune)
	ch      rune   // current rune under examination
	start   int    // start position of current token
}

// NewLexer creates a new lexer for the given input.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads the next rune from the input.
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL" character, signifies EOF
	} else {
		l.ch, _ = utf8.DecodeRuneInString(l.input[l.readPos:])
	}
	l.pos = l.readPos
	l.readPos += utf8.RuneLen(l.ch)
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	l.start = l.pos

	switch l.ch {
	case '{':
		l.readChar()
		return Token{Type: TokenLeftBrace, Literal: "{", Pos: l.start}
	case '}':
		l.readChar()
		return Token{Type: TokenRightBrace, Literal: "}", Pos: l.start}
	case '[':
		l.readChar()
		return Token{Type: TokenLeftBracket, Literal: "[", Pos: l.start}
	case ']':
		l.readChar()
		return Token{Type: TokenRightBracket, Literal: "]", Pos: l.start}
	case ':':
		l.readChar()
		return Token{Type: TokenColon, Literal: ":", Pos: l.start}
	case ',':
		l.readChar()
		return Token{Type: TokenComma, Literal: ",", Pos: l.start}
	case '"':
		return l.readString()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		return l.readNumber()
	case 0:
		return Token{Type: TokenEOF, Literal: "", Pos: l.start}
	default:
		if isLetter(l.ch) {
			ident := l.readIdentifier()
			switch ident {
			case "true":
				return Token{Type: TokenTrue, Literal: ident, Pos: l.start}
			case "false":
				return Token{Type: TokenFalse, Literal: ident, Pos: l.start}
			case "null":
				return Token{Type: TokenNull, Literal: ident, Pos: l.start}
			default:
				return Token{Type: TokenError, Literal: fmt.Sprintf("unexpected identifier %s", ident), Pos: l.start}
			}
		}
		return Token{Type: TokenError, Literal: fmt.Sprintf("unexpected character %c", l.ch), Pos: l.start}
	}
}

// skipWhitespace skips over whitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier reads an identifier (for true, false, null).
func (l *Lexer) readIdentifier() string {
	start := l.pos
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.pos]
}

// isLetter checks if a rune is a letter (a-z, A-Z).
func isLetter(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

// readNumber reads a number token.
func (l *Lexer) readNumber() Token {
	start := l.pos
	if l.ch == '-' {
		l.readChar()
	}
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	if l.ch == 'e' || l.ch == 'E' {
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	literal := l.input[start:l.pos]
	return Token{Type: TokenNumber, Literal: literal, Pos: start}
}

// isDigit checks if a rune is a digit.
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// isHexDigit checks if a rune is a hexadecimal digit.
func isHexDigit(ch rune) bool {
	return ('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

// readString reads a string token.
func (l *Lexer) readString() Token {
	l.readChar() // consume the opening quote
	var buf strings.Builder
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar() // skip backslash
			if l.ch == 0 {
				return Token{Type: TokenError, Literal: "unterminated escape sequence", Pos: l.pos}
			}
			switch l.ch {
			case '"':
				buf.WriteRune('"')
			case '\\':
				buf.WriteRune('\\')
			case '/':
				buf.WriteRune('/')
			case 'b':
				buf.WriteRune('\b')
			case 'f':
				buf.WriteRune('\f')
			case 'n':
				buf.WriteRune('\n')
			case 'r':
				buf.WriteRune('\r')
			case 't':
				buf.WriteRune('\t')
			case 'u':
				// Unicode escape - simplified: just skip 4 hex digits
				for i := 0; i < 4; i++ {
					l.readChar()
					if !isHexDigit(l.ch) {
						return Token{Type: TokenError, Literal: "invalid Unicode escape", Pos: l.pos}
					}
				}
				// For simplicity, we'll just write the placeholder
				buf.WriteString("\\uXXXX")
			default:
				return Token{Type: TokenError, Literal: fmt.Sprintf("invalid escape sequence \\%c", l.ch), Pos: l.pos}
			}
		} else {
			buf.WriteRune(l.ch)
		}
		l.readChar()
	}
	if l.ch == 0 {
		return Token{Type: TokenError, Literal: "unterminated string", Pos: l.pos}
	}
	l.readChar() // consume the closing quote
	return Token{Type: TokenString, Literal: buf.String(), Pos: l.start}
}
