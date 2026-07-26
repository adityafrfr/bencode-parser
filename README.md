# Bencode Parser

A Go package built around a general-purpose [bencode](https://en.wikipedia.org/wiki/Bencode) parser, with a convenient typed interface for BitTorrent metainfo files.

This is a learning project: the code is deliberately compact and focuses on making the bencode parsing flow easy to follow.

## Features

- Decodes bencoded integers, byte strings, lists, and dictionaries.
- Supports nested lists and dictionaries.
- Rejects malformed values, incomplete structures, and bytes after the root bencode value.
- Provides a reusable foundation for working with any bencoded data structure.
- Maps common `.torrent` fields into typed Go structs.
- Supports both single-file and multi-file torrents.

## Current status

The project has a tested bencode parsing core and a complete `DecodeTorrent` flow for common single-file and multi-file torrent metadata. The parser understands bencode's four value types—integers, byte strings, lists, and dictionaries—so the same foundation is suitable for generic bencode tools as well as torrent-specific applications.

The next public APIs and integrations are collected in [PLANNED_ADDITIONS.md](PLANNED_ADDITIONS.md). Contributions are welcome.

## Installation

```bash
go get github.com/adityafrfr/bencode-parser
```

Or add the package to an existing Go module:

```go
import bencode "github.com/adityafrfr/bencode-parser"
```

## Torrent usage

`DecodeTorrent` accepts the raw bencoded contents of a torrent file and returns a `*TorrentFile`.

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

	torrent, err := bencode.DecodeTorrent(string(data))
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

## Data model

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

## Supported torrent fields

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

## Design notes

- Bencoded byte strings, including binary `pieces` data, are represented as Go `string` values. A Go string can contain arbitrary bytes.
- `DecodeTorrent` is the current public convenience API. It builds on the generic bencode parsing core.
- The roadmap includes public generic-decoding, encoding, streaming, validation, and torrent tooling APIs.
