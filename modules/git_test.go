package modules

import (
	"testing"
)

// ============================================================
// cobraArgs (gitModule)
// ============================================================

func Test_gitModule_cobraArgs_NoArgs(t *testing.T) {
	m := &gitModule{}
	validator := m.cobraArgs()
	err := validator(nil, []string{})
	if err == nil {
		t.Error("expected error with no args")
	}
}

func Test_gitModule_cobraArgs_ValidHTTPS(t *testing.T) {
	m := &gitModule{}
	validator := m.cobraArgs()
	err := validator(nil, []string{"https://github.com/org/repo"})
	if err != nil {
		t.Errorf("unexpected error for valid HTTPS git URL: %v", err)
	}
}

func Test_gitModule_cobraArgs_ValidSSH(t *testing.T) {
	m := &gitModule{}
	validator := m.cobraArgs()
	err := validator(nil, []string{"git@github.com:org/repo.git"})
	if err != nil {
		t.Errorf("unexpected error for valid SSH git URL: %v", err)
	}
}

// ============================================================
// gitModule stub Apply / BulkApply
// ============================================================

func Test_gitModule_Save_WrongType(t *testing.T) {
	m := &gitModule{}
	if err := m.Save("not-a-spell"); err != ErrConverting {
		t.Fatalf("want ErrConverting, got %v", err)
	}
}

func Test_gitModule_BulkSave_WithEntries_FailsGracefully(t *testing.T) {
	// git.Save attempts a real clone; it will fail but exercises the loop path.
	cfg := &Config{
		Git: []gitSpell{
			{URL: "https://github.com/nonexistent-org-xyz/nonexistent-repo-abc.git"},
		},
	}
	_ = (&gitModule{}).BulkSave(cfg)
}
