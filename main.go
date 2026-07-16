// Package main provides the protoreg CLI tool.
//
// protoreg is a code generator that automatically creates Marshal and Unmarshal
// methods for Go structs annotated with protoreg tags. It supports flexible
// encoding configurations through struct field tags, including byte order
// (big-endian/little-endian), word order (high-word-first/low-word-first),
// generation modes (marshal only, unmarshal only, or both), and fixed-size
// arrays of integers, booleans, and floating-point values.
//
// Usage:
//
//	protoreg -type TypeName [options] [package]
//
// Flags:
//
//	-type string
//	    Comma-separated list of struct names to generate unmarshaler/marshaler for (required)
//	-o string
//	    Output file name. Default is <file>_<key>.go
//	-v
//	    Enable verbose (debug) logging
//	-key string
//	    Struct tag key to use (default: "protoreg")
//
// # Struct Tags
//
// Tags are placed on struct fields using the tag key (default: "protoreg").
// Tags use a key=value format separated by commas.
//
// ## Struct-Level Tags
//
// Struct-level configuration is specified on an empty field named "_":
//
//	type MyData struct {
//	  _  struct{} `protoreg:"encoding=big,wordorder=high,mode=all"`
//	  ...
//	}
//
// Supported struct-level tags:
//
//	encoding (default: "big")
//	    Byte order for multi-byte values.
//	    Values: "big" (big-endian), "little" (little-endian)
//
//	wordorder (default: "high")
//	    Order of 16-bit words within multi-word values (32-bit, 64-bit).
//	    Values: "high" (high word first), "low" (low word first)
//
//	mode (default: "all")
//	    Which functions to generate.
//	    Values: "all" (both Marshal/Unmarshal), "marshal", "unmarshal"
//
//	marshalfunc (default: "Marshal")
//	    Custom name for the generated Marshal method.
//	    Example: marshalfunc=Serialize
//
//	unmarshalfunc (default: "Unmarshal")
//	    Custom name for the generated Unmarshal method.
//	    Example: unmarshalfunc=Deserialize
//
// ## Field-Level Tags
//
// Field-level tags specify how individual fields should be marshaled/unmarshaled.
//
//	offset (required)
//	    Zero-based offset in the buffer (in uint16 units).
//	    Example: offset=0, offset=5
//
//	size (for strings)
//	    Number of uint16 elements to use for the field.
//	    Example: size=8 reserves 8 uint16s (16 bytes for char8, 8 bytes for char16)
//
//	Arrays are inferred from the Go array declaration. Fixed-size arrays of
//	integers, booleans, and floating-point values are supported; no extra tag
//	is required beyond the usual offset tag.
//
//	char (for strings, default: "8")
//	    Character width for string fields.
//	    Values: "8" (8-bit characters), "16" (16-bit characters)
//
//	charencoding (for strings, default: "ascii")
//	    Character encoding for string fields.
//	    Values: "ascii", "utf8"
//	    Note: Only applies to char16 mode
//
//	byte (for individual bytes)
//	    Which byte to extract from a uint16.
//	    Values: "high" (upper byte), "low" (lower byte)
//	    Example: byte=low extracts the lower 8 bits
//
// # Examples
//
// ## Basic integer marshaling
//
//	type SimpleData struct {
//	  _ struct{} `protoreg:"encoding=big"`
//	  Field1 uint16 `protoreg:"offset=0"`
//	  Field2 uint32 `protoreg:"offset=1"`
//	}
//
// ## Mixed encoding with word order
//
//	type MixedData struct {
//	  _ struct{} `protoreg:"encoding=little,wordorder=low"`
//	  Timestamp uint64 `protoreg:"offset=0"`
//	  Counter   uint32 `protoreg:"offset=4"`
//	}
//
// ## String handling with custom character encoding
//
//	type StringData struct {
//	  _ struct{} `protoreg:"encoding=big,char=8"`
//	  Name   string `protoreg:"offset=0,size=16,char=8"`
//	  Label  string `protoreg:"offset=16,size=8,char=16,charencoding=utf8"`
//	}
//
// ## Fixed-size array support
//
//	type ArrayData struct {
//	  _ struct{} `protoreg:"encoding=big"`
//	  Values  [4]uint16  `protoreg:"offset=0"`
//	  Counts  [3]uint32  `protoreg:"offset=4"`
//	  Samples [5]float32 `protoreg:"offset=10"`
//	  Flags   [4]bool    `protoreg:"offset=20"`
//	}
//
// ## Custom function names
//
//	type CustomData struct {
//	  _ struct{} `protoreg:"marshalfunc=Encode,unmarshalfunc=Decode"`
//	  Value uint16 `protoreg:"offset=0"`
//	}
//	// Generated methods: Encode() and Decode()
//
// ## Marshal-only generation
//
//	type SendOnly struct {
//	  _ struct{} `protoreg:"mode=marshal"`
//	  Data uint32 `protoreg:"offset=0"`
//	}
//	// Only Marshal() method is generated
//
// When used with go:generate, the tool automatically detects the calling package and file.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/nocccer/protoreg/generator"
)

