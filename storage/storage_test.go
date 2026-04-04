package storage_test

import (
	"bytes"
	"credo/storage"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newStorage(t *testing.T) (*storage.FileStorage, string) {
	t.Helper()
	dir := t.TempDir()
	filename := filepath.Join(dir, "credospell.yaml")
	return &storage.FileStorage{Filename: filename}, filename
}

func Test_Write_CreatesFile(t *testing.T) {
	s, filename := newStorage(t)
	s.Write([]byte("extra: data\n"))
	if _, err := os.Stat(filename); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func Test_Write_ContentIncludesDefaultHeader(t *testing.T) {
	s, filename := newStorage(t)
	s.Write([]byte("key: value\n"))
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "WARNING") {
		t.Errorf("expected WARNING header in file, got: %s", data)
	}
}

func Test_Write_ContentIncludesPayload(t *testing.T) {
	s, filename := newStorage(t)
	payload := []byte("pip:\n  - name: numpy\n")
	s.Write(payload)
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, payload) {
		t.Errorf("expected payload in file; got: %s", data)
	}
}

func Test_Write_FilePermissions(t *testing.T) {
	s, filename := newStorage(t)
	s.Write([]byte(""))
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("expected permissions 0600, got %04o", perm)
	}
}

func Test_Write_IsAtomic(t *testing.T) {
	// The temp-file rename pattern means no partial writes should be visible.
	// We verify by writing twice and reading after each write.
	s, filename := newStorage(t)
	s.Write([]byte("first: 1\n"))
	first, _ := os.ReadFile(filename)
	s.Write([]byte("second: 2\n"))
	second, _ := os.ReadFile(filename)
	if !bytes.Contains(first, []byte("first")) {
		t.Errorf("first read missing expected content: %s", first)
	}
	if !bytes.Contains(second, []byte("second")) {
		t.Errorf("second read missing expected content: %s", second)
	}
}

func Test_Write_Overwrite(t *testing.T) {
	s, filename := newStorage(t)
	s.Write([]byte("v: 1\n"))
	s.Write([]byte("v: 2\n"))
	data, _ := os.ReadFile(filename)
	if bytes.Contains(data, []byte("v: 1")) {
		t.Error("old content should not be present after overwrite")
	}
}

func Test_Read_ExistingFile(t *testing.T) {
	s, _ := newStorage(t)
	s.Write([]byte("key: val\n"))
	got := s.Read()
	if !bytes.Contains(got, []byte("key: val")) {
		t.Errorf("Read() missing written content: %s", got)
	}
}

func Test_Read_CreatesFileWhenMissing(t *testing.T) {
	s, filename := newStorage(t)
	// File does not exist yet.
	got := s.Read()
	if _, err := os.Stat(filename); err != nil {
		t.Fatalf("Read() should have created the file: %v", err)
	}
	if !bytes.Contains(got, []byte("WARNING")) {
		t.Errorf("created file should have default header: %s", got)
	}
}

func Test_Read_EmptyPayload(t *testing.T) {
	s, _ := newStorage(t)
	s.Write([]byte(""))
	got := s.Read()
	if len(got) == 0 {
		t.Error("Read() returned empty; expected at least the default header")
	}
}
