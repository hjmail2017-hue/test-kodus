package jsonparser

import (
	"reflect"
	"testing"
)

func TestParseLiteral(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    interface{}
		wantErr bool
	}{
		{"null", "null", nil, false},
		{"true", "true", true, false},
		{"false", "false", false, false},
		{"integer", "42", int64(42), false},
		{"negative integer", "-42", int64(-42), false},
		{"float", "3.14", 3.14, false},
		{"scientific notation", "1.2e3", 1.2e3, false},
		{"string", `"hello"`, "hello", false},
		{"empty string", `""`, "", false},
		{"string with escape", `"hello\nworld"`, "hello\nworld", false},
		{"string with quote", `"he said \"hello\""`, "he said \"hello\"", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v (%T), want %v (%T)", got, got, tt.want, tt.want)
			}
		})
	}
}

func TestParseArray(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    interface{}
		wantErr bool
	}{
		{"empty array", "[]", []interface{}{}, false},
		{"array of numbers", "[1, 2, 3]", []interface{}{int64(1), int64(2), int64(3)}, false},
		{"mixed array", `[1, "two", true]`, []interface{}{int64(1), "two", true}, false},
		{"nested array", `[[1, 2], [3, 4]]`, []interface{}{
			[]interface{}{int64(1), int64(2)},
			[]interface{}{int64(3), int64(4)},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseObject(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    interface{}
		wantErr bool
	}{
		{"empty object", "{}", map[string]interface{}{}, false},
		{"simple object", `{"key": "value"}`, map[string]interface{}{"key": "value"}, false},
		{"multiple keys", `{"a": 1, "b": true}`, map[string]interface{}{"a": int64(1), "b": true}, false},
		{"nested object", `{"obj": {"inner": "val"}}`, map[string]interface{}{
			"obj": map[string]interface{}{"inner": "val"},
		}, false},
		{"array in object", `{"arr": [1,2]}`, map[string]interface{}{
			"arr": []interface{}{int64(1), int64(2)},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unexpected character", "foo"},
		{"unterminated string", `"hello`},
		{"invalid escape", `"\x"`},
		{"missing colon", `{"key" "value"}`},
		{"missing comma in array", "[1 2]"},
		{"trailing comma in array", "[1,]"},
		{"trailing comma in object", `{"a":1,}`},
		{"invalid number", "12."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err == nil {
				t.Errorf("Parse() expected error but got %v", got)
			} else {
				t.Logf("expected error: %v", err)
			}
		})
	}
}

func TestLexer(t *testing.T) {
	input := `{"key": "value", "num": 42}`
	lexer := NewLexer(input)
	tokens := []Token{}
	for {
		tok := lexer.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF || tok.Type == TokenError {
			break
		}
	}
	// Just ensure no error tokens
	for _, tok := range tokens {
		if tok.Type == TokenError {
			t.Errorf("unexpected error token: %v", tok)
		}
	}
}
