package cmd

import (
	"strings"
	"testing"
)

func Test_RootCmd_NotNil(t *testing.T) {
	if RootCmd == nil {
		t.Fatal("RootCmd should not be nil")
	}
}

func Test_RootCmd_Use(t *testing.T) {
	if RootCmd.Use != "credo" {
		t.Errorf("want Use=credo, got %q", RootCmd.Use)
	}
}

func Test_RootCmd_HasVersionSubcommand(t *testing.T) {
	for _, cmd := range RootCmd.Commands() {
		if cmd.Use == "version" {
			return
		}
	}
	t.Error("root command missing 'version' subcommand")
}

func Test_setup_ReturnsNewInstance(t *testing.T) {
	c := setup()
	if c == nil {
		t.Fatal("setup() returned nil")
	}
	if c.Use != "credo" {
		t.Errorf("want Use=credo, got %q", c.Use)
	}
}

func Test_setup_ShortDescriptionNotEmpty(t *testing.T) {
	c := setup()
	if c.Short == "" {
		t.Error("Short description should not be empty")
	}
}

func Test_setup_LongDescriptionNotEmpty(t *testing.T) {
	c := setup()
	if c.Long == "" {
		t.Error("Long description should not be empty")
	}
}

func Test_setup_LongContainsBioinformatics(t *testing.T) {
	c := setup()
	if !strings.Contains(strings.ToLower(c.Long), "bioinformatics") {
		t.Errorf("Long description should mention bioinformatics, got: %q", c.Long)
	}
}

func Test_setup_HasVersionSubcommand(t *testing.T) {
	c := setup()
	for _, sub := range c.Commands() {
		if sub.Use == "version" {
			return
		}
	}
	t.Error("setup() command missing 'version' subcommand")
}

func Test_setup_VersionSubcmdHasShort(t *testing.T) {
	c := setup()
	for _, sub := range c.Commands() {
		if sub.Use == "version" {
			if sub.Short == "" {
				t.Error("version subcommand should have a Short description")
			}
			return
		}
	}
}

func Test_setup_ExecuteVersionCmd(t *testing.T) {
	c := setup()
	c.SetArgs([]string{"version"})
	if err := c.Execute(); err != nil {
		t.Errorf("executing 'version' subcommand returned error: %v", err)
	}
}
