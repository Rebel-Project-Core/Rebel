package modules

import (
	"credo/cache"
	"credo/logger"
	"credo/project"
	"fmt"
	"os"
	"path"

	"github.com/CREDOProject/sharedutils/types"
	"github.com/spf13/cobra"
)

// moduleDownloadPath returns the download directory for a module, creating it
// if it does not exist.
func moduleDownloadPath(moduleName string) (string, error) {
	p, err := project.ProjectPath()
	if err != nil {
		return "", err
	}
	dir := path.Join(*p, moduleName)
	if err := os.MkdirAll(dir, DirectoryPermissions); err != nil {
		return "", err
	}
	return dir, nil
}

// retrieveFromCache attempts to load a typed value from the cache.
// Returns (value, true) on a cache hit, (nil, false) on a miss or type error.
func retrieveFromCache[T any](moduleName, key string) (*T, bool) {
	raw := cache.Retrieve(moduleName, key)
	if raw == nil {
		return nil, false
	}
	v, err := types.To[T](raw)
	if err != nil {
		logger.Get().Printf("[%s/cache]: %v", moduleName, err)
		return nil, false
	}
	return v, true
}

// minArgsValidator returns a cobra args validator that requires at least one
// argument, reporting the module name in the error message.
func minArgsValidator(moduleName string) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s module requires at least one argument.", moduleName)
		}
		return nil
	}
}
