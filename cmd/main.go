package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func main() {
	// 添加新的命令行标志
	queryFlag := flag.String("q", "", "XPath-like query string to filter JSON")
	fileFlag := flag.String("f", "", "JSON file path")
	helpFlag := flag.Bool("h", false, "Show help message")

	// 自定义 usage 信息
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\n", color.CyanString("JSON Parser with XPath-like Query Support"))
		fmt.Fprintf(os.Stderr, "%s\n", color.YellowString("Usage:"))
		fmt.Fprintf(os.Stderr, "  %s %s\n", os.Args[0], color.GreenString("[-f <json_file>] [-q <query>]"))
		fmt.Fprintf(os.Stderr, "  %s %s\n\n", os.Args[0], color.GreenString("<json_file> [-q <query>]"))

		fmt.Fprintf(os.Stderr, "%s\n", color.YellowString("Options:"))
		fmt.Fprintf(os.Stderr, "  %s\t%s\n", color.GreenString("-f <json_file>"), "Specify the JSON file path")
		fmt.Fprintf(os.Stderr, "  %s\t%s\n", color.GreenString("-q <query>"), "XPath-like query string to filter JSON")
		fmt.Fprintf(os.Stderr, "  %s\t\t%s\n\n", color.GreenString("-h"), "Show this help message")

		fmt.Fprintf(os.Stderr, "%s\n", color.YellowString("Query Examples:"))
		fmt.Fprintf(os.Stderr, "  %s : Get the first user\n", color.GreenString("-q \"/data/users[0]\""))
		fmt.Fprintf(os.Stderr, "  %s : Find user with name 'Alice'\n", color.GreenString("-q \"/data/users[@name='Alice']\""))
		fmt.Fprintf(os.Stderr, "  %s : Find products with price over 1000\n", color.GreenString("-q \"/data/products[price>1000]\""))
		fmt.Fprintf(os.Stderr, "  %s : Get all notification settings\n", color.GreenString("-q \"/settings/notifications/*\""))
	}

	flag.Parse()

	// 如果没有提供任何参数或使用了 -h 标志，显示使用说明并退出
	if len(os.Args) == 1 || *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	var filePath string
	if *fileFlag != "" {
		filePath = *fileFlag
	} else if flag.NArg() > 0 {
		filePath = flag.Arg(0)
	} else {
		fmt.Println("Error: Please provide a JSON file path using -f flag or as an argument")
		flag.Usage()
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	if *queryFlag != "" {
		// 如果提供了查询，执行查询
		result, err := queryJSON(jsonData, *queryFlag)
		if err != nil {
			fmt.Printf("Error executing query: %v\n", err)
			os.Exit(1)
		}
		printJSON(result)
	} else {
		// 否则打印整个 JSON
		printJSON(jsonData)
	}
}

func queryJSON(data interface{}, query string) (interface{}, error) {
	parts := splitQuery(query)
	current := data

	for _, part := range parts {
		var err error
		current, err = processPart(current, part)
		if err != nil {
			return nil, err
		}
	}

	return current, nil
}

func splitQuery(query string) []string {
	var parts []string
	var current strings.Builder
	inBracket := false

	for _, char := range query {
		if char == '[' {
			inBracket = true
		}
		if char == ']' {
			inBracket = false
		}
		if char == '/' && !inBracket {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(char)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}

func processPart(data interface{}, part string) (interface{}, error) {
	if strings.HasPrefix(part, "@") {
		return processAttribute(data, part[1:])
	}
	if strings.Contains(part, "[") && strings.HasSuffix(part, "]") {
		return processArrayOrCondition(data, part)
	}
	switch v := data.(type) {
	case map[string]interface{}:
		if part == "*" {
			return v, nil
		}
		if value, ok := v[part]; ok {
			return value, nil
		}
		return nil, fmt.Errorf("key not found: %s", part)
	case []interface{}:
		if part == "*" {
			return v, nil
		}
		return nil, fmt.Errorf("array access without index: %s", part)
	default:
		return nil, fmt.Errorf("cannot query further: %v", data)
	}
}

func processAttribute(data interface{}, attr string) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		if value, ok := v[attr]; ok {
			return value, nil
		}
		return nil, fmt.Errorf("attribute not found: %s", attr)
	default:
		return nil, fmt.Errorf("cannot query attribute on non-object: %v", data)
	}
}

func processArrayOrCondition(data interface{}, part string) (interface{}, error) {
	bracketIndex := strings.Index(part, "[")
	nodeName := part[:bracketIndex]
	condition := part[bracketIndex+1 : len(part)-1]

	switch v := data.(type) {
	case map[string]interface{}:
		if nodeName != "" {
			if value, ok := v[nodeName]; ok {
				return processCondition(value, condition)
			}
			return nil, fmt.Errorf("key not found: %s", nodeName)
		}
		return processCondition(v, condition)
	case []interface{}:
		if nodeName != "" {
			return nil, fmt.Errorf("cannot use node name with array: %s", nodeName)
		}
		return processCondition(v, condition)
	default:
		return nil, fmt.Errorf("cannot process array or condition on: %v", data)
	}
}

