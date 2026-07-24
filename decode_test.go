package bencode

import (
	"log"
	"testing"
)

func TestEncodeInt(t *testing.T) {
	t.Run("decode simple integers", func(t *testing.T) {

		var toDecode string = "i42e"
		got, err := decodeInt(toDecode)
		want := 42

		if err != nil {
			log.Fatalf("got %q while trying to convert %v to int", err, toDecode)
		}

		if got != int64(want) {
			log.Fatalf("got %v want %v", got, want)
		}

	})
}

func TestEncodeString(t *testing.T) {

	t.Run("decode simple string", func(t *testing.T) {

		var toDecode string = "9:whatdafaq"
		got, err := decodeString(toDecode)
		want := "whatdafaq"

		if err != nil {
			log.Fatalf("got %q while trying to convert %v to string", err, toDecode)
		}

		if got != want {
			log.Fatalf("got %v want %v", got, want)
		}

	})
}
