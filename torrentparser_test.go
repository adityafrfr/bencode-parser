package bencode

import (
	"reflect"
	"testing"
)

// ---------- FillTorrentInfo ----------

func TestFillTorrentInfo_AllFieldsPresent(t *testing.T) {
	m := map[string]any{
		"announce":      "http://tracker.example.com/announce",
		"comment":       "a test torrent",
		"created by":    "mktorrent 1.1",
		"creation date": int64(1700000000),
	}
	tf := &TorrentFile{}
	FillTorrentInfo(m, tf)

	if tf.Announce != "http://tracker.example.com/announce" {
		t.Errorf("Announce = %q, want tracker URL", tf.Announce)
	}
	if tf.Comment != "a test torrent" {
		t.Errorf("Comment = %q, want 'a test torrent'", tf.Comment)
	}
	if tf.CreatedBy != "mktorrent 1.1" {
		t.Errorf("CreatedBy = %q, want 'mktorrent 1.1'", tf.CreatedBy)
	}
	if tf.CreationDate != int64(1700000000) {
		t.Errorf("CreationDate = %d, want 1700000000", tf.CreationDate)
	}
}

func TestFillTorrentInfo_MissingFieldsLeaveZeroValues(t *testing.T) {
	m := map[string]any{
		"announce": "http://tracker.example.com/announce",
		// comment, created by, creation date all absent
	}
	tf := &TorrentFile{}
	FillTorrentInfo(m, tf)

	if tf.Announce != "http://tracker.example.com/announce" {
		t.Errorf("Announce = %q, want tracker URL", tf.Announce)
	}
	if tf.Comment != "" {
		t.Errorf("Comment = %q, want empty string", tf.Comment)
	}
	if tf.CreatedBy != "" {
		t.Errorf("CreatedBy = %q, want empty string", tf.CreatedBy)
	}
	if tf.CreationDate != 0 {
		t.Errorf("CreationDate = %d, want 0", tf.CreationDate)
	}
}

func TestFillTorrentInfo_WrongTypeIsIgnoredNotPanicked(t *testing.T) {
	// A malformed/hostile input: "announce" is an int instead of a string,
	// "creation date" is a string instead of an int64. Should not panic,
	// and the mistyped fields should simply stay at their zero value.
	m := map[string]any{
		"announce":      int64(123),
		"creation date": "not-a-number",
	}
	tf := &TorrentFile{}
	FillTorrentInfo(m, tf)

	if tf.Announce != "" {
		t.Errorf("Announce = %q, want empty string (wrong type should be ignored)", tf.Announce)
	}
	if tf.CreationDate != 0 {
		t.Errorf("CreationDate = %d, want 0 (wrong type should be ignored)", tf.CreationDate)
	}
}

func TestFillTorrentInfo_EmptyMap(t *testing.T) {
	tf := &TorrentFile{}
	FillTorrentInfo(map[string]any{}, tf)

	want := &TorrentFile{}
	if !reflect.DeepEqual(tf, want) {
		t.Errorf("got %+v, want zero-value TorrentFile", tf)
	}
}

// ---------- FillRawInfo ----------

func TestFillRawInfo_SingleFileTorrent(t *testing.T) {
	rawInfo := map[string]any{
		"name":         "ubuntu.iso",
		"piece length": int64(262144),
		"pieces":       "\x01\x02\x03\x04", // stand-in for concatenated SHA1 hashes
		"length":       int64(4700372992),
	}
	tf := &TorrentFile{}
	FillRawInfo(rawInfo, tf)

	if tf.Info.Name != "ubuntu.iso" {
		t.Errorf("Info.Name = %q, want 'ubuntu.iso'", tf.Info.Name)
	}
	if tf.Info.PieceLength != int64(262144) {
		t.Errorf("Info.PieceLength = %d, want 262144", tf.Info.PieceLength)
	}
	if tf.Info.Pieces != "\x01\x02\x03\x04" {
		t.Errorf("Info.Pieces = %q, want raw hash bytes", tf.Info.Pieces)
	}
	if tf.Info.Length != int64(4700372992) {
		t.Errorf("Info.Length = %d, want 4700372992", tf.Info.Length)
	}
}

