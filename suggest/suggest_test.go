package suggest

import (
	"strings"
	"sync"
	"testing"
)

func resetSuggestions() {
	suggestions = []Suggestion{}
}

func Test_Register_AddsEntry(t *testing.T) {
	resetSuggestions()
	s := Suggestion{Module: "mod", From: "from", Suggested: "pkg"}
	Register(s)
	got := Get()
	if len(got) != 1 {
		t.Fatalf("want 1 suggestion, got %d", len(got))
	}
	if got[0] != s {
		t.Errorf("got %+v, want %+v", got[0], s)
	}
}

func Test_HasSuggestion_FalseWhenEmpty(t *testing.T) {
	resetSuggestions()
	if HasSuggestion() {
		t.Error("expected false on empty suggestions")
	}
}

func Test_HasSuggestion_TrueAfterRegister(t *testing.T) {
	resetSuggestions()
	Register(Suggestion{})
	if !HasSuggestion() {
		t.Error("expected true after Register")
	}
}

func Test_Get_ReturnsAllRegistered(t *testing.T) {
	resetSuggestions()
	Register(Suggestion{Module: "a"})
	Register(Suggestion{Module: "b"})
	Register(Suggestion{Module: "c"})
	if len(Get()) != 3 {
		t.Fatalf("want 3, got %d", len(Get()))
	}
}

func Test_Suggestion_String_ContainsAllFields(t *testing.T) {
	s := Suggestion{Module: "mymod", From: "origin", Suggested: "mypkg"}
	out := s.String()
	for _, field := range []string{"mymod", "origin", "mypkg"} {
		if !strings.Contains(out, field) {
			t.Errorf("String() missing %q: %s", field, out)
		}
	}
}

func Test_Suggestions_String_ContainsEntries(t *testing.T) {
	ss := Suggestions{
		{Module: "a", From: "b", Suggested: "c"},
		{Module: "x", From: "y", Suggested: "z"},
	}
	out := ss.String()
	for _, want := range []string{"c", "z"} {
		if !strings.Contains(out, want) {
			t.Errorf("Suggestions.String() missing %q: %s", want, out)
		}
	}
}

func Test_Suggestions_String_Empty(t *testing.T) {
	var ss Suggestions
	if ss.String() != "" {
		t.Errorf("empty Suggestions.String() should be empty, got %q", ss.String())
	}
}

func Test_Suggest_Concurrent(t *testing.T) {
	resetSuggestions()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Register(Suggestion{Module: "m", From: "f", Suggested: "s"})
		}()
	}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = Get()
			_ = HasSuggestion()
		}()
	}
	wg.Wait()
}

// Test_Suggest keeps the original test for backward compatibility.
func Test_Suggest(t *testing.T) {
	resetSuggestions()
	suggestion := Suggestion{Module: "Module", From: "From", Suggested: "Suggested"}
	Register(suggestion)
	if !HasSuggestion() {
		t.Error("Suggestion not registered.")
	}
	suggestions := Get()
	t.Log(suggestions.String())
}