// main is the entry point for the protoreg CLI tool. It parses command-line flags,
// initializes the code generator, and writes the generated marshal/unmarshal methods
// to the specified output file.
func main() {
	typeFlag := flag.String(
		"type",
		"",
		"List of struct names to generate unmarshaler/marshaler for.",
	)
	outputFlag := flag.String("o", "", "Output file name. Default is <file>_protoreg.go.")
	verbose := flag.Bool("v", false, "Enable verbose logging.")
	keyFlag := flag.String("key", "protoreg", "Struct tag key to use.")
	version := flag.Bool("version", false, "Print version information and exit.")
	flag.Parse()

	if *version {
		// Build-Informationen auslesen
		info, ok := debug.ReadBuildInfo()
		if !ok {
			fmt.Println(
				"No build information available. Please ensure the binary was built with module support.",
			)
			os.Exit(1)
		}

		// Main
		fmt.Printf("Main.Version: %s\n", info.Main.Version)
		fmt.Printf("Main.Path: %s\n", info.Main.Path)
		fmt.Printf("GoVersion: %s\n", info.GoVersion)

		// Settings
		fmt.Println("\nSettings:")
		for _, setting := range info.Settings {
			fmt.Printf(" - %s: %s\n", setting.Key, setting.Value)
		}

		// Dependencies
		fmt.Println("\nDependencies:")
		for _, dep := range info.Deps {
			fmt.Printf(" - %s@%s\n", dep.Path, dep.Version)
		}
		os.Exit(0)
	}

	level := slog.LevelInfo
	if *verbose {
		level = slog.LevelDebug
	}

	log := slog.New(tint.NewTextHandler(os.Stderr, &tint.Options{
		Level: level,
	}))

	if *typeFlag == "" {
		log.Error("'-type' flag is required")
		os.Exit(1)
	}

	log.Debug(
		"flags",
		slog.String("type", *typeFlag),
		slog.String("output", *outputFlag),
		slog.String("key", *keyFlag),
	)

	g := generator.NewGenerator(
		strings.Split(*typeFlag, ","),
		getPkg(),
		log,
		generator.WithTagKey(*keyFlag),
	)

	if err := g.Generate(); err != nil {
		log.Error(
			"failed to generate marshaler",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	if *outputFlag == "" {
		file := os.Getenv("GOFILE")
		*outputFlag = fmt.Sprintf("%s_%s.go", strings.Split(filepath.Base(file), ".")[0], *keyFlag)
	}

	if err := g.WriteToFile(*outputFlag); err != nil {
		log.Error(
			"failed to write marshaler code to file",
			slog.String("file", *outputFlag),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}

// getPkg determines the package directory to scan for struct definitions.
// If no arguments are provided, it defaults to the current directory.
// If a single directory argument is provided, that directory is returned.
// Otherwise, the directory containing the first argument file is returned.
func getPkg() string {
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	if len(args) == 1 && isDirectory(args[0]) {
		return args[0]
	}

	return filepath.Dir(args[0])
}

// isDirectory reports whether the given path is a valid directory.
// It returns false if the path does not exist or is not a directory.
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
