package version_test

import (
	"bytes"
	"credo/version"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(fn func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func Test_PrintVersion_NoError(t *testing.T) {
	if err := version.PrintVersion("testapp"); err != nil {
		t.Fatalf("PrintVersion returned unexpected error: %v", err)
	}
}

func Test_PrintVersion_ContainsName(t *testing.T) {
	out := captureStdout(func() {
		_ = version.PrintVersion("mycredo")
	})
	if !strings.Contains(out, "mycredo") {
		t.Errorf("output does not contain app name: %q", out)
	}
}

func Test_PrintVersion_ContainsVersionFields(t *testing.T) {
	out := captureStdout(func() {
		_ = version.PrintVersion("app")
	})
	for _, field := range []string{"Version", "Commit", "Build Date"} {
		if !strings.Contains(out, field) {
			t.Errorf("output missing field %q: %q", field, out)
		}
	}
}

func Test_PrintVersion_DefaultValues(t *testing.T) {
	out := captureStdout(func() {
		_ = version.PrintVersion("app")
	})
	// Default values set in version.go
	if !strings.Contains(out, version.Version) {
		t.Errorf("output missing version %q: %q", version.Version, out)
	}
}
