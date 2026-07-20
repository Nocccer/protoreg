package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestIsCacheSameUsesCache(t *testing.T) {
	t.Parallel()

	outputPath := t.TempDir() + "/generated.go"
	cachePath := outputPath + ".cache"

	if err := os.WriteFile(outputPath, []byte("generated"), 0o600); err != nil {
		t.Fatalf("write output: %v", err)
	}
	if err := os.WriteFile(cachePath, []byte("same-key"), 0o600); err != nil {
		t.Fatalf("write cache: %v", err)
	}

	skipped, err := isCacheSame(outputPath, cachePath, "same-key")
	if err != nil {
		t.Fatalf("isCacheSame returned error: %v", err)
	}
	if !skipped {
		t.Fatalf("expected generation to be skipped when cache matches")
	}
}

func TestIsCacheSameWithoutMatchingCache(t *testing.T) {
	t.Parallel()

	outputPath := t.TempDir() + "/generated.go"
	cachePath := outputPath + ".cache"

	if err := os.WriteFile(outputPath, []byte("generated"), 0o600); err != nil {
		t.Fatalf("write output: %v", err)
	}
	if err := os.WriteFile(cachePath, []byte("old-key"), 0o600); err != nil {
		t.Fatalf("write cache: %v", err)
	}

	skipped, err := isCacheSame(outputPath, cachePath, "new-key")
	if err != nil {
		t.Fatalf("isCacheSame returned error: %v", err)
	}
	if skipped {
		t.Fatalf("expected generation to run when cache differs")
	}
}

func TestCleanCacheRemovesProtoregCacheDirectory(t *testing.T) {
	t.Parallel()

	gocache := t.TempDir()
	cacheDir := filepath.Join(gocache, "protoreg")
	cacheFile := filepath.Join(cacheDir, "generated.go.cache")

	if err := os.MkdirAll(cacheDir, 0o750); err != nil {
		t.Fatalf("create cache dir: %v", err)
	}
	if err := os.WriteFile(cacheFile, []byte("cached"), 0o600); err != nil {
		t.Fatalf("write cache file: %v", err)
	}

	if err := cleanCache(gocache); err != nil {
		t.Fatalf("cleanCache returned error: %v", err)
	}

	if _, err := os.Stat(cacheDir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected cache dir to be removed, got err=%v", err)
	}
}
