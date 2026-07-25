package bencode

import (
	"reflect"
	"testing"
)

func TestDecodeTorrent_ValidCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *TorrentFile
	}{
		{
			name:  "Single File Torrent",
			input: "d8:announce35:http://tracker.example.com/announce4:infod6:lengthi1024e4:name8:test.txt12:piece lengthi512e6:pieces20:12345678901234567890ee",
			want: &TorrentFile{
				Announce: "http://tracker.example.com/announce",
				Info: InfoDict{
					Name:        "test.txt",
					Length:      1024,
					PieceLength: 512,
					Pieces:      "12345678901234567890",
				},
			},
		},
		{
			name:  "Multi File Torrent with All Optional Metadata",
			input: "d8:announce21:http://tracker.com/an13:announce-listll21:http://tracker.com/anee7:comment4:cool10:created by2:Go13:creation datei1700000000e4:infod5:filesld6:lengthi500e4:pathl9:file1.txteed6:lengthi300e4:pathl4:docs9:file2.txteee4:name10:rootfolder12:piece lengthi256e6:pieces20:12345678901234567890ee",
			want: &TorrentFile{
				Announce:     "http://tracker.com/an",
				AnnounceList: [][]string{{"http://tracker.com/an"}},
				Comment:      "cool",
				CreatedBy:    "Go",
				CreationDate: 1700000000,
				Info: InfoDict{
					Name:        "rootfolder",
					PieceLength: 256,
					Pieces:      "12345678901234567890",
					Files: []File{
						{Length: 500, Path: []string{"file1.txt"}},
						{Length: 300, Path: []string{"docs", "file2.txt"}},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeTorrent(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:  %+v\nwant: %+v", got, tt.want)
			}
		})
	}
}

func TestDecodeTorrent_InvalidCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Empty Input",
			input: "",
		},
		{
			name:  "Root Is Not Dict",
			input: "i42e",
		},
		{
			name:  "Missing Info Dict",
			input: "d8:announce21:http://tracker.com/ane",
		},
		{
			name:  "Truncated String Payload",
			input: "d8:announce35:http://tracker.example.come",
		},
		{
			name:  "Unclosed Integer",
			input: "d8:announce21:http://tracker.com/an4:infod12:piece lengthi512ee",
		},
		{
			name:  "Non-String Dict Key",
			input: "di42e4:infoe",
		},
		{
			name:  "Unclosed Dict",
			input: "d8:announce21:http://tracker.com/an",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeTorrent(tt.input)
			if err == nil {
				t.Errorf("expected error for input %q, got nil", tt.input)
			}
		})
	}
}

func TestParseNext_Primitives(t *testing.T) {
	t.Run("Negative Integer", func(t *testing.T) {
		val, consumed, err := parseNext("i-42e", 0)
		if err != nil || val.(int64) != -42 || consumed != 5 {
			t.Errorf("got (%v, %d, %v), want (-42, 5, nil)", val, consumed, err)
		}
	})

	t.Run("Zero Integer", func(t *testing.T) {
		val, consumed, err := parseNext("i0e", 0)
		if err != nil || val.(int64) != 0 || consumed != 3 {
			t.Errorf("got (%v, %d, %v), want (0, 3, nil)", val, consumed, err)
		}
	})

	t.Run("Empty String", func(t *testing.T) {
		val, consumed, err := parseNext("0:", 0)
		if err != nil || val.(string) != "" || consumed != 2 {
			t.Errorf("got (%v, %d, %v), want (\"\", 2, nil)", val, consumed, err)
		}
	})
}