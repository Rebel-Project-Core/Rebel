package config

import (
	"credo/modules"
	"credo/storage"
	"os"
	"path/filepath"
	"testing"
)

// setStorageFile injects a temp-file-backed FileStorage into the singleton,
// returning a cleanup function that restores the original.
func setStorageFile(t *testing.T) (filename string, cleanup func()) {
	t.Helper()
	dir := t.TempDir()
	filename = filepath.Join(dir, "credospell.yaml")
	prev := _internalFileStorage
	_internalFileStorage = &storage.FileStorage{Filename: filename}
	return filename, func() { _internalFileStorage = prev }
}

// --- fromFile ---

func Test_fromFile_ValidYAML(t *testing.T) {
	yaml := []byte("git:\n  - url: https://example.com/repo\n    version: v1\n")
	cfg, err := fromFile(yaml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Git) != 1 {
		t.Fatalf("want 1 git entry, got %d", len(cfg.Git))
	}
}

func Test_fromFile_InvalidYAML(t *testing.T) {
	_, err := fromFile([]byte("{invalid: [yaml"))
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func Test_fromFile_EmptyInput(t *testing.T) {
	cfg, err := fromFile([]byte(""))
	if err != nil {
		t.Fatalf("unexpected error for empty input: %v", err)
	}
	if len(cfg.Git)+len(cfg.Pip)+len(cfg.Apt)+len(cfg.Cran)+len(cfg.Conda)+len(cfg.Url) != 0 {
		t.Error("expected empty config for empty input")
	}
}

func Test_fromFile_DefaultContent(t *testing.T) {
	// The default header written by FileStorage is valid YAML (it starts with "---").
	header := []byte("# comment\n---\n")
	cfg, err := fromFile(header)
	if err != nil {
		t.Fatalf("unexpected error for header-only content: %v", err)
	}
	_ = cfg
}

// --- FileProvider.Get ---

func Test_FileProvider_Get_ReturnsEmptyConfigOnNewFile(t *testing.T) {
	_, cleanup := setStorageFile(t)
	defer cleanup()

	fp := &FileProvider{}
	cfg, err := fp.Get()
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if cfg == nil {
		t.Fatal("Get() returned nil config")
	}
}

func Test_FileProvider_Get_ParsesExistingFile(t *testing.T) {
	filename, cleanup := setStorageFile(t)
	defer cleanup()

	// Write a minimal YAML file with the required header.
	content := "# !!! WARNING !!!\n# generated\n---\npip:\n  - name: requests\n"
	os.WriteFile(filename, []byte(content), 0600)

	fp := &FileProvider{}
	cfg, err := fp.Get()
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if len(cfg.Pip) != 1 {
		t.Fatalf("want 1 pip entry, got %d", len(cfg.Pip))
	}
}

// --- FileProvider.Write ---

func Test_FileProvider_Write_NoError(t *testing.T) {
	_, cleanup := setStorageFile(t)
	defer cleanup()

	fp := &FileProvider{}
	cfg := &modules.Config{}
	if err := fp.Write(cfg); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}
}

func Test_FileProvider_Write_FileCreated(t *testing.T) {
	filename, cleanup := setStorageFile(t)
	defer cleanup()

	fp := &FileProvider{}
	_ = fp.Write(&modules.Config{})
	if _, err := os.Stat(filename); err != nil {
		t.Fatalf("Write() did not create file: %v", err)
	}
}

// --- Round-trip: Write then Get ---

func Test_FileProvider_RoundTrip(t *testing.T) {
	_, cleanup := setStorageFile(t)
	defer cleanup()

	fp := &FileProvider{}
	want := &modules.Config{}
	// We can't easily construct pip/git/cran spells from outside the modules package,
	// but we can verify that an empty config survives a round-trip without error.
	if err := fp.Write(want); err != nil {
		t.Fatal(err)
	}
	got, err := fp.Get()
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("round-trip returned nil")
	}
}

// --- getFileStorage singleton ---

func Test_getFileStorage_Singleton(t *testing.T) {
	prev := _internalFileStorage
	_internalFileStorage = nil
	defer func() { _internalFileStorage = prev }()

	a := getFileStorage()
	b := getFileStorage()
	if a != b {
		t.Error("getFileStorage() should return the same singleton")
	}
}

func Test_getFileStorage_DefaultFilename(t *testing.T) {
	prev := _internalFileStorage
	_internalFileStorage = nil
	defer func() { _internalFileStorage = prev }()

	s := getFileStorage()
	if s.Filename != "credospell.yaml" {
		t.Errorf("expected credospell.yaml, got %q", s.Filename)
	}
}
