package modules

import (
	"testing"
)

func Test_installAptPackages_NoAptModule_ReturnsNil(t *testing.T) {
	if _, registered := Modules[aptModuleName]; registered {
		t.Skip("apt module is registered (Debian/Ubuntu); test only valid on other systems")
	}
	// When apt is not registered, installAptPackages should be a no-op.
	if err := installAptPackages(&Config{}, []string{"curl", "wget"}); err != nil {
		t.Errorf("expected nil when apt not registered, got: %v", err)
	}
}

func Test_installAptPackages_EmptyPackageList_ReturnsNil(t *testing.T) {
	if err := installAptPackages(&Config{}, []string{}); err != nil {
		t.Errorf("expected nil for empty package list, got: %v", err)
	}
}
