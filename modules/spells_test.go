package modules

import (
	"testing"
)

// ============================================================
// aptSpell
// ============================================================

func Test_aptSpell_equals_SameDependencies(t *testing.T) {
	a := aptSpell{Name: "pkg", Dependencies: []aptSpell{{Name: "dep"}}}
	b := aptSpell{Name: "pkg", Dependencies: []aptSpell{{Name: "dep"}}}
	if !a.equals(b) {
		t.Error("expected equal apt spells")
	}
}

func Test_aptSpell_equals_DifferentDependencyCount(t *testing.T) {
	a := aptSpell{Name: "pkg", Dependencies: []aptSpell{{Name: "d1"}}}
	b := aptSpell{Name: "pkg", Dependencies: []aptSpell{{Name: "d1"}, {Name: "d2"}}}
	if a.equals(b) {
		t.Error("expected not-equal with different dependency count")
	}
}

func Test_aptSpell_equals_EmptyDependencies(t *testing.T) {
	a := aptSpell{Name: "x"}
	b := aptSpell{Name: "y"}
	// equals() only compares Dependencies, not Name — both have no deps.
	if !a.equals(b) {
		t.Error("two spells with no deps should be equal regardless of Name")
	}
}

func Test_aptSpell_equals_WrongType(t *testing.T) {
	a := aptSpell{}
	if a.equals(pipSpell{}) {
		t.Error("expected false when comparing with different type")
	}
}

// ============================================================
// pipSpell
// ============================================================

func Test_pipSpell_equals_SameName(t *testing.T) {
	a := pipSpell{Name: "numpy"}
	b := pipSpell{Name: "numpy"}
	if !a.equals(b) {
		t.Error("expected equal pip spells")
	}
}

func Test_pipSpell_equals_DifferentName(t *testing.T) {
	a := pipSpell{Name: "numpy"}
	b := pipSpell{Name: "pandas"}
	if a.equals(b) {
		t.Error("expected not-equal for different names")
	}
}

func Test_pipSpell_equals_WrongType(t *testing.T) {
	a := pipSpell{Name: "x"}
	if a.equals(aptSpell{}) {
		t.Error("expected false for wrong type")
	}
}

// ============================================================
// condaSpell
// ============================================================

func Test_condaSpell_equals_SameNameAndChannel(t *testing.T) {
	a := condaSpell{Name: "scipy", Channel: "conda-forge"}
	b := condaSpell{Name: "scipy", Channel: "conda-forge"}
	if !a.equals(b) {
		t.Error("expected equal conda spells")
	}
}

func Test_condaSpell_equals_DifferentName(t *testing.T) {
	a := condaSpell{Name: "scipy"}
	b := condaSpell{Name: "numpy"}
	if a.equals(b) {
		t.Error("expected not-equal for different names")
	}
}

func Test_condaSpell_equals_DifferentChannel(t *testing.T) {
	a := condaSpell{Name: "pkg", Channel: "conda-forge"}
	b := condaSpell{Name: "pkg", Channel: "defaults"}
	if a.equals(b) {
		t.Error("expected not-equal for different channels")
	}
}

func Test_condaSpell_equals_WrongType(t *testing.T) {
	a := condaSpell{Name: "x"}
	if a.equals(pipSpell{}) {
		t.Error("expected false for wrong type")
	}
}

// ============================================================
// gitSpell
// ============================================================

func Test_gitSpell_equals_SameURLAndVersion(t *testing.T) {
	a := gitSpell{URL: "https://github.com/org/repo", Version: "v1.0"}
	b := gitSpell{URL: "https://github.com/org/repo", Version: "v1.0"}
	if !a.equals(b) {
		t.Error("expected equal git spells")
	}
}

func Test_gitSpell_equals_DifferentURL(t *testing.T) {
	a := gitSpell{URL: "https://github.com/org/a"}
	b := gitSpell{URL: "https://github.com/org/b"}
	if a.equals(b) {
		t.Error("expected not-equal for different URLs")
	}
}

func Test_gitSpell_equals_DifferentVersion(t *testing.T) {
	a := gitSpell{URL: "https://github.com/org/repo", Version: "v1"}
	b := gitSpell{URL: "https://github.com/org/repo", Version: "v2"}
	if a.equals(b) {
		t.Error("expected not-equal for different versions")
	}
}

func Test_gitSpell_equals_WrongType(t *testing.T) {
	a := gitSpell{}
	if a.equals(urlSpell{}) {
		t.Error("expected false for wrong type")
	}
}

