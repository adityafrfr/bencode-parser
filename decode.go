package bencode

import (
	"fmt"
	"strconv"
	"strings"
)

func parseInt(s string, offset int) (any, int, error) {

	end := strings.IndexByte(s[offset:], 'e')
	if end == -1 {
		return nil, 0, fmt.Errorf("invalid int, missing the terminating 'e'")
	}

	numStr := s[offset+1 : offset+end]
	val, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return nil, 0, fmt.Errorf("Invalid integer: %w", err)
	}

	return val, end + 1, nil
}

func parseList(offset int, s string) (any, int, error) {
	current := offset + 1

	var list []any

	for current < len(s) && s[current] != 'e' {
		item, consumed, err := parseNext(s, current)
		if err != nil {
			return nil, 0, err
		}

		list = append(list, item)
		current += consumed
	}

	if current >= len(s) || s[current] != 'e' {
		return nil, 0, fmt.Errorf("invalid list: missing 'e'")
	}

	return list, current - offset, nil
}

func parseDict(offset int, s string) (any, int, error) {
	curr := offset + 1
	dict := make(map[string]any)
	for curr < len(s) && s[curr] != 'e' {
		keyVal, consumedKey, err := parseNext(s, curr)
		if err != nil {
			return nil, 0, err
		}
		key, ok := keyVal.(string)
		if !ok {
			return nil, 0, fmt.Errorf("dict key must be a string")
		}
		curr += consumedKey

		val, consumedVal, err := parseNext(s, curr)
		if err != nil {
			return nil, 0, err
		}
		dict[key] = val
		curr += consumedVal
	}
	if curr >= len(s) || s[curr] != 'e' {
		return nil, 0, fmt.Errorf("invalid dict: missing 'e'")
	}
	return dict, (curr - offset) + 1, nil
}

func parseString(s string, offset int) (any, int, error) {
	colon := strings.IndexByte(s[offset:], ':')
	if colon == -1 {
		return nil, 0, fmt.Errorf("invalid string: missing ':'")
	}

	length, err := strconv.Atoi(s[offset : offset+colon])
	if err != nil || length < 0 {
		return nil, 0, fmt.Errorf("invalid string length")
	}

	start := offset + colon + 1
	end := start + length

	if end > len(s) {
		return nil, 0, fmt.Errorf("string length exceeds input size")
	}

	return s[start:end], colon + length + 1, nil
}

// THIS FUNCTION CAN PUT ANYTHING INSIDE THE RETURN BOX
func parseNext(s string, offset int) (any, int, error) {
	if offset >= len(s) {
		return nil, 0, fmt.Errorf("unexpected end of input")
	}

	switch s[offset] {

	case 'i':
		return parseInt(s, offset)

	case 'l':
		return parseList(offset, s)

	case 'd':
		return parseDict(offset, s)

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return parseString(s, offset)

	default:
		return nil, 0, fmt.Errorf("unexpected byte: %c", s[offset])
	}
}
