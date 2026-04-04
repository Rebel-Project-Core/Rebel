package cache

import (
	"fmt"
	"sync"
	"testing"
)

// resetCache clears the global cache between tests.
func resetCache() {
	mu.Lock()
	defer mu.Unlock()
	cache = make(map[string]map[string]any)
}

func Test_Insert_NewEntry(t *testing.T) {
	resetCache()
	if err := Insert("mod", "key", "value"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_Insert_Duplicate(t *testing.T) {
	resetCache()
	_ = Insert("mod", "dup", 1)
	if err := Insert("mod", "dup", 2); err != ErrAlreadyCached {
		t.Fatalf("want ErrAlreadyCached, got %v", err)
	}
}

func Test_Insert_SameModuleDifferentKeys(t *testing.T) {
	resetCache()
	if err := Insert("mod", "k1", 1); err != nil {
		t.Fatal(err)
	}
	if err := Insert("mod", "k2", 2); err != nil {
		t.Fatal(err)
	}
}

func Test_Insert_DifferentModulesSameKey(t *testing.T) {
	resetCache()
	if err := Insert("mod1", "key", "a"); err != nil {
		t.Fatal(err)
	}
	if err := Insert("mod2", "key", "b"); err != nil {
		t.Fatal(err)
	}
}

func Test_Retrieve_ExistingEntry(t *testing.T) {
	resetCache()
	type item struct{ X int }
	want := item{X: 42}
	_ = Insert("m", "k", want)
	got := Retrieve("m", "k")
	if got == nil {
		t.Fatal("want value, got nil")
	}
	if got.(item) != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func Test_Retrieve_MissingModule(t *testing.T) {
	resetCache()
	if v := Retrieve("nonexistent", "k"); v != nil {
		t.Fatalf("want nil, got %v", v)
	}
}

func Test_Retrieve_MissingKey(t *testing.T) {
	resetCache()
	_ = Insert("m", "exists", true)
	if v := Retrieve("m", "missing"); v != nil {
		t.Fatalf("want nil, got %v", v)
	}
}

func Test_Cache_ConcurrentInserts(t *testing.T) {
	resetCache()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			_ = Insert("concurrent", fmt.Sprintf("k%d", i), i)
		}()
	}
	wg.Wait()
}

func Test_Cache_ConcurrentReadsAndWrites(t *testing.T) {
	resetCache()
	_ = Insert("rw", "shared", "val")
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		i := i
		go func() {
			defer wg.Done()
			_ = Retrieve("rw", "shared")
		}()
		go func() {
			defer wg.Done()
			_ = Insert("rw", fmt.Sprintf("new%d", i), i)
		}()
	}
	wg.Wait()
}

// Test_Cache keeps the original test for backward compatibility.
func Test_Cache(t *testing.T) {
	resetCache()
	type TestType struct{ Value string }
	initial := TestType{Value: "Hello, world"}
	if err := Insert("a", "b", initial); err != nil {
		t.Error("unexpected error")
	}
	if err := Insert("a", "b", initial); err != ErrAlreadyCached {
		t.Error("expected ErrAlreadyCached")
	}
	got := Retrieve("a", "b")
	if got == nil {
		t.Error("expected value, got nil")
	}
	if Retrieve("a", "c") != nil {
		t.Error("expected nil for missing key")
	}
}