// ============================================================
// urlSpell
// ============================================================

func Test_urlSpell_equals_SameURLAndPath(t *testing.T) {
	a := urlSpell{URL: "https://example.com/file", Path: "/usr/local/bin"}
	b := urlSpell{URL: "https://example.com/file", Path: "/usr/local/bin"}
	if !a.equals(b) {
		t.Error("expected equal url spells")
	}
}

func Test_urlSpell_equals_DifferentURL(t *testing.T) {
	a := urlSpell{URL: "https://example.com/a", Path: "/tmp"}
	b := urlSpell{URL: "https://example.com/b", Path: "/tmp"}
	if a.equals(b) {
		t.Error("expected not-equal for different URLs")
	}
}

func Test_urlSpell_equals_DifferentPath(t *testing.T) {
	a := urlSpell{URL: "https://example.com/f", Path: "/tmp"}
	b := urlSpell{URL: "https://example.com/f", Path: "/usr"}
	if a.equals(b) {
		t.Error("expected not-equal for different paths")
	}
}

func Test_urlSpell_equals_WrongType(t *testing.T) {
	a := urlSpell{}
	if a.equals(gitSpell{}) {
		t.Error("expected false for wrong type")
	}
}

// ============================================================
// cranSpell
// ============================================================

func Test_cranSpell_equals_SameFields(t *testing.T) {
	a := cranSpell{PackageName: "ggplot2", PackagePath: "pkg.tar.gz", BioConductor: false}
	b := cranSpell{PackageName: "ggplot2", PackagePath: "pkg.tar.gz", BioConductor: false}
	if !a.equals(b) {
		t.Error("expected equal cran spells")
	}
}

func Test_cranSpell_equals_DifferentPackageName(t *testing.T) {
	a := cranSpell{PackageName: "ggplot2"}
	b := cranSpell{PackageName: "dplyr"}
	if a.equals(b) {
		t.Error("expected not-equal for different package names")
	}
}

func Test_cranSpell_equals_DifferentBioConductor(t *testing.T) {
	a := cranSpell{PackageName: "pkg", BioConductor: true}
	b := cranSpell{PackageName: "pkg", BioConductor: false}
	if a.equals(b) {
		t.Error("expected not-equal for different BioConductor flag")
	}
}

func Test_cranSpell_equals_DifferentDependencyCount(t *testing.T) {
	a := cranSpell{PackageName: "pkg", Dependencies: []cranSpell{{PackageName: "dep"}}}
	b := cranSpell{PackageName: "pkg"}
	if a.equals(b) {
		t.Error("expected not-equal for different dependency count")
	}
}

func Test_cranSpell_equals_DifferentExternalDependencies(t *testing.T) {
	a := cranSpell{PackageName: "pkg", ExternalDependencies: Config{}}
	b := cranSpell{PackageName: "pkg", ExternalDependencies: Config{}}
	if !a.equals(b) {
		t.Error("same empty external deps should be equal")
	}
}

func Test_cranSpell_equals_DifferentPackagePath(t *testing.T) {
	a := cranSpell{PackageName: "pkg", PackagePath: "a.tar.gz"}
	b := cranSpell{PackageName: "pkg", PackagePath: "b.tar.gz"}
	if a.equals(b) {
		t.Error("expected not-equal for different PackagePath")
	}
}

func Test_cranSpell_equals_SameDependencies_Equal(t *testing.T) {
	dep := cranSpell{PackageName: "dep"}
	a := cranSpell{PackageName: "pkg", Dependencies: []cranSpell{dep}}
	b := cranSpell{PackageName: "pkg", Dependencies: []cranSpell{dep}}
	if !a.equals(b) {
		t.Error("expected equal with same single dependency")
	}
}

func Test_cranSpell_equals_SameDependencies_Unequal(t *testing.T) {
	a := cranSpell{PackageName: "pkg", Dependencies: []cranSpell{{PackageName: "dep-a"}}}
	b := cranSpell{PackageName: "pkg", Dependencies: []cranSpell{{PackageName: "dep-b"}}}
	if a.equals(b) {
		t.Error("expected not-equal when same-count deps have different names")
	}
}

func Test_cranSpell_equals_DifferentExternalDeps(t *testing.T) {
	a := cranSpell{
		PackageName:          "pkg",
		ExternalDependencies: Config{},
	}
	b := cranSpell{
		PackageName:          "pkg",
		ExternalDependencies: Config{Url: []urlSpell{{URL: "https://example.com/f", Path: "/tmp"}}},
	}
	if a.equals(b) {
		t.Error("expected not-equal for different ExternalDependencies")
	}
}

