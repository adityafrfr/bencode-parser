package bencode

import "fmt"

// THIS FUNCTION CAN PUT ANYTHING INSIDE THE RETURN BOX
func parseNext(s string, offset int) (any, int, error)	{
if offset >= len(s) {
		return nil, 0, fmt.Errorf("unexpected end of input")
	}

	switch s[offset] {
	case 'i':
		// TODO integer
	case 'l':
		// TODO list
	case 'd':
		// TODO dict
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// TODO string
	default:
		return nil, 0, fmt.Errorf("unexpected byte: %c", s[offset])
	}
}