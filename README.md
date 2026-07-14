# protoreg

`protoreg` is a small code generator that creates `Marshal` and `Unmarshal`
methods for Go structs annotated with `protoreg` tags. It targets register
based protocols (slices of `uint16` or `uint32`) commonly used in IoT and
industrial protocols such as Modbus.

This repository contains the generator implementation and tests. See
[main.go](main.go) for the CLI usage and available options.

**Highlights**

- Generate `Marshal()` / `Unmarshal()` for plain and custom-typed fields.
- Control byte order and word order per-struct via tags.
- Support for fixed-size strings, numeric types, single-bit boolean fields, and fixed-size arrays of integers, floats, and booleans.

**Status**

Basic generation and a comprehensive test-suite are included. See the
`tests/` folder for generated example code and unit tests.

## Installation

```bash
go install github.com/nocccer/protoreg@latest
```

## Usage

Use the CLI to generate code for a given type in a package:

```bash
protoreg -type=MyStruct
```

Common flags (see `cmd/main.go` for full docs):

- `-type` (required): comma-separated list of struct names to generate for.
- `-o`: output file name (default: `<file>_protoreg.go`).
- `-v`: verbose logging.
- `-key`: struct tag key to use (default: `protoreg`).

When used with `go:generate` the tool automatically detects the calling
package and file. Example `go:generate` usage is present in the `tests/`
package.

## Struct tags

Configuration is provided using struct tags with the selected tag key
(`protoreg` by default). Tags use `key=value` pairs separated by commas.

Examples and the complete tag reference are documented in
[main.go](main.go). Important tags include:

- `offset` (required): zero-based buffer offset in `uint16` units.
- `size` (strings): number of `uint16` elements reserved for the field.
- `encoding`: `big` (default) or `little` — byte order for multi-byte values.
- `wordorder`: `high` (default) or `low` — 16-bit word order for multi-word values.
- `byte`: `high` or `low` — which byte to read from a `uint16` (for 8-bit fields).
- `bit`: single-bit position (0-15) for boolean fields.
- `mode`: `all` (default), `marshal`, or `unmarshal` to control generation.

## Examples

Basic integer field:

```go
type Simple struct {
  _ struct{} `protoreg:"encoding=big"`
  Value uint16 `protoreg:"offset=0"`
}
```

Boolean packed into a single bit:

```go
type Flags struct {
  _ struct{} `protoreg:"encoding=big"`
  Enabled bool `protoreg:"offset=10,bit=3"`
}
```

The generator emits a fixed hex mask for known `bit` positions (e.g.
`bit=3` -> `0x0008`) so generated `Unmarshal` code uses constants instead of
computing shifts at runtime.

Fixed-size arrays of integers, floats, and booleans:

```go
type Packet struct {
  _ struct{} `protoreg:"encoding=big"`
  Values  [4]uint16  `protoreg:"offset=0"`
  Counts  [3]uint32  `protoreg:"offset=4"`
  Samples [5]float32 `protoreg:"offset=10"`
  Flags   [4]bool    `protoreg:"offset=20"`
}
```

Arrays use the Go array length to determine how many elements are encoded and
rely on the element type for the underlying width.

## Running tests

Run the provided unit tests:

```bash
make test
```

## Contributing

Contributions welcome. Please open issues or PRs. Follow the existing test
patterns in `tests/` when adding new behavior.

## License

See the `LICENSE` file at the repository root for licensing terms.