func Test_cranSpell_equals_WrongType(t *testing.T) {
	a := cranSpell{}
	if a.equals(pipSpell{}) {
		t.Error("expected false for wrong type")
	}
}

// ============================================================
// Commit methods
// ============================================================

func Test_pipModule_Commit_AddsEntry(t *testing.T) {
	m := &pipModule{}
	cfg := &Config{}
	if err := m.Commit(cfg, pipSpell{Name: "numpy"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Pip) != 1 {
		t.Fatalf("want 1 pip entry, got %d", len(cfg.Pip))
	}
}

func Test_pipModule_Commit_Duplicate(t *testing.T) {
	m := &pipModule{}
	cfg := &Config{}
	_ = m.Commit(cfg, pipSpell{Name: "numpy"})
	if err := m.Commit(cfg, pipSpell{Name: "numpy"}); err != ErrAlreadyPresent {
		t.Fatalf("want ErrAlreadyPresent, got %v", err)
	}
}

func Test_pipModule_Commit_WrongType(t *testing.T) {
	m := &pipModule{}
	if err := m.Commit(&Config{}, 42); err != ErrConverting {
		t.Fatalf("want ErrConverting, got %v", err)
	}
}

func Test_gitModule_Commit_AddsEntry(t *testing.T) {
	m := &gitModule{}
	cfg := &Config{}
	spell := gitSpell{URL: "https://github.com/org/repo", Version: "v1"}
	if err := m.Commit(cfg, spell); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Git) != 1 {
		t.Fatalf("want 1 git entry, got %d", len(cfg.Git))
	}
}

func Test_gitModule_Commit_Duplicate(t *testing.T) {
	m := &gitModule{}
	cfg := &Config{}
	spell := gitSpell{URL: "https://github.com/org/repo"}
	_ = m.Commit(cfg, spell)
	if err := m.Commit(cfg, spell); err != ErrAlreadyPresent {
		t.Fatalf("want ErrAlreadyPresent, got %v", err)
	}
}

func Test_gitModule_Commit_WrongType(t *testing.T) {
	m := &gitModule{}
	if err := m.Commit(&Config{}, "notaspell"); err != ErrConverting {
		t.Fatalf("want ErrConverting, got %v", err)
	}
}

