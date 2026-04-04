package modules

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

// ---- helpers ----

type testSpell struct{ Name string }

func (t testSpell) equals(e equatable) bool {
	o, err := to[testSpell](e)
	if err != nil {
		return false
	}
	return t.Name == o.Name
}

// to is a local generic helper mirroring sharedutils/types.To.
func to[T any](v equatable) (*T, error) {
	if t, ok := v.(T); ok {
		return &t, nil
	}
	return nil, errors.New("type assertion failed")
}

type stubModule struct {
	bulkSaveErr  error
	bulkApplyErr error
}

func (s stubModule) Commit(*Config, any) error        { return nil }
func (s stubModule) Save(any) error                   { return nil }
func (s stubModule) BulkSave(*Config) error           { return s.bulkSaveErr }
func (s stubModule) Apply(any) error                  { return nil }
func (s stubModule) BulkApply(*Config) error          { return s.bulkApplyErr }
func (s stubModule) CliConfig(*Config) *cobra.Command { return &cobra.Command{Use: "stub"} }

// ---- Contains ----

func Test_Contains_EmptySlice(t *testing.T) {
	if Contains([]testSpell{}, testSpell{"x"}) {
		t.Error("Contains on empty slice should return false")
	}
}

func Test_Contains_NilSlice(t *testing.T) {
	if Contains[testSpell](nil, testSpell{"x"}) {
		t.Error("Contains on nil slice should return false")
	}
}

func Test_Contains_Match(t *testing.T) {
	s := []testSpell{{"a"}, {"b"}, {"c"}}
	if !Contains(s, testSpell{"b"}) {
		t.Error("Contains should return true for existing element")
	}
}

func Test_Contains_NoMatch(t *testing.T) {
	s := []testSpell{{"a"}, {"b"}}
	if Contains(s, testSpell{"z"}) {
		t.Error("Contains should return false for missing element")
	}
}

// ---- Register ----

func Test_Register_NewModule(t *testing.T) {
	const testName = "test-register-unique"
	defer delete(Modules, testName)

	Register(testName, func() Module { return stubModule{} })
	if _, ok := Modules[testName]; !ok {
		t.Error("module not found in registry after Register")
	}
}

// ---- RegisterModulesCli ----

func Test_RegisterModulesCli_AddsCmds(t *testing.T) {
	root := &cobra.Command{Use: "root"}
	cfg := &Config{}
	RegisterModulesCli(root, cfg)
	if len(root.Commands()) == 0 {
		t.Error("RegisterModulesCli should add at least one subcommand")
	}
}

// ---- DeepSave ----

func Test_DeepSave_EmptyConfig_NoError(t *testing.T) {
	if err := DeepSave(&Config{}); err != nil {
		t.Fatalf("DeepSave on empty config should not error, got: %v", err)
	}
}

func Test_DeepSave_PropagatesError(t *testing.T) {
	const testName = "test-deepsave-err"
	want := errors.New("save failed")
	Modules[testName] = func() Module { return stubModule{bulkSaveErr: want} }
	defer delete(Modules, testName)

	if err := DeepSave(&Config{}); !errors.Is(err, want) {
		t.Errorf("DeepSave should propagate module error, got: %v", err)
	}
}

// ---- DeepApply ----

func Test_DeepApply_EmptyConfig_NoError(t *testing.T) {
	if err := DeepApply(&Config{}); err != nil {
		t.Fatalf("DeepApply on empty config should not error, got: %v", err)
	}
}

func Test_DeepApply_PropagatesError(t *testing.T) {
	const testName = "test-deepapply-err"
	want := errors.New("apply failed")
	Modules[testName] = func() Module { return stubModule{bulkApplyErr: want} }
	defer delete(Modules, testName)

	if err := DeepApply(&Config{}); !errors.Is(err, want) {
		t.Errorf("DeepApply should propagate module error, got: %v", err)
	}
}