func TestFillRawInfo_MissingLengthStaysZero(t *testing.T) {
	// Multi-file torrents omit "length" at the info level (it lives per-file
	// instead), so this must not panic and must leave Length at 0.
	rawInfo := map[string]any{
		"name":         "my_folder",
		"piece length": int64(16384),
	}
	tf := &TorrentFile{}
	FillRawInfo(rawInfo, tf)

	if tf.Info.Name != "my_folder" {
		t.Errorf("Info.Name = %q, want 'my_folder'", tf.Info.Name)
	}
	if tf.Info.Length != 0 {
		t.Errorf("Info.Length = %d, want 0 when absent", tf.Info.Length)
	}
}

func TestFillRawInfo_OverwritesExistingInfo(t *testing.T) {
	// tf.Info already has data; FillRawInfo should replace it wholesale,
	// not merge into it (matches the current "info := InfoDict{}; ...;
	// tf.Info = info" implementation).
	tf := &TorrentFile{
		Info: InfoDict{Name: "stale", Length: 999},
	}
	rawInfo := map[string]any{
		"name": "fresh",
	}
	FillRawInfo(rawInfo, tf)

	if tf.Info.Name != "fresh" {
		t.Errorf("Info.Name = %q, want 'fresh'", tf.Info.Name)
	}
	if tf.Info.Length != 0 {
		t.Errorf("Info.Length = %d, want 0 (stale data should not leak through)", tf.Info.Length)
	}
}

// ---------- FillAnnounceList ----------

func TestFillAnnounceList_MultipleTiers(t *testing.T) {
	m := map[string]any{
		"announce-list": []any{
			[]any{"http://tracker1.example.com/a"},
			[]any{"http://tracker2.example.com/a", "http://tracker3.example.com/a"},
		},
	}
	tf := &TorrentFile{}
	FillAnnounceList(m, tf)

	want := [][]string{
		{"http://tracker1.example.com/a"},
		{"http://tracker2.example.com/a", "http://tracker3.example.com/a"},
	}
	if !reflect.DeepEqual(tf.AnnounceList, want) {
		t.Errorf("AnnounceList = %v, want %v", tf.AnnounceList, want)
	}
}

func TestFillAnnounceList_Absent(t *testing.T) {
	tf := &TorrentFile{}
	FillAnnounceList(map[string]any{}, tf)

	if tf.AnnounceList != nil {
		t.Errorf("AnnounceList = %v, want nil when key absent", tf.AnnounceList)
	}
}

func TestFillAnnounceList_EmptyTierList(t *testing.T) {
	m := map[string]any{
		"announce-list": []any{},
	}
	tf := &TorrentFile{}
	FillAnnounceList(m, tf)

	if tf.AnnounceList != nil {
		t.Errorf("AnnounceList = %v, want nil for an empty announce-list", tf.AnnounceList)
	}
}

func TestFillAnnounceList_WrongShapeIsSkipped(t *testing.T) {
	// announce-list present but not a list of lists (e.g. flat list of strings) -
	// should be ignored rather than panicking.
	m := map[string]any{
		"announce-list": []any{"not-a-tier"},
	}
	tf := &TorrentFile{}
	FillAnnounceList(m, tf)

	if tf.AnnounceList != nil {
		t.Errorf("AnnounceList = %v, want nil (malformed tier should be skipped)", tf.AnnounceList)
	}
}

// ---------- FillFiles ----------

func TestFillFiles_MultiFileTorrent(t *testing.T) {
	rawInfo := map[string]any{
		"files": []any{
			map[string]any{
				"length": int64(1000),
				"path":   []any{"subdir", "a.txt"},
			},
			map[string]any{
				"length": int64(2000),
				"path":   []any{"b.txt"},
			},
		},
	}
	tf := &TorrentFile{}
	FillFiles(rawInfo, tf)

	want := []File{
		{Length: 1000, Path: []string{"subdir", "a.txt"}},
		{Length: 2000, Path: []string{"b.txt"}},
	}
	if !reflect.DeepEqual(tf.Info.Files, want) {
		t.Errorf("Info.Files = %+v, want %+v", tf.Info.Files, want)
	}
}

func TestFillFiles_Absent(t *testing.T) {
	// Single-file torrents have no "files" key at all.
	tf := &TorrentFile{}
	FillFiles(map[string]any{}, tf)

	if tf.Info.Files != nil {
		t.Errorf("Info.Files = %v, want nil when 'files' absent", tf.Info.Files)
	}
}

