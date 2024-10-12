package json

import (
	"fmt"
	"strings"
)

// KeyNotFoundError is returned when a key is not found in an object.
type KeyNotFoundError struct {
	Key string
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("key not found: %s", e.Key)
}

func (e *KeyNotFoundError) Is(target error) bool {
	_, ok := target.(*KeyNotFoundError)
	return ok
}

// InvalidPathError is returned when the provided path is invalid.
type InvalidPathError struct {
	Path []string
}

func (e *InvalidPathError) Error() string {
	return fmt.Sprintf("invalid path: %s", strings.Join(e.Path, "."))
}

func (e *InvalidPathError) Is(target error) bool {
	_, ok := target.(*InvalidPathError)
	return ok
}

// IndexOutOfRangeError is returned when an array index is out of range.
type IndexOutOfRangeError struct {
	Index  int
	Length int
}

func (e *IndexOutOfRangeError) Error() string {
	return fmt.Sprintf("index out of range: index %d, length %d", e.Index, e.Length)
}

func (e *IndexOutOfRangeError) Is(target error) bool {
	_, ok := target.(*IndexOutOfRangeError)
	return ok
}

// ParseError is returned when there's an error parsing the JSON input.
type ParseError struct {
	Line   int
	Column int
	Msg    string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at line %d, column %d: %s", e.Line, e.Column, e.Msg)
}

func (e *ParseError) Is(target error) bool {
	_, ok := target.(*ParseError)
	return ok
}

// MaxDepthExceededError is returned when the JSON structure exceeds the maximum allowed depth.
type MaxDepthExceededError struct {
	MaxDepth int
}

func (e *MaxDepthExceededError) Error() string {
	return fmt.Sprintf("maximum nesting depth of %d exceeded", e.MaxDepth)
}

func (e *MaxDepthExceededError) Is(target error) bool {
	_, ok := target.(*MaxDepthExceededError)
	return ok
}

// ... 其他错误类型
