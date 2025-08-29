package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nocccer/protoreg/generator"
	"github.com/lmittmann/tint"
)

func main() {
	typeFlag := flag.String(
		"type",
		"",
		"List of struct names to generate unmarshaler/marshaler for.",
	)
	outputFlag := flag.String("o", "", "Output file name. Default is <file>_protoreg.go.")
	verbose := flag.Bool("v", false, "Enable verbose logging.")
	flag.Parse()

	level := slog.LevelInfo
	if *verbose {
		level = slog.LevelDebug
	}

	log := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
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
	)

	g := generator.NewGenerator(
		strings.Split(*typeFlag, ","),
		getPkg(),
		log,
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
		*outputFlag = fmt.Sprintf("%s_protoreg.go", strings.Split(filepath.Base(file), ".")[0])
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

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
