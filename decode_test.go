package bencode

import (
	"reflect"
	"testing"
)

// ---------- parseInt ----------

func TestParseInt_Positive(t *testing.T) {
	val, consumed, err := parseInt("i42e", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(42) {
		t.Errorf("got %v, want int64(42)", val)
	}
	if consumed != 4 {
		t.Errorf("consumed = %d, want 4", consumed)
	}
}

func TestParseInt_Negative(t *testing.T) {
	val, consumed, err := parseInt("i-13e", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(-13) {
		t.Errorf("got %v, want int64(-13)", val)
	}
	if consumed != 5 {
		t.Errorf("consumed = %d, want 5", consumed)
	}
}

func TestParseInt_Zero(t *testing.T) {
	val, consumed, err := parseInt("i0e", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(0) {
		t.Errorf("got %v, want int64(0)", val)
	}
	if consumed != 3 {
		t.Errorf("consumed = %d, want 3", consumed)
	}
}

func TestParseInt_LargeValue(t *testing.T) {
	// Bigger than int32, must fit in int64.
	val, _, err := parseInt("i9223372036854775807e", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(9223372036854775807) {
		t.Errorf("got %v, want max int64", val)
	}
}

func TestParseInt_MissingTerminator(t *testing.T) {
	_, _, err := parseInt("i42", 0)
	if err == nil {
		t.Fatal("expected error for missing 'e', got nil")
	}
}

func TestParseInt_NotANumber(t *testing.T) {
	_, _, err := parseInt("iabce", 0)
	if err == nil {
		t.Fatal("expected error for non-numeric content, got nil")
	}
}

func TestParseInt_Empty(t *testing.T) {
	_, _, err := parseInt("ie", 0)
	if err == nil {
		t.Fatal("expected error for empty integer 'ie', got nil")
	}
}

func TestParseInt_AtOffset(t *testing.T) {
	// Make sure offset handling works when not starting at position 0.
	s := "xxxi99e"
	val, consumed, err := parseInt(s, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(99) {
		t.Errorf("got %v, want int64(99)", val)
	}
	if consumed != 4 {
		t.Errorf("consumed = %d, want 4", consumed)
	}
}

// ---------- parseString ----------

func TestParseString_Basic(t *testing.T) {
	val, consumed, err := parseString("4:spam", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "spam" {
		t.Errorf("got %v, want 'spam'", val)
	}
	if consumed != 6 {
		t.Errorf("consumed = %d, want 6", consumed)
	}
}

func TestParseString_Empty(t *testing.T) {
	val, consumed, err := parseString("0:", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "" {
		t.Errorf("got %v, want empty string", val)
	}
	if consumed != 2 {
		t.Errorf("consumed = %d, want 2", consumed)
	}
}

func TestParseString_MissingColon(t *testing.T) {
	_, _, err := parseString("4spam", 0)
	if err == nil {
		t.Fatal("expected error for missing ':', got nil")
	}
}

func TestParseString_LengthExceedsInput(t *testing.T) {
	_, _, err := parseString("10:abc", 0)
	if err == nil {
		t.Fatal("expected error for string length exceeding input size, got nil")
	}
}

func TestParseString_InvalidLength(t *testing.T) {
	_, _, err := parseString("abc:def", 0)
	if err == nil {
		t.Fatal("expected error for non-numeric length prefix, got nil")
	}
}

func TestParseString_AtOffset(t *testing.T) {
	s := "yy3:cat"
	val, consumed, err := parseString(s, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "cat" {
		t.Errorf("got %v, want 'cat'", val)
	}
	if consumed != 5 {
		t.Errorf("consumed = %d, want 5", consumed)
	}
}

func TestParseString_BinaryContent(t *testing.T) {
	// pieces fields hold raw (non-UTF8-safe) bytes; make sure we don't
	// assume printable text anywhere in the string path.
	raw := "\x00\x01\xff\xfe"
	s := "4:" + raw
	val, _, err := parseString(s, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != raw {
		t.Errorf("got %q, want %q", val, raw)
	}
}

// ---------- parseList ----------

func TestParseList_Empty(t *testing.T) {
	val, consumed, err := parseList(0, "le")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list, ok := val.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", val)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %v", list)
	}
	if consumed != 2 {
		t.Errorf("consumed = %d, want 2 (must include the closing 'e')", consumed)
	}
}

func TestParseList_Strings(t *testing.T) {
	val, consumed, err := parseList(0, "l4:spam4:eggse")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []any{"spam", "eggs"}
	if !reflect.DeepEqual(val, want) {
		t.Errorf("got %v, want %v", val, want)
	}
	// "l4:spam4:eggse" is 14 bytes total; consumed must cover the whole
	// construct including the trailing 'e', or callers (parseDict/parseList)
	// mis-advance their cursor past this list and corrupt further parsing.
	if consumed != 14 {
		t.Errorf("consumed = %d, want 14", consumed)
	}
}

func TestParseList_Mixed(t *testing.T) {
	val, _, err := parseList(0, "l4:spami42ee")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []any{"spam", int64(42)}
	if !reflect.DeepEqual(val, want) {
		t.Errorf("got %v, want %v", val, want)
	}
}

func TestParseList_Nested(t *testing.T) {
	val, _, err := parseList(0, "ll4:spamee")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []any{[]any{"spam"}}
	if !reflect.DeepEqual(val, want) {
		t.Errorf("got %v, want %v", val, want)
	}
}

func TestParseList_MissingTerminator(t *testing.T) {
	_, _, err := parseList(0, "l4:spam")
	if err == nil {
		t.Fatal("expected error for missing 'e', got nil")
	}
}

func TestParseList_InvalidElement(t *testing.T) {
	_, _, err := parseList(0, "l4:spXe")
	if err == nil {
		t.Fatal("expected error for malformed element, got nil")
	}
}

// ---------- parseDict ----------

func TestParseDict_Empty(t *testing.T) {
	val, consumed, err := parseDict(0, "de")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dict, ok := val.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", val)
	}
	if len(dict) != 0 {
		t.Errorf("expected empty dict, got %v", dict)
	}
	if consumed != 2 {
		t.Errorf("consumed = %d, want 2", consumed)
	}
}

func TestParseDict_SingleKey(t *testing.T) {
	val, consumed, err := parseDict(0, "d3:cow3:mooe")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]any{"cow": "moo"}
	if !reflect.DeepEqual(val, want) {
		t.Errorf("got %v, want %v", val, want)
	}
	if consumed != 12 {
		t.Errorf("consumed = %d, want 12", consumed)
	}
}

func TestParseDict_MultipleKeys(t *testing.T) {
	val, _, err := parseDict(0, "d3:cow3:moo4:spam4:eggse")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]any{"cow": "moo", "spam": "eggs"}
	if !reflect.DeepEqual(val, want) {
		t.Errorf("got %v, want %v", val, want)
	}
}

func TestParseDict_NestedDict(t *testing.T) {
	// Inner dict needs its own closing 'e', then the outer dict needs one too.
	val, _, err := parseDict(0, "d4:infod4:name4:testee")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]any{
		"info": map[string]any{"name": "test"},
	}
	if !reflect.DeepEqual(val, want) {
		t.Errorf("got %v, want %v", val, want)
	}
}

func TestParseDict_ValueIsList(t *testing.T) {
	val, _, err := parseDict(0, "d4:listl1:a1:bee")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]any{
		"list": []any{"a", "b"},
	}
	if !reflect.DeepEqual(val, want) {
		t.Errorf("got %v, want %v", val, want)
	}
}

func TestParseDict_NonStringKey(t *testing.T) {
	// key encoded as an integer instead of a string -> must error
	_, _, err := parseDict(0, "di1e3:fooe")
	if err == nil {
		t.Fatal("expected error for non-string dict key, got nil")
	}
}

// Regression test for a real bug found while testing this package:
// parseList's returned "consumed" count did not include the list's own
// closing 'e'. That meant any dict (or list) containing a list value,
// followed by MORE keys/elements after it, would have its cursor land
// one byte short — landing back on the list's own closing 'e' instead of
// moving past it. That stray 'e' then looked like the *parent's* own
// terminator, so the parent dict silently stopped early and every key
// after the list-valued one was dropped, with no error raised.
// This is exactly the shape of a real torrent file: "announce-list" is a
// list of lists, and "comment"/"info"/etc. commonly follow it.
func TestParseDict_KeyAfterListValue_Regression(t *testing.T) {
	s := "d4:listl4:spame3:cow3:mooe"
	val, consumed, err := parseDict(0, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]any{
		"list": []any{"spam"},
		"cow":  "moo",
	}
	if !reflect.DeepEqual(val, want) {
		t.Errorf("got %v, want %v (a key after a list value was likely dropped)", val, want)
	}
	if consumed != len(s) {
		t.Errorf("consumed = %d, want %d (entire input)", consumed, len(s))
	}
}

func TestParseDict_MissingTerminator(t *testing.T) {
	_, _, err := parseDict(0, "d3:cow3:moo")
	if err == nil {
		t.Fatal("expected error for missing 'e', got nil")
	}
}

func TestParseDict_DanglingKeyNoValue(t *testing.T) {
	_, _, err := parseDict(0, "d3:cowe")
	if err == nil {
		t.Fatal("expected error when key has no matching value, got nil")
	}
}

// ---------- parseNext (dispatch) ----------

func TestParseNext_DispatchesToInt(t *testing.T) {
	val, _, err := parseNext("i5e", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(5) {
		t.Errorf("got %v, want int64(5)", val)
	}
}

func TestParseNext_DispatchesToString(t *testing.T) {
	val, _, err := parseNext("3:cat", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "cat" {
		t.Errorf("got %v, want 'cat'", val)
	}
}

func TestParseNext_DispatchesToList(t *testing.T) {
	val, _, err := parseNext("l1:ae", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(val, []any{"a"}) {
		t.Errorf("got %v, want [a]", val)
	}
}

func TestParseNext_DispatchesToDict(t *testing.T) {
	val, _, err := parseNext("d1:a1:be", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(val, map[string]any{"a": "b"}) {
		t.Errorf("got %v, want map[a:b]", val)
	}
}

func TestParseNext_UnexpectedByte(t *testing.T) {
	_, _, err := parseNext("x", 0)
	if err == nil {
		t.Fatal("expected error for unrecognized leading byte, got nil")
	}
}

func TestParseNext_EmptyInput(t *testing.T) {
	_, _, err := parseNext("", 0)
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestParseNext_OffsetPastEnd(t *testing.T) {
	_, _, err := parseNext("i5e", 10)
	if err == nil {
		t.Fatal("expected error for offset beyond input length, got nil")
	}
}

// ---------- Realistic end-to-end style case ----------

func TestParseNext_RealisticTorrentLikeDict(t *testing.T) {
	// A miniature but structurally realistic bencoded torrent dict.
	input := "d8:announce16:http://tracker/a7:comment4:test4:infod4:name5:file1" +
		"12:piece lengthi16384e6:lengthi1000eee"

	val, consumed, err := parseNext(input, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if consumed != len(input) {
		t.Errorf("consumed = %d, want %d (entire input)", consumed, len(input))
	}

	dict, ok := val.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", val)
	}

	if dict["announce"] != "http://tracker/a" {
		t.Errorf("announce = %v, want http://tracker/a", dict["announce"])
	}
	if dict["comment"] != "test" {
		t.Errorf("comment = %v, want test", dict["comment"])
	}

	info, ok := dict["info"].(map[string]any)
	if !ok {
		t.Fatalf("expected info to be map[string]any, got %T", dict["info"])
	}
	if info["name"] != "file1" {
		t.Errorf("info.name = %v, want file1", info["name"])
	}
	if info["piece length"] != int64(16384) {
		t.Errorf("info.piece length = %v, want 16384", info["piece length"])
	}
	if info["length"] != int64(1000) {
		t.Errorf("info.length = %v, want 1000", info["length"])
	}
}