func TestFillFiles_AppendsToExistingFiles(t *testing.T) {
	// FillFiles appends onto tf.Info.Files rather than replacing it, so
	// calling it after FillRawInfo must not clobber other Info fields,
	// and calling it twice should accumulate.
	tf := &TorrentFile{Info: InfoDict{Name: "keepme"}}
	rawInfo := map[string]any{
		"files": []any{
			map[string]any{"length": int64(5), "path": []any{"one.txt"}},
		},
	}
	FillFiles(rawInfo, tf)

	if tf.Info.Name != "keepme" {
		t.Errorf("Info.Name = %q, want 'keepme' (FillFiles must not clobber other Info fields)", tf.Info.Name)
	}
	if len(tf.Info.Files) != 1 {
		t.Fatalf("len(Info.Files) = %d, want 1", len(tf.Info.Files))
	}
}

func TestFillFiles_MissingPathDefaultsToNil(t *testing.T) {
	rawInfo := map[string]any{
		"files": []any{
			map[string]any{"length": int64(42)},
		},
	}
	tf := &TorrentFile{}
	FillFiles(rawInfo, tf)

	if len(tf.Info.Files) != 1 {
		t.Fatalf("len(Info.Files) = %d, want 1", len(tf.Info.Files))
	}
	if tf.Info.Files[0].Length != 42 {
		t.Errorf("Files[0].Length = %d, want 42", tf.Info.Files[0].Length)
	}
	if tf.Info.Files[0].Path != nil {
		t.Errorf("Files[0].Path = %v, want nil when 'path' absent", tf.Info.Files[0].Path)
	}
}

// ---------- MapToTorrent ----------

