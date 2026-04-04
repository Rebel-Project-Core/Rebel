package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// resetGPath resets the cached project path so the next call computes it fresh.
func resetGPath() {
	gPath = nil
}

func Test_ProjectPath_ReturnsNonEmpty(t *testing.T) {
	resetGPath()
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer func() { os.Chdir(orig); resetGPath() }()
	os.Chdir(dir)

	p, err := ProjectPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == "" {
		t.Fatal("ProjectPath() returned empty string")
	}
}

func Test_ProjectPath_EndsWithProjectDirectoryName(t *testing.T) {
	resetGPath()
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer func() { os.Chdir(orig); resetGPath() }()
	os.Chdir(dir)

	p, err := ProjectPath()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(p, ProjectDirectoryName) {
		t.Errorf("expected path ending with %q, got %q", ProjectDirectoryName, p)
	}
}

func Test_ProjectPath_CreatesDirectory(t *testing.T) {
	resetGPath()
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer func() { os.Chdir(orig); resetGPath() }()
	os.Chdir(dir)

	p, err := ProjectPath()
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(p)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Errorf("expected directory, got file at %s", p)
	}
}

func Test_ProjectPath_IsInsideWorkingDir(t *testing.T) {
	resetGPath()
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer func() { os.Chdir(orig); resetGPath() }()
	os.Chdir(dir)

	p, err := ProjectPath()
	if err != nil {
		t.Fatal(err)
	}
	// Resolve both through filepath.EvalSymlinks to handle macOS /private/ symlinks.
	realDir, _ := filepath.EvalSymlinks(dir)
	realP, _ := filepath.EvalSymlinks(p)
	if !strings.HasPrefix(realP, realDir) {
		t.Errorf("project path %q should be inside working dir %q", realP, realDir)
	}
}

func Test_ProjectPath_CachesResult(t *testing.T) {
	resetGPath()
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer func() { os.Chdir(orig); resetGPath() }()
	os.Chdir(dir)

	first, err := ProjectPath()
	if err != nil {
		t.Fatal(err)
	}
	// Change directory — second call should still return the cached (first) value.
	os.Chdir(t.TempDir())
	second, err := ProjectPath()
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Errorf("expected cached result; got %q then %q", first, second)
	}
}

func Test_ProjectPath_ReturnsCopy(t *testing.T) {
	resetGPath()
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer func() { os.Chdir(orig); resetGPath() }()
	os.Chdir(dir)

	p1, _ := ProjectPath()
	p2, _ := ProjectPath()
	// Modifying one should not affect the other (they are value copies, not pointers).
	if &p1 == &p2 {
		t.Error("ProjectPath() should return value copies, not pointers to the same variable")
	}
}
