package modules

import (
	"credo/cache"
	"credo/logger"
	"fmt"
	"os"
	"strings"

	goconda "github.com/CREDOProject/go-conda"
	condautils "github.com/CREDOProject/go-conda/utils"
	"github.com/CREDOProject/sharedutils/types"
	"github.com/spf13/cobra"
)

const condaModuleName = "conda"

const condaModuleShort = "Retrieves a conda package and its dependencies."

const condaModuleExample = `
Install a conda package:
	credo conda numpy

Install a conda package from a channel:
	credo conda scipy --channel=bioconda
`

// Registers the condaModule.
func init() { Register(condaModuleName, func() Module { return &condaModule{} }) }

// condaModule is used to manage the conda scope in the credospell configuration.
type condaModule struct{}

// Apply implements Module.
func (c *condaModule) Apply(any) error {
	return nil
}

// BulkApply implements Module.
func (c *condaModule) BulkApply(config *Config) error {
	return nil
}

type condaSpell struct {
	Name                 string `yaml:"name"`
	Channel              string `yaml:"channel,omitempty"`
	ExternalDependencies Config `yaml:"external_dependencies,omitempty"`
}

// equals returns true if c and t have the same Name and Channel.
func (c condaSpell) equals(t equatable) bool {
	o, err := types.To[condaSpell](t)
	if err != nil {
		return false
	}
	return strings.Compare(c.Name, o.Name) == 0 &&
		strings.Compare(c.Channel, o.Channel) == 0
}

// BulkSave implements Module.
func (c *condaModule) BulkSave(config *Config) error {
	for _, cs := range config.Conda {
		if err := c.Save(cs); err != nil {
			return err
		}
	}
	return nil
}

// Function used to run the module from the command line.
// It serves as an entry point to the bare run of the condaModule.
//
// Intended to be used by cobra.
func (c *condaModule) cobraRun(config *Config) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		channel, _ := cmd.Flags().GetString("channel")
		spell, err := c.bareRun(condaSpell{
			Name:    args[0],
			Channel: channel,
		})
		if err != nil {
			logger.Get().Fatal(err)
		}
		err = c.Commit(config, spell)
		if err != nil {
			logger.Get().Fatal(err)
		}
	}
}

// CliConfig implements Module.
func (c *condaModule) CliConfig(config *Config) *cobra.Command {
	command := &cobra.Command{
		Short:   condaModuleShort,
		Example: condaModuleExample,
		Use:     condaModuleName,
		Run:     c.cobraRun(config),
		Args:    minArgsValidator(condaModuleName),
	}
	command.PersistentFlags().String("channel", "", "Conda channel to use.")
	return command
}

// Commit implements Module.
func (c *condaModule) Commit(config *Config, result any) error {
	newEntry, err := types.To[condaSpell](result)
	if err != nil {
		return ErrConverting
	}
	if Contains(config.Conda, *newEntry) {
		return ErrAlreadyPresent
	}
	config.Conda = append(config.Conda, *newEntry)
	return nil
}

func (c *condaModule) bareRun(p condaSpell) (condaSpell, error) {
	if spell, ok := retrieveFromCache[condaSpell](condaModuleName, p.Name); ok {
		return *spell, nil
	}
	condaBinary, err := condautils.DetectCondaBinary()
	if err != nil {
		return condaSpell{}, fmt.Errorf("conda binary not found: %v", err)
	}
	cmd, err := goconda.New(condaBinary, "", "").
		Install(&goconda.PackageInfo{
			PackageName: p.Name,
			Channel:     p.Channel,
		}).
		DryRun().
		Seal()
	err = cmd.Run(&goconda.RunOptions{
		Output: os.Stdout,
	})
	if err != nil {
		return condaSpell{}, err
	}
	_ = cache.Insert(condaModuleName, p.Name, p)
	return p, nil
}

// Save implements Module.
func (c *condaModule) Save(anySpell any) error {
	spell, err := types.To[condaSpell](anySpell)
	if err != nil {
		return ErrConverting
	}
	if cache.Retrieve(condaModuleName, spell.Name) != nil {
		return nil
	}
	condaBinary, err := condautils.DetectCondaBinary()
	if err != nil {
		return err
	}
	downloadPath, err := moduleDownloadPath(condaModuleName)
	if err != nil {
		return err
	}
	cmd, err := goconda.
		New(condaBinary, downloadPath, downloadPath).
		Download(&goconda.PackageInfo{
			PackageName: spell.Name,
			Channel:     spell.Channel,
		}, downloadPath).Seal()

	err = cmd.Run(&goconda.RunOptions{
		Output: os.Stdout,
	})
	if err == nil {
		_ = cache.Insert(condaModuleName, spell.Name, true)
	}
	return err
}
