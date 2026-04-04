package modules

import (
	"testing"
)

// ============================================================
// applyModule (stub)
// ============================================================

func Test_applyModule_Apply_ReturnsNil(t *testing.T) {
	a := applyModule{}
	if err := a.Apply(nil); err != nil {
		t.Errorf("Apply() = %v, want nil", err)
	}
}

func Test_applyModule_BulkApply_ReturnsNil(t *testing.T) {
	a := applyModule{}
	if err := a.BulkApply(&Config{}); err != nil {
		t.Errorf("BulkApply() = %v, want nil", err)
	}
}

func Test_applyModule_BulkSave_ReturnsNil(t *testing.T) {
	a := applyModule{}
	if err := a.BulkSave(&Config{}); err != nil {
		t.Errorf("BulkSave() = %v, want nil", err)
	}
}

func Test_applyModule_Commit_ReturnsNil(t *testing.T) {
	a := applyModule{}
	if err := a.Commit(&Config{}, nil); err != nil {
		t.Errorf("Commit() = %v, want nil", err)
	}
}

func Test_applyModule_Save_ReturnsNil(t *testing.T) {
	a := applyModule{}
	if err := a.Save(nil); err != nil {
		t.Errorf("Save() = %v, want nil", err)
	}
}

func Test_applyModule_CliConfig_Use(t *testing.T) {
	a := applyModule{}
	cmd := a.CliConfig(&Config{})
	if cmd.Use != applyModuleName {
		t.Errorf("want Use=%q, got %q", applyModuleName, cmd.Use)
	}
}

// ============================================================
// saveModule (stub)
// ============================================================

func Test_saveModule_Apply_ReturnsNil(t *testing.T) {
	s := &saveModule{}
	if err := s.Apply(nil); err != nil {
		t.Errorf("Apply() = %v, want nil", err)
	}
}

func Test_saveModule_BulkApply_ReturnsNil(t *testing.T) {
	s := &saveModule{}
	if err := s.BulkApply(&Config{}); err != nil {
		t.Errorf("BulkApply() = %v, want nil", err)
	}
}

func Test_saveModule_BulkSave_ReturnsNil(t *testing.T) {
	s := &saveModule{}
	if err := s.BulkSave(&Config{}); err != nil {
		t.Errorf("BulkSave() = %v, want nil", err)
	}
}

func Test_saveModule_Commit_ReturnsNil(t *testing.T) {
	s := &saveModule{}
	if err := s.Commit(&Config{}, nil); err != nil {
		t.Errorf("Commit() = %v, want nil", err)
	}
}

func Test_saveModule_Save_ReturnsNil(t *testing.T) {
	s := &saveModule{}
	if err := s.Save(nil); err != nil {
		t.Errorf("Save() = %v, want nil", err)
	}
}

func Test_saveModule_CliConfig_Use(t *testing.T) {
	s := &saveModule{}
	cmd := s.CliConfig(&Config{})
	if cmd.Use != saveModuleName {
		t.Errorf("want Use=%q, got %q", saveModuleName, cmd.Use)
	}
}

// ============================================================
// Execute CliConfig Run functions with empty config (no-op)
// ============================================================

func Test_applyModule_CliConfig_RunWithEmptyConfig(t *testing.T) {
	a := applyModule{}
	cfg := &Config{}
	cmd := a.CliConfig(cfg)
	// Run should succeed with empty config (DeepApply iterates empty slices).
	cmd.Run(cmd, []string{})
}

func Test_saveModule_CliConfig_RunWithEmptyConfig(t *testing.T) {
	s := &saveModule{}
	cfg := &Config{}
	cmd := s.CliConfig(cfg)
	// Run should succeed with empty config (DeepSave iterates empty slices).
	cmd.Run(cmd, []string{})
}
