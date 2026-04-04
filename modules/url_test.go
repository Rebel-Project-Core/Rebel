package modules

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"credo/cache"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ============================================================
// isArchive
// ============================================================

func Test_isArchive_Zip(t *testing.T) {
	for _, url := range []string{
		"https://example.com/file.zip",
		"https://example.com/File.ZIP",
		"https://example.com/path/to/archive.zip",
	} {
		if !isArchive(url) {
			t.Errorf("isArchive(%q) = false, want true", url)
		}
	}
}

func Test_isArchive_TarGz(t *testing.T) {
	for _, url := range []string{
		"https://example.com/archive.tar.gz",
		"https://example.com/archive.TAR.GZ",
	} {
		if !isArchive(url) {
			t.Errorf("isArchive(%q) = false, want true", url)
		}
	}
}

func Test_isArchive_Tgz(t *testing.T) {
	if !isArchive("https://example.com/archive.tgz") {
		t.Error("isArchive(.tgz) = false, want true")
	}
}

func Test_isArchive_NonArchive(t *testing.T) {
	for _, url := range []string{
		"https://example.com/binary",
		"https://example.com/file.exe",
		"https://example.com/data.json",
		"https://example.com/script.sh",
	} {
		if isArchive(url) {
			t.Errorf("isArchive(%q) = true, want false", url)
		}
	}
}

// ============================================================
// copyExecutable
// ============================================================

func Test_copyExecutable_CopiesFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dst := filepath.Join(dir, "dst")
	os.WriteFile(src, []byte("binary content"), 0644)

	if err := copyExecutable(src, dst); err != nil {
		t.Fatalf("copyExecutable error: %v", err)
	}
	got, _ := os.ReadFile(dst)
	if string(got) != "binary content" {
		t.Errorf("unexpected content: %q", got)
	}
}

func Test_copyExecutable_SetsExecutableBit(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "tool")
	dst := filepath.Join(dir, "tool-out")
	os.WriteFile(src, []byte("#!/bin/sh"), 0644)

	if err := copyExecutable(src, dst); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(dst)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&0100 == 0 {
		t.Errorf("expected executable bit set, got %04o", info.Mode())
	}
}

func Test_copyExecutable_IntoDirectory(t *testing.T) {
	dir := t.TempDir()
	destDir := filepath.Join(dir, "dest")
	os.MkdirAll(destDir, 0755)
	src := filepath.Join(dir, "mytool")
	os.WriteFile(src, []byte("data"), 0644)

	if err := copyExecutable(src, destDir); err != nil {
		t.Fatalf("copyExecutable into dir error: %v", err)
	}
	expected := filepath.Join(destDir, "mytool")
	if _, err := os.Stat(expected); err != nil {
		t.Errorf("file not placed inside directory: %v", err)
	}
}

func Test_copyExecutable_MissingSrc(t *testing.T) {
	dir := t.TempDir()
	if err := copyExecutable(filepath.Join(dir, "nope"), filepath.Join(dir, "out")); err == nil {
		t.Error("expected error for missing source")
	}
}

// ============================================================
// extractZip
// ============================================================

func makeZip(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "test.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	for name, content := range files {
		fw, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		fw.Write([]byte(content))
	}
	w.Close()
	f.Close()
	return zipPath
}

func Test_extractZip_ExtractsFiles(t *testing.T) {
	zipPath := makeZip(t, map[string]string{
		"hello.txt":     "hello world",
		"dir/nested.txt": "nested content",
	})
	dest := t.TempDir()
	if err := extractZip(zipPath, dest); err != nil {
		t.Fatalf("extractZip error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dest, "hello.txt"))
	if err != nil {
		t.Fatalf("extracted file not found: %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("unexpected content: %q", data)
	}
}

func Test_extractZip_CreatesNestedDirs(t *testing.T) {
	zipPath := makeZip(t, map[string]string{
		"a/b/c.txt": "deep",
	})
	dest := t.TempDir()
	if err := extractZip(zipPath, dest); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dest, "a/b/c.txt")); err != nil {
		t.Errorf("nested file not extracted: %v", err)
	}
}

func Test_extractZip_PathTraversal(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "evil.zip")
	f, _ := os.Create(zipPath)
	w := zip.NewWriter(f)
	fw, _ := w.Create("../evil.txt")
	fw.Write([]byte("pwned"))
	w.Close()
	f.Close()

	dest := t.TempDir()
	err := extractZip(zipPath, dest)
	if err == nil {
		t.Error("expected error for path traversal in zip")
	}
	if !strings.Contains(err.Error(), "unsafe") {
		t.Errorf("expected 'unsafe' in error, got: %v", err)
	}
}

