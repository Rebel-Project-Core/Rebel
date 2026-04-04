package modules

import (
	"os"
	"testing"
)

// TestMain sets up a temporary working directory for the entire modules test
// suite so that project.ProjectPath() creates directories in a safe location.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "rebel-modules-test-*")
	if err != nil {
		panic("modules TestMain: " + err.Error())
	}
	defer os.RemoveAll(dir)

	orig, err := os.Getwd()
	if err != nil {
		panic("modules TestMain getwd: " + err.Error())
	}
	if err := os.Chdir(dir); err != nil {
		panic("modules TestMain chdir: " + err.Error())
	}
	defer os.Chdir(orig)

	os.Exit(m.Run())
}
