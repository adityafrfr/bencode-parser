package bencode

import (
	"bytes"
	"fmt"
	"io"
)

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func Unmarshal(data []byte) (any, error) {
	return NewDecoder(bytes.NewReader(data)).Decode()
}

func (d *Decoder) Decode() (any, error) {
	data, err := io.ReadAll(d.r)

	if err != nil {
		return nil, fmt.Errorf("bencode: failed to read input: %w", err)
	}
	if len(data) == 0 {
		return nil, io.EOF
	}

	val, consumed, err := parseNext(string(data), 0)

	if err != nil {
		return nil, err
	}
	if consumed < len(data) {
		return nil, fmt.Errorf("bencode: unexpected trailing data at the offset %d", consumed)
	}

	return val, nil
}
