package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"json-parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: jsonparser <json_string> [jsonpath]")
		fmt.Println("   or: echo <json_string> | jsonparser - [jsonpath]")
		fmt.Println("\nExamples:")
		fmt.Println("  jsonparser '{\"name\": \"Alice\"}'")
		fmt.Println("  jsonparser '{\"user\": {\"name\": \"Bob\"}}' '$.user.name'")
		fmt.Println("  echo '{\"items\": [1,2,3]}' | jsonparser - '$.items[0]'")
		os.Exit(1)
	}

	input := os.Args[1]
	if input == "-" {
		// Read from stdin
		var buf strings.Builder
		_, err := io.Copy(&buf, os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		input = buf.String()
	}

	result, err := jsonparser.Parse(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	// If JSONPath is provided, extract the value
	if len(os.Args) >= 3 {
		jsonPath := os.Args[2]
		extracted, err := jsonparser.Get(result, jsonPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "JSONPath error: %v\n", err)
			os.Exit(1)
		}
		result = extracted
	}

	// Pretty print the result using standard JSON encoder
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling result: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}

