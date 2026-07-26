package bencode

import(
	"testing"
	"fmt"
)

func TestScratch(t *testing.T) {
	v, n, err := parseNext("i42e", 0)
	fmt.Println(v, n, err) // want: 42 4 <nil>

	v, n, err = parseNext("4:spam", 0)
	fmt.Println(v, n, err) // want: spam 6 <nil>
}