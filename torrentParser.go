package bencode

import (
	"fmt"
)

type File struct {
	Length int64
	Path   []string
}

type InfoDict struct {
	Name        string
	PieceLength int64
	Pieces      string
	Length      int64
	Files       []File
}

type TorrentFile struct {
	Announce     string
	AnnounceList [][]string
	Comment      string
	CreatedBy    string
	CreationDate int64
	Info         InfoDict
}

func FillTorrentInfo(m map[string]any, tf *TorrentFile) {
	if v, ok := m["announce"].(string); ok {
		tf.Announce = v
	}
	if v, ok := m["comment"].(string); ok {
		tf.Comment = v
	}
	if v, ok := m["created by"].(string); ok {
		tf.CreatedBy = v
	}
	if v, ok := m["creation date"].(int64); ok {
		tf.CreationDate = v
	}
}

func FillRawInfo(rawInfo map[string]any, tf *TorrentFile) {
	info := InfoDict{}

	if v, ok := rawInfo["name"].(string); ok {
		info.Name = v
	}
	if v, ok := rawInfo["piece length"].(int64); ok {
		info.PieceLength = v
	}
	if v, ok := rawInfo["pieces"].(string); ok {
		info.Pieces = v
	}
	if v, ok := rawInfo["length"].(int64); ok {
		info.Length = v
	}

	tf.Info = info
}

func FillAnnounceList(m map[string]any, tf *TorrentFile) {
	if rawList, ok := m["announce-list"].([]any); ok {
		for _, subList := range rawList {
			if stringList, ok := subList.([]any); ok {
				var tier []string
				for _, s := range stringList {
					if str, ok := s.(string); ok {
						tier = append(tier, str)
					}
				}
				tf.AnnounceList = append(tf.AnnounceList, tier)
			}
		}
	}
}

func FillFiles(rawInfo map[string]any, tf *TorrentFile) {
	if rawFiles, ok := rawInfo["files"].([]any); ok {
		for _, rf := range rawFiles {
			if fileMap, ok := rf.(map[string]any); ok {
				f := File{}
				if v, ok := fileMap["length"].(int64); ok {
					f.Length = v
				}
				if pathList, ok := fileMap["path"].([]any); ok {
					for _, p := range pathList {
						if str, ok := p.(string); ok {
							f.Path = append(f.Path, str)
						}
					}
				}
				tf.Info.Files = append(tf.Info.Files, f)
			}
		}
	}
}

func MapToTorrent(m map[string]any) (*TorrentFile, error) {
	tf := &TorrentFile{}

	FillTorrentInfo(m, tf)

	rawInfo, ok := m["info"].(map[string]any)

	if !ok {
		return nil, fmt.Errorf("missing 'info' dictionary")
	}

	FillRawInfo(rawInfo, tf)
	FillAnnounceList(m, tf)
	FillFiles(rawInfo, tf)

	return tf, nil
}

func DecodeTorrent(input string) (*TorrentFile, error) {
	val, _, err := parseNext(input, 0)
	if err != nil {
		return nil, err
	}

	dict, ok := val.(map[string]any)

	if !ok {
		return nil, fmt.Errorf(" root bencode element must be a dictionary")
	}

	return MapToTorrent(dict)
}

/*


d
  8:announce   <tracker URL string>
  13:announce-list  l l <url> <url> e l <url> e e
  7:comment    <string>
  10:created by <string>
  13:creation date <integer>
  4:info d
    4:name         <string>
    12:piece length <integer>
    6:pieces        <string of concatenated SHA1 hashes>
    6:length        <integer>        (only if single file)
    5:files  l d 6:length <int> 4:path l <str> <str> e e ... e   (only if multiple files)
  e
e


*/
