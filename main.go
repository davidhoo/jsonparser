// The json-parser program is a command-line tool for parsing, validating,
// and manipulating JSON data. It demonstrates the usage of the json package
// which provides JSON parsing and manipulation functionality.
package main

import (
	"fmt"
	"json-parser/pkg/json"
	"os"
)

func main() {
	// Example JSON string to demonstrate the functionality
	jsonStr := `{
		"name": "John Doe",
		"age": 30,
		"city": "New York",
		"hobbies": ["reading", "swimming", "coding"]
	}`

	fmt.Println("Original JSON:")
	fmt.Println(jsonStr)

	// Demonstrate pretty printing
	fmt.Println("\nPretty printed JSON:")
	prettyJSON, err := json.PrettyPrint(jsonStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error pretty printing JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(prettyJSON)

	// Demonstrate getting values by path
	fmt.Println("\nGetting values by path:")
	paths := [][]string{
		{"name"},
		{"age"},
		{"hobbies", "1"},
		{"address"},
		{"hobbies", "3"},
	}

	for _, path := range paths {
		value, err := json.GetValueByPath(jsonStr, path)
		if err != nil {
			fmt.Printf("Error getting value for path %v: %v\n", path, err)
		} else {
			fmt.Printf("Value at path %v: %v\n", path, value)
		}
	}
}
