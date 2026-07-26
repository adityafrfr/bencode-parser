# Planned additions

This roadmap collects useful ways to interact with the parser as it grows from a focused torrent decoder into a fuller bencode library. Pick an item that interests you, open an issue to discuss the shape of the API, and contribute a focused pull request with tests.

## Generic bencode APIs

- [ ] Export a `Decode` function that returns a generic bencode value (`int64`, `string`, `[]any`, or `map[string]any`).
- [ ] Add `DecodeBytes([]byte)` so callers can work directly with binary data.
- [ ] Add a `DecodeReader(io.Reader)` API for decoding from files, network connections, and streams.
- [ ] Define public value types or helpers for inspecting nested dictionaries and lists ergonomically.
- [ ] Add path/query helpers, such as retrieving `info.name` from a decoded dictionary.

## Encoding and Go integration

- [ ] Add `Encode` for converting generic Go values back to bencode.
- [ ] Add `Marshal` and `Unmarshal` APIs inspired by Go's `encoding/json` package.
- [ ] Support struct tags for mapping Go structs to bencode dictionary keys.
- [ ] Add custom marshal/unmarshal interfaces for application-defined types.
- [ ] Provide examples for configuration files, distributed-hash-table messages, and other non-torrent bencode data.

## Torrent tooling

- [ ] Add helpers to calculate an info hash from the exact bencoded `info` dictionary.
- [ ] Generate magnet links from decoded torrent metadata.
- [ ] Add torrent validation helpers for piece length, piece hashes, file entries, and tracker tiers.
- [ ] Add helpers for displaying a torrent's total size and file tree.
- [ ] Add APIs for creating `.torrent` metadata from files and directories.

## Developer experience

- [ ] Publish package documentation on pkg.go.dev with runnable examples.
- [ ] Add fuzz tests for the decoder.
- [ ] Add benchmark coverage for large strings, deeply nested values, and large torrent files.
- [ ] Add a command-line tool for inspecting bencode and `.torrent` files.
- [ ] Add GitHub Actions for formatting, tests, race detection, and static analysis.

## Contribution ideas

Good first contributions include a small API with documentation, a focused test suite, an example program, benchmark coverage, or a new command-line subcommand. Keep pull requests scoped, explain the intended public API, and include tests for both normal and malformed input.
