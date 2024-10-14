# JSON Parser with XPath-like Query Support

[![Go Version](https://img.shields.io/badge/Go-1.16%2B-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

A powerful command-line tool for parsing, querying, and displaying JSON data with colorized output. This tool provides an intuitive way to explore and extract information from complex JSON structures using XPath-like queries.

## Features

- **Parse and pretty-print** JSON files with colorized output.
- **XPath-like query support** for JSON data.
- Support for various **comparison operators** in queries.
- Easy-to-use **command-line interface**.
- Efficient handling of **large JSON files**.

## Installation

### From Source

1. Ensure you have Go installed (version 1.16 or later).
2. Clone the repository:
   ```bash
   git clone https://github.com/davidhoo/jsonparser.git
   ```
3. Navigate to the project directory:
   ```bash
   cd jsonparser
   ```
4. Install dependencies:
   ```bash
   go mod tidy
   ```
5. Build the project:
   ```bash
   go build -o jp main.go
   ```

### Homebrew Installation (Experimental)

You can also install this tool using Homebrew by tapping the repository. This method is experimental and does not require submission to the Homebrew core repository.

1. Tap the repository:
   ```bash
   brew tap davidhoo/jsonparser
   ```
2. Install the tool:
   ```bash
   brew install jsonparser
   ```

## Usage

Basic syntax:

```bash
./jp [-f <json_file>] [-q <query>]
```

or

```bash
./jp <json_file> [-q <query>]
```

### Options

- `-f <json_file>`: Specify the JSON file path.
- `-q <query>`: XPath-like query string to filter JSON.
- `-h`: Show help message.

If no query is provided, the entire JSON will be printed with colorized output.

## Query Syntax

- Use `/` to separate path elements.
- Use `@` to access object properties.
- Use `[]` for array indexing or filtering.
- Use `*` as a wildcard to select all elements.

### Supported Operators

- **Equality**: `=`
- **Inequality**: `!=`
- **Greater than**: `>`
- **Greater than or equal to**: `>=`
- **Less than**: `<`
- **Less than or equal to**: `<=`

## Query Examples

1. Get the first user:
   ```bash
   ./jp -f complex.json -q "/data/users[0]"
   ```

2. Find user with name 'Alice':
   ```bash
   ./jp -f complex.json -q "/data/users[@name='Alice']"
   ```

3. Find products with price over 1000:
   ```bash
   ./jp -f complex.json -q "/data/products[price>1000]"
   ```

4. Get all notification settings:
   ```bash
   ./jp -f complex.json -q "/settings/notifications/*"
   ```

## Color Scheme

The tool uses the following color scheme for JSON elements:

- **Keys**: Cyan
- **String values**: Green
- **Number values**: Yellow
- **Boolean values**: Blue
- **Null values**: Red
- **Brackets and braces**: Magenta
- **Colons and commas**: White

## Error Handling

The tool provides informative error messages for:

- File reading errors.
- JSON parsing errors.
- Query execution errors.
- Invalid query syntax.

## Performance

The tool is designed to handle large JSON files efficiently. However, performance may vary depending on the complexity of the query and the size of the JSON data.

## Limitations

- The tool currently does not support writing modified JSON back to a file.
- Some advanced XPath features (like functions) are not implemented.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgements

- This project uses the [github.com/fatih/color](https://github.com/fatih/color) library for colorized output.