func Test_extractZip_InvalidFile(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.zip")
	os.WriteFile(bad, []byte("not a zip"), 0644)
	if err := extractZip(bad, t.TempDir()); err == nil {
		t.Error("expected error for invalid zip")
	}
}

// ============================================================
// extractTarGz
// ============================================================

func makeTarGz(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	tgzPath := filepath.Join(dir, "test.tar.gz")
	f, err := os.Create(tgzPath)
	if err != nil {
		t.Fatal(err)
	}
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	for name, content := range files {
		hdr := &tar.Header{
			Name: name,
			Mode: 0644,
			Size: int64(len(content)),
		}
		tw.WriteHeader(hdr)
		tw.Write([]byte(content))
	}
	tw.Close()
	gz.Close()
	f.Close()
	return tgzPath
}

func makeTarGzWithDir(t *testing.T, dir string, files map[string]string) string {
	t.Helper()
	outDir := t.TempDir()
	tgzPath := filepath.Join(outDir, "test.tar.gz")
	f, err := os.Create(tgzPath)
	if err != nil {
		t.Fatal(err)
	}
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	// Add directory entry.
	tw.WriteHeader(&tar.Header{
		Typeflag: tar.TypeDir,
		Name:     dir + "/",
		Mode:     0755,
	})
	for name, content := range files {
		hdr := &tar.Header{
			Name: dir + "/" + name,
			Mode: 0644,
			Size: int64(len(content)),
		}
		tw.WriteHeader(hdr)
		tw.Write([]byte(content))
	}
	tw.Close()
	gz.Close()
	f.Close()
	return tgzPath
}

func Test_extractTarGz_ExtractsFiles(t *testing.T) {
	tgz := makeTarGz(t, map[string]string{"file.txt": "tar content"})
	dest := t.TempDir()
	if err := extractTarGz(tgz, dest); err != nil {
		t.Fatalf("extractTarGz error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dest, "file.txt"))
	if err != nil {
		t.Fatalf("file not extracted: %v", err)
	}
	if string(data) != "tar content" {
		t.Errorf("unexpected content: %q", data)
	}
}

func Test_extractTarGz_ExtractsDirectory(t *testing.T) {
	tgz := makeTarGzWithDir(t, "mydir", map[string]string{"inner.txt": "inner"})
	dest := t.TempDir()
	if err := extractTarGz(tgz, dest); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dest, "mydir", "inner.txt")); err != nil {
		t.Errorf("directory entry not extracted: %v", err)
	}
}

func Test_extractTarGz_PathTraversal(t *testing.T) {
	dir := t.TempDir()
	tgzPath := filepath.Join(dir, "evil.tar.gz")
	f, _ := os.Create(tgzPath)
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	tw.WriteHeader(&tar.Header{
		Name: "../evil.txt",
		Mode: 0644,
		Size: 5,
	})
	tw.Write([]byte("pwned"))
	tw.Close()
	gz.Close()
	f.Close()

	dest := t.TempDir()
	err := extractTarGz(tgzPath, dest)
	if err == nil {
		t.Error("expected error for path traversal in tar.gz")
	}
	if !strings.Contains(err.Error(), "unsafe") {
		t.Errorf("expected 'unsafe' in error, got: %v", err)
	}
}

func Test_extractTarGz_InvalidFile(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.tar.gz")
	os.WriteFile(bad, []byte("not gzip"), 0644)
	if err := extractTarGz(bad, t.TempDir()); err == nil {
		t.Error("expected error for invalid tar.gz")
	}
}

// ============================================================
// extractArchive dispatch
// ============================================================

func Test_extractArchive_DispatchesZip(t *testing.T) {
	zipPath := makeZip(t, map[string]string{"f.txt": "hi"})
	dest := t.TempDir()
	if err := extractArchive(zipPath, dest); err != nil {
		t.Fatalf("extractArchive(.zip) error: %v", err)
	}
}

func Test_extractArchive_DispatchesTarGz(t *testing.T) {
	tgz := makeTarGz(t, map[string]string{"f.txt": "hi"})
	dest := t.TempDir()
	if err := extractArchive(tgz, dest); err != nil {
		t.Fatalf("extractArchive(.tar.gz) error: %v", err)
	}
}

// ============================================================
// downloadFile
// ============================================================

