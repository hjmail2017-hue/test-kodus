package jsonparser

import (
	"reflect"
	"testing"
)

func TestParseJSONPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"root only", "$", false},
		{"simple dot notation", "$.key", false},
		{"nested dot notation", "$.foo.bar", false},
		{"array index", "$[0]", false},
		{"nested array index", "$.items[0]", false},
		{"bracket notation", `$["key"]`, false},
		{"bracket with spaces", `$[ "key" ]`, true}, // spaces not supported
		{"empty path", "", true},
		{"missing dollar", "key", true},
		{"invalid bracket", "$[", true},
		{"unclosed bracket", "$[0", true},
		{"invalid array index", "$[foo]", true},
		{"mixed notation", `$.items[0].name`, false},
		{"deep nesting", `$.a.b.c[0].d[1].e`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseJSONPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJSONPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJSONPathEvaluate(t *testing.T) {
	testData := map[string]interface{}{
		"name": "Alice",
		"age":  30,
		"address": map[string]interface{}{
			"city":    "Beijing",
			"street":  "Main St",
			"zipcode": "100000",
		},
		"scores": []interface{}{85, 92, 78},
		"items": []interface{}{
			map[string]interface{}{"id": 1, "name": "item1"},
			map[string]interface{}{"id": 2, "name": "item2"},
		},
	}

	tests := []struct {
		name     string
		path     string
		want     interface{}
		wantErr  bool
	}{
		{"root", "$", testData, false},
		{"simple key", "$.name", "Alice", false},
		{"nested key", "$.address.city", "Beijing", false},
		{"array index", "$.scores[0]", 85, false},
		{"negative array index", "$.scores[-1]", nil, true},
		{"out of bounds", "$.scores[10]", nil, true},
		{"nested array", "$.items[0].name", "item1", false},
		{"deep nesting", "$.items[1].id", 2, false},
		{"non-existent key", "$.nonexistent", nil, true},
		{"access array as object", "$.scores.name", nil, true},
		{"access object as array", "$.name[0]", nil, true},
		{"bracket notation", `$["address"]["city"]`, "Beijing", false},
		{"mixed notation", `$.items[0]["name"]`, "item1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jp, err := ParseJSONPath(tt.path)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ParseJSONPath() unexpected error = %v", err)
				}
				return
			}

			got, err := jp.Evaluate(testData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFunction(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Bob",
			"prefs": map[string]interface{}{
				"theme": "dark",
				"lang":  "en",
			},
		},
	}

	tests := []struct {
		name    string
		path    string
		want    interface{}
		wantErr bool
	}{
		{"valid path", "$.user.name", "Bob", false},
		{"nested path", "$.user.prefs.theme", "dark", false},
		{"invalid path", "$.user.age", nil, true},
		{"malformed path", "user.name", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(data, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}