func processCondition(data interface{}, condition string) (interface{}, error) {
	if index, err := strconv.Atoi(condition); err == nil {
		return processArrayIndex(data, index)
	}
	return processFilterCondition(data, condition)
}

func processArrayIndex(data interface{}, index int) (interface{}, error) {
	switch v := data.(type) {
	case []interface{}:
		if index >= 0 && index < len(v) {
			return v[index], nil
		}
		return nil, fmt.Errorf("array index out of bounds: %d", index)
	default:
		return nil, fmt.Errorf("cannot use array index on non-array: %v", data)
	}
}

func processFilterCondition(data interface{}, condition string) (interface{}, error) {
	re := regexp.MustCompile(`(@?\w+)\s*([=!<>]+)\s*(.+)`)
	matches := re.FindStringSubmatch(condition)
	if matches == nil {
		return nil, fmt.Errorf("invalid filter condition: %s", condition)
	}

	attr, op, value := matches[1], matches[2], matches[3]
	value = strings.Trim(value, "'\"")

	switch v := data.(type) {
	case []interface{}:
		var result []interface{}
		for _, item := range v {
			if matches, _ := evaluateCondition(item, attr, op, value); matches {
				result = append(result, item)
			}
		}
		return result, nil
	case map[string]interface{}:
		if matches, _ := evaluateCondition(v, attr, op, value); matches {
			return v, nil
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("cannot apply filter condition to: %v", data)
	}
}

func evaluateCondition(item interface{}, attr, op, value string) (bool, error) {
	var itemValue interface{}
	if strings.HasPrefix(attr, "@") {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("cannot access attribute on non-object")
		}
		itemValue = itemMap[attr[1:]]
	} else {
		itemValue = item
	}

	switch op {
	case "=":
		return fmt.Sprintf("%v", itemValue) == value, nil
	case "!=":
		return fmt.Sprintf("%v", itemValue) != value, nil
	case ">":
		return compareValues(itemValue, value) > 0, nil
	case ">=":
		return compareValues(itemValue, value) >= 0, nil
	case "<":
		return compareValues(itemValue, value) < 0, nil
	case "<=":
		return compareValues(itemValue, value) <= 0, nil
	default:
		return false, fmt.Errorf("unsupported operator: %s", op)
	}
}

func compareValues(a, b interface{}) int {
	aFloat, aErr := strconv.ParseFloat(fmt.Sprintf("%v", a), 64)
	bFloat, bErr := strconv.ParseFloat(fmt.Sprintf("%v", b), 64)

	if aErr == nil && bErr == nil {
		if aFloat < bFloat {
			return -1
		} else if aFloat > bFloat {
			return 1
		}
		return 0
	}

	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)
	return strings.Compare(aStr, bStr)
}

func printJSON(data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting JSON: %v\n", err)
		return
	}

	coloredJSON := colorizeJSON(string(jsonBytes))
	fmt.Println(coloredJSON)
}

func colorizeJSON(jsonStr string) string {
	var result strings.Builder
	var inString bool
	var inKey bool
	var colonCount int

	for i := 0; i < len(jsonStr); i++ {
		char := rune(jsonStr[i])
		switch {
		case char == '"':
			if i > 0 && jsonStr[i-1] != '\\' {
				inString = !inString
				if !inString {
					colonCount = 0
				}
				if inKey {
					result.WriteString(color.CyanString(string(char)))
				} else {
					result.WriteString(color.GreenString(string(char)))
				}
			}
		case inString:
			if inKey {
				result.WriteString(color.CyanString(string(char)))
			} else {
				result.WriteString(color.GreenString(string(char)))
			}
		case char == ':':
			result.WriteString(color.WhiteString(string(char)))
			colonCount++
			if colonCount == 1 {
				inKey = false
			}
		case char == '{' || char == '}' || char == '[' || char == ']':
			result.WriteString(color.MagentaString(string(char)))
			inKey = true
			colonCount = 0
		case char >= '0' && char <= '9' || char == '.' || char == '-':
			result.WriteString(color.YellowString(string(char)))
		case char == 't' && strings.HasPrefix(jsonStr[i:], "true"):
			result.WriteString(color.BlueString("true"))
			i += 3
		case char == 'f' && strings.HasPrefix(jsonStr[i:], "false"):
			result.WriteString(color.BlueString("false"))
			i += 4
		case char == 'n' && strings.HasPrefix(jsonStr[i:], "null"):
			result.WriteString(color.RedString("null"))
			i += 3
		case char == ',':
			result.WriteString(color.WhiteString(string(char)))
			inKey = true
			colonCount = 0
		default:
			result.WriteRune(char)
		}
	}

	return result.String()
}
