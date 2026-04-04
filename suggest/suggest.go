package suggest

import (
	"fmt"
	"sync"
)

var mu sync.RWMutex

// Represents the internal list of suggestions.
var suggestions Suggestions = []Suggestion{}

// Suggestion represents a suggestion made from a credo Module.
type Suggestion struct {
	// Module is the name of the module that the suggestion came from.
	Module string
	// From is the origin source of the suggestion, a.k.a. the package that
	// originated it.
	From string
	// Suggested is the suggested package.
	Suggested string
}

// Register adds a new suggestion to the internal list of suggestions.
func Register(suggest Suggestion) {
	mu.Lock()
	defer mu.Unlock()
	suggestions = append(suggestions, suggest)
}

// Get returns the internal list of all registered suggestions.
func Get() Suggestions {
	mu.RLock()
	defer mu.RUnlock()
	return suggestions
}

// Returns a formatted string representation of a suggestion.
func (s Suggestion) String() string {
	return fmt.Sprintf("%s suggested from %s coming from %s.",
		s.Suggested,
		s.From,
		s.Module)
}

// Suggestions represents a list of Suggetion
type Suggestions []Suggestion

// Returns a formatted string representation of the suggestions.
func (suggestions Suggestions) String() (output string) {
	for _, s := range suggestions {
		output += fmt.Sprintf("\t- %s\n", s.String())
	}
	return
}

// HasSuggestion checks if there are any registered suggestions.
func HasSuggestion() bool {
	mu.RLock()
	defer mu.RUnlock()
	return len(suggestions) > 0
}
