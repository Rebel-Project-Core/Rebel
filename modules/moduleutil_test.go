package modules

import (
	"credo/cache"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

// ---- minArgsValidator ----

func Test_minArgsValidator_NoArgs(t *testing.T) {
	v := minArgsValidator("testmod")
	if err := v(&cobra.Command{}, []string{}); err == nil {
		t.Error("expected error with no args")
	}
}

func Test_minArgsValidator_OneArg(t *testing.T) {
	v := minArgsValidator("testmod")
	if err := v(&cobra.Command{}, []string{"arg1"}); err != nil {
		t.Errorf("unexpected error with one arg: %v", err)
	}
}

func Test_minArgsValidator_MultipleArgs(t *testing.T) {
	v := minArgsValidator("testmod")
	if err := v(&cobra.Command{}, []string{"a", "b", "c"}); err != nil {
		t.Errorf("unexpected error with multiple args: %v", err)
	}
}

func Test_minArgsValidator_ErrorMessageContainsModuleName(t *testing.T) {
	v := minArgsValidator("mypkg")
	err := v(&cobra.Command{}, []string{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !containsString(err.Error(), "mypkg") {
		t.Errorf("error message should contain module name, got: %s", err)
	}
}

// ---- externalDependenciesEqual ----

func Test_externalDependenciesEqual_EmptyConfigs(t *testing.T) {
	if !externalDependenciesEqual(Config{}, Config{}) {
		t.Error("two empty configs should be equal")
	}
}

func Test_externalDependenciesEqual_ReflectEquality(t *testing.T) {
	a := Config{}
	b := Config{}
	if !reflect.DeepEqual(a, b) {
		t.Skip("DeepEqual sanity check failed")
	}
	if !externalDependenciesEqual(a, b) {
		t.Error("equal configs should return true")
	}
}

// ---- retrieveFromCache ----

func Test_retrieveFromCache_Miss(t *testing.T) {
	v, ok := retrieveFromCache[urlSpell]("util-test-miss", "nokey")
	if ok || v != nil {
		t.Error("expected cache miss, got hit")
	}
}

func Test_retrieveFromCache_Hit(t *testing.T) {
	const mod = "util-test-hit"
	const key = "mykey"
	want := urlSpell{URL: "https://example.com", Path: "/tmp"}
	_ = cache.Insert(mod, key, want)

	got, ok := retrieveFromCache[urlSpell](mod, key)
	if !ok || got == nil {
		t.Fatal("expected cache hit")
	}
	if got.URL != want.URL || got.Path != want.Path {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func Test_retrieveFromCache_WrongType(t *testing.T) {
	const mod = "util-test-wrongtype"
	const key = "wrongkey"
	// Insert an int, then try to retrieve as urlSpell.
	_ = cache.Insert(mod, key, 42)
	v, ok := retrieveFromCache[urlSpell](mod, key)
	if ok || v != nil {
		t.Error("expected miss for wrong type, got hit")
	}
}

// ---- moduleDownloadPath ----

func Test_moduleDownloadPath_CreatesDirectory(t *testing.T) {
	dir, err := moduleDownloadPath("test-download-path")
	if err != nil {
		t.Fatalf("moduleDownloadPath error: %v", err)
	}
	if dir == "" {
		t.Fatal("expected non-empty path")
	}
	// Second call should succeed (directory already exists).
	dir2, err := moduleDownloadPath("test-download-path")
	if err != nil {
		t.Fatalf("second call error: %v", err)
	}
	if dir != dir2 {
		t.Errorf("expected same path on second call; got %q then %q", dir, dir2)
	}
}

// helper
func containsString(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
