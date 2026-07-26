# Bencode Parser

A Go project for decoding and eventually encoding [bencode](https://en.wikipedia.org/wiki/Bencode), with typed parsing for BitTorrent `.torrent` metainfo files.

The immediate goal is a small, dependable general-purpose bencode package that is useful on its own. The torrent API is built above that foundation, so callers can choose the depth they need: raw bencode values when they want full control, or typed torrent metadata when they want convenience.

## Current Status

The project currently provides:

- Decoding of bencoded integers, byte strings, lists, and dictionaries.
- Support for nested bencode values.
- Rejection of malformed input, incomplete structures, and trailing bytes after the root value.
- `Decoder` and `Unmarshal` APIs for generic bencode decoding.
- `ParseTorrentFile` for converting a bencoded torrent metainfo document into typed Go structs.
- Common single-file and multi-file torrent fields, including announce tiers, files, piece length, and pieces.

The implementation is still evolving. The public API is being shaped around clear package boundaries, so consumers should expect additions and occasional renames before a stable v1 release.

## What This Aims To Be

The long-term project has four layers, each useful independently:

```text
application backend
        |
BitTorrent client
        |
torrent metainfo API
        |
generic bencode codec
```

The generic codec will remain independent of torrents. The torrent layer will interpret
`.torrent` metadata and validate it. A future client layer will use that validated
metadata to communicate with trackers and peers. A backend, if added, will be an
application built on those reusable libraries rather than part of the codec itself.

## Future Goals

- Add bencode encoding through `Encoder` and `Marshal` APIs.
- Establish a stable generic value model, including a deliberate binary byte-string API.
- Split generic bencode and torrent metainfo concerns into separate Go packages.
- Improve torrent metainfo validation and support canonical info-hash calculation.
- Add richer torrent metadata support as the typed API matures.
- Build a BitTorrent client incrementally: storage, peer-wire messages, peers, pieces, and trackers.
- Build an optional backend above the client for managing torrent jobs and exposing an application API.

## Installation

```bash
go get github.com/adityafrfr/bencode-parser
```

Or add the package to an existing Go module:

```go
import bencode "github.com/adityafrfr/bencode-parser"
```

## Generic Bencode Usage

`Unmarshal` decodes one complete bencoded value into Go values:

- bencode integers become `int64`
- bencode byte strings become `string`
- bencode lists become `[]any`
- bencode dictionaries become `map[string]any`

```go
value, err := bencode.Unmarshal([]byte("d3:fooi42ee"))
if err != nil {
	panic(err)
}

dict := value.(map[string]any)
fmt.Println(dict["foo"]) // 42
```

## Torrent Usage

`ParseTorrentFile` accepts the raw bencoded contents of a torrent file and returns a `*TorrentFile`.

```go
package main

import (
	"fmt"
	"os"

	bencode "github.com/adityafrfr/bencode-parser"
)

func main() {
	data, err := os.ReadFile("example.torrent")
	if err != nil {
		panic(err)
	}

	torrent, err := bencode.ParseTorrentFile(string(data))
	if err != nil {
		panic(err)
	}

	fmt.Println("Name:", torrent.Info.Name)
	fmt.Println("Announce URL:", torrent.Announce)
	fmt.Println("Piece length:", torrent.Info.PieceLength)
}
```

For a multi-file torrent, file details are available through `torrent.Info.Files`:

```go
for _, file := range torrent.Info.Files {
	fmt.Printf("%d bytes: %v\n", file.Length, file.Path)
}
```

## Data Model

```go
type TorrentFile struct {
	Announce     string
	AnnounceList [][]string
	Comment      string
	CreatedBy    string
	CreationDate int64
	Info         InfoDict
}

type InfoDict struct {
	Name        string
	PieceLength int64
	Pieces      string
	Length      int64
	Files       []File
}

type File struct {
	Length int64
	Path   []string
}
```

`Info.Length` is used by single-file torrents. Multi-file torrents instead populate `Info.Files`.

## Supported Torrent Fields

| Location | Fields |
| --- | --- |
| Root dictionary | `announce`, `announce-list`, `comment`, `created by`, `creation date`, `info` |
| `info` dictionary | `name`, `piece length`, `pieces`, `length`, `files` |
| Each `files` entry | `length`, `path` |

Unknown fields are safely ignored when creating the typed `TorrentFile` value.

## Testing

Run the test suite:

```bash
go test ./...
```

Optionally run tests with the race detector:

```bash
go test -race ./...
```

## Design Notes

- Bencoded byte strings, including binary `pieces` data, are represented as Go `string` values. A Go string can contain arbitrary bytes.
- `ParseTorrentFile` is the current public torrent convenience API. It builds on the generic bencode parsing core.
- The current parser reads the full input before decoding; streaming decode is a future consideration, not a present guarantee.
- Generic bencode parsing and torrent processing are intentionally separate concerns, even while they currently live in the same Go package.
