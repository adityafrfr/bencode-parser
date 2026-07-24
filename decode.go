package bencode

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func decodeInt(input string) (int64, error) {
	numStr := input[1 : len(input)-1]

	val, err := strconv.Atoi(numStr)

	if err != nil {
		return -1, fmt.Errorf("got error %q while trying to convert %v to int", err, numStr)
	}

	return int64(val), nil
}

func decodeString(input string) (string, error) {
	colonIdx := strings.IndexByte(input, ':')

	length, err := strconv.Atoi(input[:colonIdx])

	if err != nil {
		log.Fatalf("got err %q while trying to extract string %q", err, input)
		return "", err
	}

	start := colonIdx + 1

	val := input[start : start+length]

	return val, nil
}

func decodeList(input string) ([]interface{}, int, error) {
	if (input[0] != 'l')	{
		return nil,0, fmt.Errorf("invalid list start")
	}
}