func Test_downloadFile_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("file contents"))
	}))
	defer ts.Close()

	dest := filepath.Join(t.TempDir(), "out.bin")
	if err := downloadFile(ts.URL, dest); err != nil {
		t.Fatalf("downloadFile error: %v", err)
	}
	data, _ := os.ReadFile(dest)
	if !bytes.Equal(data, []byte("file contents")) {
		t.Errorf("unexpected content: %q", data)
	}
}

func Test_downloadFile_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	if err := downloadFile(ts.URL, filepath.Join(t.TempDir(), "out")); err == nil {
		t.Error("expected error for 404 response")
	}
}

func Test_downloadFile_InvalidURL(t *testing.T) {
	if err := downloadFile("http://127.0.0.1:0/notlistening", t.TempDir()+"/out"); err == nil {
		t.Error("expected error for unreachable URL")
	}
}

// ============================================================
// bareRun (urlModule)
// ============================================================

func Test_urlModule_bareRun_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer ts.Close()

	m := &urlModule{}
	spell := urlSpell{URL: ts.URL + "/file.bin", Path: "/tmp"}
	got, err := m.bareRun(spell)
	if err != nil {
		t.Fatalf("bareRun error: %v", err)
	}
	if got.URL != spell.URL {
		t.Errorf("got URL %q, want %q", got.URL, spell.URL)
	}
}

func Test_urlModule_bareRun_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	m := &urlModule{}
	_, err := m.bareRun(urlSpell{URL: ts.URL + "/secret", Path: "/tmp"})
	if err == nil {
		t.Error("expected error for non-2xx response")
	}
}

func Test_urlModule_bareRun_InvalidURL(t *testing.T) {
	m := &urlModule{}
	_, err := m.bareRun(urlSpell{URL: "http://127.0.0.1:0/nope", Path: "/tmp"})
	if err == nil {
		t.Error("expected error for unreachable URL")
	}
}

// ============================================================
// extractZip with directory entries
// ============================================================

func makeZipWithDir(t *testing.T, dirName string, files map[string]string) string {
	t.Helper()
	d := t.TempDir()
	zipPath := filepath.Join(d, "test.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	// Explicit directory entry.
	w.Create(dirName + "/")
	for name, content := range files {
		fw, _ := w.Create(dirName + "/" + name)
		fw.Write([]byte(content))
	}
	w.Close()
	f.Close()
	return zipPath
}

func Test_extractZip_DirectoryEntry(t *testing.T) {
	zipPath := makeZipWithDir(t, "subdir", map[string]string{"data.txt": "hello"})
	dest := t.TempDir()
	if err := extractZip(zipPath, dest); err != nil {
		t.Fatalf("extractZip with directory entry error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "subdir")); err != nil {
		t.Errorf("directory not created: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "subdir", "data.txt")); err != nil {
		t.Errorf("file inside directory not extracted: %v", err)
	}
}

// ============================================================
// urlModule.bareRun — cache hit
// ============================================================

func Test_urlModule_bareRun_CacheHit(t *testing.T) {
	const cacheURL = "https://example.com/already-cached-file"
	// Prime the cache so bareRun returns without making any HTTP call.
	_ = cache.Insert(urlModuleName, cacheURL, urlSpell{URL: cacheURL, Path: "/cached"})

	m := &urlModule{}
	got, err := m.bareRun(urlSpell{URL: cacheURL, Path: "/cached"})
	if err != nil {
		t.Fatalf("bareRun cache-hit error: %v", err)
	}
	if got.URL != cacheURL {
		t.Errorf("got URL %q, want %q", got.URL, cacheURL)
	}
}

// ============================================================
// urlModule.Save and Apply (integration via temp project dir)
// ============================================================

func Test_urlModule_Save_DownloadsFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("binary data"))
	}))
	defer ts.Close()

	fileURL := ts.URL + "/myfile.bin"
	m := &urlModule{}
	if err := m.Save(urlSpell{URL: fileURL, Path: "/unused"}); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	// Verify the file ended up in the download directory.
	downloadDir, err := moduleDownloadPath(urlModuleName)
	if err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(downloadDir, "myfile.bin")
	if _, err := os.Stat(dest); err != nil {
		t.Errorf("downloaded file not found at %s: %v", dest, err)
	}
}

func Test_urlModule_Save_WrongType(t *testing.T) {
	m := &urlModule{}
	if err := m.Save("not-a-spell"); err != ErrConverting {
		t.Fatalf("want ErrConverting, got %v", err)
	}
}

