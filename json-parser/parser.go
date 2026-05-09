package jsonparser

import (
	"fmt"
	"strconv"
)

// Parser parses JSON tokens into Go values.
type Parser struct {
	lexer *Lexer
	cur   Token
	peek  Token
}

// NewParser creates a new parser for the given input.
func NewParser(input string) *Parser {
	lexer := NewLexer(input)
	p := &Parser{lexer: lexer}
	p.nextToken()
	p.nextToken() // initialize both cur and peek
	return p
}

// nextToken advances the tokens.
func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.lexer.NextToken()
}

// Parse parses the entire JSON input.
func (p *Parser) Parse() (interface{}, error) {
	value, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	// Ensure there's no extra tokens after the JSON value
	if p.cur.Type != TokenEOF {
		return nil, fmt.Errorf("unexpected token %s at position %d", p.cur.String(), p.cur.Pos)
	}
	return value, nil
}

// parseValue parses a JSON value.
func (p *Parser) parseValue() (interface{}, error) {
	switch p.cur.Type {
	case TokenLeftBrace:
		return p.parseObject()
	case TokenLeftBracket:
		return p.parseArray()
	case TokenString:
		return p.parseString()
	case TokenNumber:
		return p.parseNumber()
	case TokenTrue:
		p.nextToken()
		return true, nil
	case TokenFalse:
		p.nextToken()
		return false, nil
	case TokenNull:
		p.nextToken()
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected token %s at position %d", p.cur.String(), p.cur.Pos)
	}
}

// parseObject parses a JSON object.
func (p *Parser) parseObject() (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	p.nextToken() // consume '{'
	// Handle empty object
	if p.cur.Type == TokenRightBrace {
		p.nextToken()
		return obj, nil
	}
	for {
		// Parse key (must be a string)
		if p.cur.Type != TokenString {
			return nil, fmt.Errorf("expected string key, got %s at position %d", p.cur.String(), p.cur.Pos)
		}
		key := p.cur.Literal
		p.nextToken()
		// Expect colon
		if p.cur.Type != TokenColon {
			return nil, fmt.Errorf("expected colon after key, got %s at position %d", p.cur.String(), p.cur.Pos)
		}
		p.nextToken()
		// Parse value
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		obj[key] = value
		// Check for comma or closing brace
		if p.cur.Type == TokenRightBrace {
			p.nextToken()
			break
		}
		if p.cur.Type != TokenComma {
			return nil, fmt.Errorf("expected comma or closing brace, got %s at position %d", p.cur.String(), p.cur.Pos)
		}
		p.nextToken()
	}
	return obj, nil
}

// parseArray parses a JSON array.
func (p *Parser) parseArray() ([]interface{}, error) {
	arr := make([]interface{}, 0)
	p.nextToken() // consume '['
	// Handle empty array
	if p.cur.Type == TokenRightBracket {
		p.nextToken()
		return arr, nil
	}
	for {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)
		// Check for comma or closing bracket
		if p.cur.Type == TokenRightBracket {
			p.nextToken()
			break
		}
		if p.cur.Type != TokenComma {
			return nil, fmt.Errorf("expected comma or closing bracket, got %s at position %d", p.cur.String(), p.cur.Pos)
		}
		p.nextToken()
	}
	return arr, nil
}

// parseString returns the string value (already captured as literal).
func (p *Parser) parseString() (string, error) {
	val := p.cur.Literal
	p.nextToken()
	return val, nil
}

// parseNumber parses a number token into float64.
func (p *Parser) parseNumber() (interface{}, error) {
	val := p.cur.Literal
	// Validate JSON number format
	if !isValidJSONNumber(val) {
		return nil, fmt.Errorf("invalid number %s at position %d", val, p.cur.Pos)
	}
	p.nextToken()
	// Try parsing as integer first
	if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		return i, nil
	}
	// Try parsing as float
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f, nil
	}
	return nil, fmt.Errorf("invalid number %s at position %d", val, p.cur.Pos)
}

// isValidJSONNumber checks if the string is a valid JSON number.
func isValidJSONNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	i := 0
	// Optional minus sign
	if s[i] == '-' {
		i++
		if i >= len(s) {
			return false
		}
	}
	// Integer part
	if s[i] == '0' {
		i++
	} else if '1' <= s[i] && s[i] <= '9' {
		i++
		for i < len(s) && '0' <= s[i] && s[i] <= '9' {
			i++
		}
	} else {
		return false
	}
	// Fractional part (optional)
	if i < len(s) && s[i] == '.' {
		i++
		if i >= len(s) || s[i] < '0' || s[i] > '9' {
			return false // must have at least one digit after decimal point
		}
		for i < len(s) && '0' <= s[i] && s[i] <= '9' {
			i++
		}
	}
	// Exponent part (optional)
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		if i < len(s) && (s[i] == '+' || s[i] == '-') {
			i++
		}
		if i >= len(s) || s[i] < '0' || s[i] > '9' {
			return false // must have at least one digit after exponent
		}
		for i < len(s) && '0' <= s[i] && s[i] <= '9' {
			i++
		}
	}
	// Must consume entire string
	return i == len(s)
}

// Parse parses a JSON string and returns the corresponding Go value.
func Parse(input string) (interface{}, error) {
	parser := NewParser(input)
	return parser.Parse()
}