func TestMapToTorrent_SingleFileTorrent(t *testing.T) {
	m := map[string]any{
		"announce": "http://tracker.example.com/announce",
		"comment":  "single file test",
		"info": map[string]any{
			"name":         "ubuntu.iso",
			"piece length": int64(262144),
			"pieces":       "hash-bytes",
			"length":       int64(1000000),
		},
	}

	tf, err := MapToTorrent(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tf.Announce != "http://tracker.example.com/announce" {
		t.Errorf("Announce = %q", tf.Announce)
	}
	if tf.Info.Name != "ubuntu.iso" {
		t.Errorf("Info.Name = %q", tf.Info.Name)
	}
	if tf.Info.Length != 1000000 {
		t.Errorf("Info.Length = %d, want 1000000", tf.Info.Length)
	}
	if tf.Info.Files != nil {
		t.Errorf("Info.Files = %v, want nil for a single-file torrent", tf.Info.Files)
	}
}

func TestMapToTorrent_MultiFileTorrent(t *testing.T) {
	m := map[string]any{
		"announce": "http://tracker.example.com/announce",
		"info": map[string]any{
			"name":         "my_folder",
			"piece length": int64(16384),
			"pieces":       "hash-bytes",
			"files": []any{
				map[string]any{
					"length": int64(100),
					"path":   []any{"a.txt"},
				},
				map[string]any{
					"length": int64(200),
					"path":   []any{"subdir", "b.txt"},
				},
			},
		},
	}

	tf, err := MapToTorrent(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tf.Info.Length != 0 {
		t.Errorf("Info.Length = %d, want 0 for a multi-file torrent", tf.Info.Length)
	}
	wantFiles := []File{
		{Length: 100, Path: []string{"a.txt"}},
		{Length: 200, Path: []string{"subdir", "b.txt"}},
	}
	if !reflect.DeepEqual(tf.Info.Files, wantFiles) {
		t.Errorf("Info.Files = %+v, want %+v", tf.Info.Files, wantFiles)
	}
}

func TestMapToTorrent_MissingInfoDict(t *testing.T) {
	m := map[string]any{
		"announce": "http://tracker.example.com/announce",
		// no "info" key at all
	}

	tf, err := MapToTorrent(m)
	if err == nil {
		t.Fatal("expected error for missing 'info' dictionary, got nil")
	}
	if tf != nil {
		t.Errorf("expected nil *TorrentFile on error, got %+v", tf)
	}
}

func TestMapToTorrent_InfoWrongType(t *testing.T) {
	m := map[string]any{
		"announce": "http://tracker.example.com/announce",
		"info":     "not-a-dict",
	}

	tf, err := MapToTorrent(m)
	if err == nil {
		t.Fatal("expected error for 'info' not being a dictionary, got nil")
	}
	if tf != nil {
		t.Errorf("expected nil *TorrentFile on error, got %+v", tf)
	}
}

func TestMapToTorrent_FullFieldSet(t *testing.T) {
	m := map[string]any{
		"announce": "http://primary.example.com/announce",
		"announce-list": []any{
			[]any{"http://primary.example.com/announce"},
			[]any{"http://backup.example.com/announce"},
		},
		"comment":       "full field test",
		"created by":    "gotorrent 0.1",
		"creation date": int64(1690000000),
		"info": map[string]any{
			"name":         "file.bin",
			"piece length": int64(32768),
			"pieces":       "abc",
			"length":       int64(555),
		},
	}

	tf, err := MapToTorrent(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := &TorrentFile{
		Announce: "http://primary.example.com/announce",
		AnnounceList: [][]string{
			{"http://primary.example.com/announce"},
			{"http://backup.example.com/announce"},
		},
		Comment:      "full field test",
		CreatedBy:    "gotorrent 0.1",
		CreationDate: 1690000000,
		Info: InfoDict{
			Name:        "file.bin",
			PieceLength: 32768,
			Pieces:      "abc",
			Length:      555,
		},
	}

	if !reflect.DeepEqual(tf, want) {
		t.Errorf("got %+v\nwant %+v", tf, want)
	}
}

// ---------- DecodeTorrent (full pipeline: raw bencode bytes -> *TorrentFile) ----------

func TestDecodeTorrent_SingleFile(t *testing.T) {
	input := "d8:announce16:http://tracker/a7:comment4:test4:infod" +
		"6:lengthi1000e4:name5:file112:piece lengthi16384e6:pieces4:abcde" +
		"e"

	tf, err := DecodeTorrent(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tf.Announce != "http://tracker/a" {
		t.Errorf("Announce = %q", tf.Announce)
	}
	if tf.Comment != "test" {
		t.Errorf("Comment = %q", tf.Comment)
	}
	if tf.Info.Name != "file1" {
		t.Errorf("Info.Name = %q", tf.Info.Name)
	}
	if tf.Info.Length != 1000 {
		t.Errorf("Info.Length = %d, want 1000", tf.Info.Length)
	}
	if tf.Info.PieceLength != 16384 {
		t.Errorf("Info.PieceLength = %d, want 16384", tf.Info.PieceLength)
	}
}

func TestDecodeTorrent_WithAnnounceListAndTrailingKeys(t *testing.T) {
	// This is the scenario that exposed the parseList consumed-count bug:
	// announce-list (a list of lists) is followed by more top-level keys.
	// If that bug were still present, "comment" and "info" below would
	// silently disappear from the decoded map.
	input := "d13:announce-listll17:http://tracker1/ael17:http://tracker2/aee" +
		"7:comment4:test4:infod4:name1:x12:piece lengthi1ee" +
		"e"

	tf, err := DecodeTorrent(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantList := [][]string{
		{"http://tracker1/a"},
		{"http://tracker2/a"},
	}
	if !reflect.DeepEqual(tf.AnnounceList, wantList) {
		t.Errorf("AnnounceList = %v, want %v", tf.AnnounceList, wantList)
	}
	if tf.Comment != "test" {
		t.Errorf("Comment = %q, want 'test' (would be empty if the list-consumed bug regressed)", tf.Comment)
	}
	if tf.Info.Name != "x" {
		t.Errorf("Info.Name = %q, want 'x'", tf.Info.Name)
	}
}

func TestDecodeTorrent_MultiFileWithFilesList(t *testing.T) {
	input := "d8:announce11:http://tr/a4:infod4:name6:folder" +
		"12:piece lengthi16384e5:filesld6:lengthi10e4:pathl5:a.txtee" +
		"d6:lengthi20e4:pathl6:subdir5:b.txteeeee"

	tf, err := DecodeTorrent(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantFiles := []File{
		{Length: 10, Path: []string{"a.txt"}},
		{Length: 20, Path: []string{"subdir", "b.txt"}},
	}
	if !reflect.DeepEqual(tf.Info.Files, wantFiles) {
		t.Errorf("Info.Files = %+v, want %+v", tf.Info.Files, wantFiles)
	}
}

func TestDecodeTorrent_NotADictionaryAtRoot(t *testing.T) {
	// Root element is a bencoded string, not a dictionary.
	_, err := DecodeTorrent("4:spam")
	if err == nil {
		t.Fatal("expected error when root element is not a dictionary, got nil")
	}
}

func TestDecodeTorrent_MalformedBencode(t *testing.T) {
	_, err := DecodeTorrent("d8:announce")
	if err == nil {
		t.Fatal("expected error for truncated/malformed bencode input, got nil")
	}
}

func TestDecodeTorrent_EmptyInput(t *testing.T) {
	_, err := DecodeTorrent("")
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}