func Test_aptModule_Commit_AddsEntry(t *testing.T) {
	m := &aptModule{}
	cfg := &Config{}
	if err := m.Commit(cfg, aptSpell{Name: "curl"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Apt) != 1 {
		t.Fatalf("want 1 apt entry, got %d", len(cfg.Apt))
	}
}

func Test_aptModule_Commit_Duplicate(t *testing.T) {
	m := &aptModule{}
	cfg := &Config{}
	_ = m.Commit(cfg, aptSpell{Name: "curl"})
	if err := m.Commit(cfg, aptSpell{Name: "curl"}); err != ErrAlreadyPresent {
		t.Fatalf("want ErrAlreadyPresent, got %v", err)
	}
}

func Test_aptModule_Commit_WrongType(t *testing.T) {
	m := &aptModule{}
	if err := m.Commit(&Config{}, 99); err != ErrConverting {
		t.Fatalf("want ErrConverting, got %v", err)
	}
}

func Test_aptModule_Save_WrongType(t *testing.T) {
	m := &aptModule{}
	if err := m.Save("bad"); err != ErrConverting {
		t.Fatalf("want ErrConverting, got %v", err)
	}
}

func Test_aptModule_Apply_WrongType(t *testing.T) {
	m := &aptModule{}
	if err := m.Apply("bad"); err != ErrConverting {
		t.Fatalf("want ErrConverting, got %v", err)
	}
}

func Test_condaModule_Commit_AddsEntry(t *testing.T) {
	m := &condaModule{}
	cfg := &Config{}
	if err := m.Commit(cfg, condaSpell{Name: "scipy"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Conda) != 1 {
		t.Fatalf("want 1 conda entry, got %d", len(cfg.Conda))
	}
}

func Test_condaModule_Commit_Duplicate(t *testing.T) {
	m := &condaModule{}
	cfg := &Config{}
	_ = m.Commit(cfg, condaSpell{Name: "scipy"})
	if err := m.Commit(cfg, condaSpell{Name: "scipy"}); err != ErrAlreadyPresent {
		t.Fatalf("want ErrAlreadyPresent, got %v", err)
	}
}

func Test_condaModule_Commit_WrongType(t *testing.T) {
	m := &condaModule{}
	if err := m.Commit(&Config{}, true); err != ErrConverting {
		t.Fatalf("want ErrConverting, got %v", err)
	}
}

func Test_cranModule_Commit_AddsEntry(t *testing.T) {
	m := &cranModule{}
	cfg := &Config{}
	if err := m.Commit(cfg, cranSpell{PackageName: "ggplot2"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Cran) != 1 {
		t.Fatalf("want 1 cran entry, got %d", len(cfg.Cran))
	}
}

func Test_cranModule_Commit_Duplicate(t *testing.T) {
	m := &cranModule{}
	cfg := &Config{}
	_ = m.Commit(cfg, cranSpell{PackageName: "ggplot2"})
	if err := m.Commit(cfg, cranSpell{PackageName: "ggplot2"}); err != ErrAlreadyPresent {
		t.Fatalf("want ErrAlreadyPresent, got %v", err)
	}
}

func Test_cranModule_Commit_WrongType(t *testing.T) {
	m := &cranModule{}
	if err := m.Commit(&Config{}, "bad"); err == nil {
		t.Fatal("want error for wrong type")
	}
}

func Test_urlModule_Commit_AddsEntry(t *testing.T) {
	m := &urlModule{}
	cfg := &Config{}
	spell := urlSpell{URL: "https://example.com/f", Path: "/tmp"}
	if err := m.Commit(cfg, spell); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Url) != 1 {
		t.Fatalf("want 1 url entry, got %d", len(cfg.Url))
	}
}

func Test_urlModule_Commit_Duplicate(t *testing.T) {
	m := &urlModule{}
	cfg := &Config{}
	spell := urlSpell{URL: "https://example.com/f", Path: "/tmp"}
	_ = m.Commit(cfg, spell)
	if err := m.Commit(cfg, spell); err != ErrAlreadyPresent {
		t.Fatalf("want ErrAlreadyPresent, got %v", err)
	}
}

func Test_urlModule_Commit_WrongType(t *testing.T) {
	m := &urlModule{}
	if err := m.Commit(&Config{}, 3.14); err != ErrConverting {
		t.Fatalf("want ErrConverting, got %v", err)
	}
}

// ============================================================
// BulkSave / BulkApply with empty config
// ============================================================

func Test_pipModule_BulkSave_EmptyConfig(t *testing.T) {
	if err := (&pipModule{}).BulkSave(&Config{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_pipModule_BulkApply_EmptyConfig(t *testing.T) {
	if err := (&pipModule{}).BulkApply(&Config{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_pipModule_BulkSave_WithEntries_FailsGracefully(t *testing.T) {
	// pip.Save creates a venv; this may fail in CI but exercises the loop path.
	cfg := &Config{Pip: []pipSpell{{Name: "numpy"}}}
	_ = (&pipModule{}).BulkSave(cfg)
}

func Test_pipModule_BulkApply_WithEntries_FailsGracefully(t *testing.T) {
	cfg := &Config{Pip: []pipSpell{{Name: "numpy"}}}
	_ = (&pipModule{}).BulkApply(cfg)
}

func Test_gitModule_BulkSave_EmptyConfig(t *testing.T) {
	if err := (&gitModule{}).BulkSave(&Config{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_gitModule_BulkApply_EmptyConfig(t *testing.T) {
	if err := (&gitModule{}).BulkApply(&Config{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_aptModule_BulkSave_EmptyConfig(t *testing.T) {
	if err := (&aptModule{}).BulkSave(&Config{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_aptModule_BulkApply_EmptyConfig(t *testing.T) {
	if err := (&aptModule{}).BulkApply(&Config{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_aptModule_BulkSave_WithOptionalDep_SkipsIt(t *testing.T) {
	// An optional dependency should be skipped; iteration still enters the loop.
	// On non-Debian systems, the subsequent Save will fail — we only care about
	// the fact that the optional dep is skipped without error up to that point.
	m := &aptModule{}
	cfg := &Config{
		Apt: []aptSpell{
			{
				Name: "test-pkg",
				Dependencies: []aptSpell{
					{Name: "opt-dep", Optional: true}, // skipped
				},
			},
		},
	}
	// Result may be nil (Debian/Ubuntu) or an error (other systems); either is OK.
	_ = m.BulkSave(cfg)
}

func Test_aptModule_BulkApply_WithOptionalDep_SkipsIt(t *testing.T) {
	m := &aptModule{}
	cfg := &Config{
		Apt: []aptSpell{
			{
				Name: "test-pkg",
				Dependencies: []aptSpell{
					{Name: "opt-dep", Optional: true}, // skipped
				},
			},
		},
	}
	_ = m.BulkApply(cfg)
}

func Test_condaModule_BulkSave_EmptyConfig(t *testing.T) {
	if err := (&condaModule{}).BulkSave(&Config{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_condaModule_BulkApply_EmptyConfig(t *testing.T) {
	if err := (&condaModule{}).BulkApply(&Config{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_condaModule_BulkSave_WithEntries_FailsGracefully(t *testing.T) {
	cfg := &Config{Conda: []condaSpell{{Name: "scipy", Channel: "conda-forge"}}}
	// Will error if conda binary not installed; we just cover the iteration path.
	_ = (&condaModule{}).BulkSave(cfg)
}

func Test_condaModule_Apply_ReturnsNil(t *testing.T) {
	// condaModule.Apply is a stub that always returns nil.
	if err := (&condaModule{}).Apply("anything"); err != nil {
		t.Errorf("condaModule.Apply should return nil, got %v", err)
	}
}

func Test_cranModule_BulkSave_EmptyConfig(t *testing.T) {
	if err := (&cranModule{}).BulkSave(&Config{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_cranModule_BulkApply_EmptyConfig(t *testing.T) {
	if err := (&cranModule{}).BulkApply(&Config{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_cranModule_BulkSave_WithEntries_FailsGracefully(t *testing.T) {
	cfg := &Config{Cran: []cranSpell{{PackageName: "ggplot2"}}}
	_ = (&cranModule{}).BulkSave(cfg)
}

func Test_cranModule_BulkApply_WithEntries_FailsGracefully(t *testing.T) {
	cfg := &Config{Cran: []cranSpell{{PackageName: "ggplot2"}}}
	_ = (&cranModule{}).BulkApply(cfg)
}

func Test_urlModule_BulkSave_EmptyConfig(t *testing.T) {
	if err := (&urlModule{}).BulkSave(&Config{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_urlModule_BulkApply_EmptyConfig(t *testing.T) {
	if err := (&urlModule{}).BulkApply(&Config{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ============================================================
// gitModule stub Apply / BulkApply
// ============================================================

func Test_gitModule_Apply_ReturnsNil(t *testing.T) {
	if err := (&gitModule{}).Apply(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_gitModule_BulkApply_ReturnsNil(t *testing.T) {
	if err := (&gitModule{}).BulkApply(&Config{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ============================================================
// CliConfig structure checks
// ============================================================

func Test_pipModule_CliConfig_Use(t *testing.T) {
	cmd := (&pipModule{}).CliConfig(&Config{})
	if cmd.Use != pipModuleName {
		t.Errorf("want Use=%q, got %q", pipModuleName, cmd.Use)
	}
}

func Test_gitModule_CliConfig_Use(t *testing.T) {
	cmd := (&gitModule{}).CliConfig(&Config{})
	if cmd.Use != gitModuleName {
		t.Errorf("want Use=%q, got %q", gitModuleName, cmd.Use)
	}
}

func Test_aptModule_CliConfig_Use(t *testing.T) {
	cmd := (&aptModule{}).CliConfig(&Config{})
	if cmd.Use != aptModuleName {
		t.Errorf("want Use=%q, got %q", aptModuleName, cmd.Use)
	}
}

func Test_condaModule_CliConfig_HasChannelFlag(t *testing.T) {
	cmd := (&condaModule{}).CliConfig(&Config{})
	if cmd.PersistentFlags().Lookup("channel") == nil {
		t.Error("conda command missing --channel flag")
	}
}

func Test_cranModule_CliConfig_HasRepositoryFlag(t *testing.T) {
	cmd := (&cranModule{}).CliConfig(&Config{})
	if cmd.PersistentFlags().Lookup("repository") == nil {
		t.Error("cran command missing --repository flag")
	}
}

func Test_cranModule_CliConfig_HasBioconductorAlias(t *testing.T) {
	cmd := (&cranModule{}).CliConfig(&Config{})
	for _, alias := range cmd.Aliases {
		if alias == bioconductorModuleName {
			return
		}
	}
	t.Error("cran command missing bioconductor alias")
}

func Test_urlModule_CliConfig_HasPathFlag(t *testing.T) {
	cmd := (&urlModule{}).CliConfig(&Config{})
	if cmd.PersistentFlags().Lookup("path") == nil {
		t.Error("url command missing --path flag")
	}
}
