package jsonparser

import (
	"fmt"
	"strconv"
	"strings"
)

// JSONPath represents a parsed JSONPath expression.
type JSONPath struct {
	segments []segment
}

type segment interface {
	apply(data interface{}) (interface{}, error)
}

type rootSegment struct{}

type childSegment struct {
	key string
}

type arrayIndexSegment struct {
	index int
}

// ParseJSONPath parses a JSONPath expression.
func ParseJSONPath(path string) (*JSONPath, error) {
	if path == "" {
		return nil, fmt.Errorf("empty JSONPath expression")
	}

	// Must start with $
	if !strings.HasPrefix(path, "$") {
		return nil, fmt.Errorf("JSONPath must start with $")
	}

	segments := []segment{&rootSegment{}}
	remaining := path[1:]

	for remaining != "" {
		switch {
		case strings.HasPrefix(remaining, "."):
			// Child segment: .key or .key1.key2
			remaining = remaining[1:]
			if remaining == "" {
				return nil, fmt.Errorf("expected key after .")
			}

			// Check if it's a bracket notation: ["key"]
			if strings.HasPrefix(remaining, "[") {
				// Parse bracket notation
				end := strings.Index(remaining, "]")
				if end == -1 {
					return nil, fmt.Errorf("unclosed bracket")
				}
				content := remaining[1:end]

				if strings.HasPrefix(content, `"`) && strings.HasSuffix(content, `"`) {
					// Quoted key: ["key"]
					key := content[1 : len(content)-1]
					segments = append(segments, &childSegment{key: key})
				} else {
					// Array index: [0] or [1]
					idx, err := strconv.Atoi(content)
					if err != nil {
						return nil, fmt.Errorf("invalid array index: %s", content)
					}
					segments = append(segments, &arrayIndexSegment{index: idx})
				}
				remaining = remaining[end+1:]
			} else {
				// Dot notation: .key
				// Extract key (until next dot or bracket or end)
				end := len(remaining)
				for i, ch := range remaining {
					if ch == '.' || ch == '[' {
						end = i
						break
					}
				}
				key := remaining[:end]
				if key == "" {
					return nil, fmt.Errorf("empty key in dot notation")
				}
				segments = append(segments, &childSegment{key: key})
				remaining = remaining[end:]
			}

		case strings.HasPrefix(remaining, "["):
			// Array index or quoted key
			end := strings.Index(remaining, "]")
			if end == -1 {
				return nil, fmt.Errorf("unclosed bracket")
			}
			content := remaining[1:end]

			if strings.HasPrefix(content, `"`) && strings.HasSuffix(content, `"`) {
				// Quoted key: ["key"]
				key := content[1 : len(content)-1]
				segments = append(segments, &childSegment{key: key})
			} else {
				// Array index: [0]
				idx, err := strconv.Atoi(content)
				if err != nil {
					return nil, fmt.Errorf("invalid array index: %s", content)
				}
				segments = append(segments, &arrayIndexSegment{index: idx})
			}
			remaining = remaining[end+1:]

		default:
			return nil, fmt.Errorf("invalid JSONPath syntax at: %s", remaining)
		}
	}

	return &JSONPath{segments: segments}, nil
}

// Evaluate evaluates the JSONPath against the given data.
func (jp *JSONPath) Evaluate(data interface{}) (interface{}, error) {
	var current interface{} = data

	for i, seg := range jp.segments {
		if i == 0 {
			// First segment must be root, just continue
			if _, ok := seg.(*rootSegment); !ok {
				return nil, fmt.Errorf("first segment must be root")
			}
			continue
		}

		result, err := seg.apply(current)
		if err != nil {
			return nil, fmt.Errorf("segment %d: %w", i, err)
		}
		current = result
	}

	return current, nil
}

// apply for rootSegment does nothing
func (r *rootSegment) apply(data interface{}) (interface{}, error) {
	return data, nil
}

// apply for childSegment extracts a child from an object
func (c *childSegment) apply(data interface{}) (interface{}, error) {
	obj, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected object, got %T", data)
	}

	value, exists := obj[c.key]
	if !exists {
		return nil, fmt.Errorf("key '%s' not found", c.key)
	}

	return value, nil
}

// apply for arrayIndexSegment extracts an element from an array
func (a *arrayIndexSegment) apply(data interface{}) (interface{}, error) {
	arr, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", data)
	}

	if a.index < 0 || a.index >= len(arr) {
		return nil, fmt.Errorf("array index %d out of bounds (length %d)", a.index, len(arr))
	}

	return arr[a.index], nil
}

// Get uses JSONPath to extract a value from JSON data.
func Get(data interface{}, path string) (interface{}, error) {
	jp, err := ParseJSONPath(path)
	if err != nil {
		return nil, err
	}
	return jp.Evaluate(data)
}