func Test_urlModule_Save_CacheHit_NoDownload(t *testing.T) {
	const cachedURL = "https://example.com/cached-save"
	_ = cache.Insert(urlModuleName, cachedURL, true)

	// No HTTP server — if Save tries to download, it will fail.
	m := &urlModule{}
	if err := m.Save(urlSpell{URL: cachedURL, Path: "/unused"}); err != nil {
		t.Fatalf("Save with cached URL should skip download, got: %v", err)
	}
}

func Test_urlModule_Apply_CopiesFile(t *testing.T) {
	// Pre-populate the download directory with a file.
	downloadDir, err := moduleDownloadPath(urlModuleName)
	if err != nil {
		t.Fatal(err)
	}
	fileName := "applytest.bin"
	srcPath := filepath.Join(downloadDir, fileName)
	if err := os.WriteFile(srcPath, []byte("exec content"), 0644); err != nil {
		t.Fatal(err)
	}

	dest := t.TempDir()
	spell := urlSpell{
		URL:  "https://example.com/" + fileName,
		Path: dest,
	}
	m := &urlModule{}
	if err := m.Apply(spell); err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, fileName)); err != nil {
		t.Errorf("applied file not found: %v", err)
	}
}

func Test_urlModule_Apply_ExtractsZip(t *testing.T) {
	downloadDir, err := moduleDownloadPath(urlModuleName)
	if err != nil {
		t.Fatal(err)
	}
	zipPath := filepath.Join(downloadDir, "apply.zip")
	// Create a zip in the download directory.
	f, _ := os.Create(zipPath)
	w := zip.NewWriter(f)
	fw, _ := w.Create("contents.txt")
	fw.Write([]byte("zipped"))
	w.Close()
	f.Close()

	dest := t.TempDir()
	spell := urlSpell{
		URL:  "https://example.com/apply.zip",
		Path: dest,
	}
	m := &urlModule{}
	if err := m.Apply(spell); err != nil {
		t.Fatalf("Apply (zip) error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "contents.txt")); err != nil {
		t.Errorf("zip contents not extracted: %v", err)
	}
}

func Test_urlModule_Apply_WrongType(t *testing.T) {
	m := &urlModule{}
	if err := m.Apply(42); err != ErrConverting {
		t.Fatalf("want ErrConverting, got %v", err)
	}
}

func Test_urlModule_BulkSave_WithEntries_CallsSave(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("data"))
	}))
	defer ts.Close()

	cfg := &Config{
		Url: []urlSpell{
			{URL: ts.URL + "/bulk1.bin", Path: "/tmp"},
		},
	}
	m := &urlModule{}
	if err := m.BulkSave(cfg); err != nil {
		t.Fatalf("BulkSave error: %v", err)
	}
}

func Test_urlModule_Save_DownloadFails(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	m := &urlModule{}
	// Use a URL not in cache so downloadFile is actually called.
	err := m.Save(urlSpell{URL: ts.URL + "/fail-unique-9999.bin", Path: "/tmp"})
	if err == nil {
		t.Error("expected error for failed download")
	}
}

func Test_urlModule_BulkApply_WithEntries(t *testing.T) {
	downloadDir, err := moduleDownloadPath(urlModuleName)
	if err != nil {
		t.Fatal(err)
	}
	fileName := "bulkapply.bin"
	os.WriteFile(filepath.Join(downloadDir, fileName), []byte("data"), 0644)

	dest := t.TempDir()
	cfg := &Config{
		Url: []urlSpell{
			{URL: "https://example.com/" + fileName, Path: dest},
		},
	}
	m := &urlModule{}
	if err := m.BulkApply(cfg); err != nil {
		t.Fatalf("BulkApply error: %v", err)
	}
}

// ============================================================
// cobraArgs (urlModule)
// ============================================================

func Test_urlModule_cobraArgs_NoArgs(t *testing.T) {
	m := &urlModule{}
	cmd := m.CliConfig(&Config{})
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with no args")
	}
}

func Test_urlModule_cobraArgs_NoPath(t *testing.T) {
	m := &urlModule{}
	cmd := m.CliConfig(&Config{})
	// Set URL arg but no --path flag.
	err := cmd.Args(cmd, []string{"https://example.com/f"})
	if err == nil {
		t.Error("expected error when --path not set")
	}
}

func Test_urlModule_cobraArgs_Valid(t *testing.T) {
	m := &urlModule{}
	cmd := m.CliConfig(&Config{})
	// Parse flags so --path is visible via cmd.Flags().GetString("path").
	cmd.ParseFlags([]string{"--path", "/tmp/dest"})
	err := cmd.Args(cmd, []string{"https://example.com/f"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
