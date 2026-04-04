package modules

import "fmt"

// installAptPackages resolves, saves, and applies a list of apt packages,
// committing each to config. Does nothing if the apt module is not registered.
func installAptPackages(config *Config, packages []string) error {
	if _, ok := Modules["apt"]; !ok {
		return nil
	}
	apt := aptModule{}
	for _, v := range packages {
		spell, err := apt.bareRun(aptSpell{Name: v})
		if err != nil {
			return fmt.Errorf("installAptPackages bareRun %q: %v", v, err)
		}
		if err = apt.Commit(config, spell); err != nil && err != ErrAlreadyPresent {
			return fmt.Errorf("installAptPackages commit %q: %v", v, err)
		}
		if err = apt.Save(spell); err != nil {
			return fmt.Errorf("installAptPackages save %q: %v", v, err)
		}
		if err = apt.Apply(spell); err != nil {
			return fmt.Errorf("installAptPackages apply %q: %v", v, err)
		}
	}
	return nil
}
