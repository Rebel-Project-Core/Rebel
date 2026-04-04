package modules

import (
	"credo/cache"
	"credo/logger"
	"credo/project"
	"fmt"
	"os"
	"path"
	"strings"

	gopip "github.com/CREDOProject/go-pip"
	"github.com/CREDOProject/go-pip/utils"
	pythonvenv "github.com/CREDOProject/go-pythonvenv"
	"github.com/CREDOProject/sharedutils/types"
	"github.com/spf13/cobra"
)

const pipModuleName = "pip"

const pipModuleShort = "Retrieves a python package and its dependencies."

const pipModuleExample = `
Install a pip package:
	credo pip numpy

Install a pip package pinning it to a version:
	credo pip numpy==1.26.0
`

// Registers the pipModule.
func init() { Register(pipModuleName, func() Module { return &pipModule{} }) }

// pipModule is used to manage the pip scope in the credospell configuration.
type pipModule struct{}

// Apply implements Module.
func (m *pipModule) Apply(anySpell any) error {
	converted, err := types.To[pipSpell](anySpell)
	if err != nil {
		return fmt.Errorf("Error converting pip spell, %v", err)
	}
	pipBinary, err := getPipBinary()
	if err != nil {
		return err
	}
	downloadPath, err := moduleDownloadPath(pipModuleName)
	if err != nil {
		return err
	}
	cmd, err := gopip.New(*pipBinary).
		Install(converted.Name).
		FindLinks(downloadPath).
		Seal()
	err = cmd.Run(&gopip.RunOptions{
		Output: os.Stdout,
	})
	if err != nil {
		return err
	}
	return nil
}

// BulkApply implements Module.
func (m *pipModule) BulkApply(config *Config) error {
	for _, ps := range config.Pip {
		err := m.Apply(ps)
		if err != nil {
			return err
		}
	}
	return nil
}

type pipSpell struct {
	Name                 string `yaml:"name"`
	ExternalDependencies Config `yaml:"external_dependencies,omitempty"`
}

// equals returns true if s and t have the same Name.
func (s pipSpell) equals(t equatable) bool {
	o, err := types.To[pipSpell](t)
	if err != nil {
		return false
	}
	return strings.Compare(s.Name, o.Name) == 0

}

// Commit implements Module.
func (m *pipModule) Commit(config *Config, result any) error {
	newEntry, err := types.To[pipSpell](result)
	if err != nil {
		return ErrConverting
	}
	if Contains(config.Pip, *newEntry) {
		return ErrAlreadyPresent
	}
	config.Pip = append(config.Pip, *newEntry)
	return nil
}

func setupPythonVenv(path string) (string, error) {
	venv, err := pythonvenv.Create(path)
	if err != nil {
		return "", err
	}
	return venv.Path, nil
}

func getPipBinary() (*string, error) {
	// Obtain the project path
	projectPath, err := project.ProjectPath()
	if err != nil {
		return nil, fmt.Errorf("getPipBinary, obtaining project path: %v", err)
	}

	venvPath, err := setupPythonVenv(path.Join(*projectPath, PipVenvDirectoryName))
	if err != nil {
		return nil, fmt.Errorf("getPipBinary, setting up venv: %v", err)
	}

	pipBinary, err := utils.PipBinaryFrom(path.Join(venvPath, "bin"))
	if err != nil {
		return nil, fmt.Errorf("getPipBinary, binary from path: %v", err)
	}
	return &pipBinary, nil
}

func (c *pipModule) installApt(config *Config) error {
	return installAptPackages(config, PipAptDependencies)
}

func (m *pipModule) bareRun(p pipSpell) (pipSpell, error) {
	if spell, ok := retrieveFromCache[pipSpell](pipModuleName, p.Name); ok {
		return *spell, nil
	}
	pipBinary, err := getPipBinary()
	if err != nil {
		return pipSpell{}, fmt.Errorf("bareRun, retrieving pip binary: %v", err)
	}

	cmd, err := gopip.New(*pipBinary).Install(p.Name).DryRun().Seal()
	if err != nil {
		return pipSpell{}, fmt.Errorf("bareRun, creating pip command: %v", err)
	}

	err = cmd.Run(&gopip.RunOptions{
		Output: os.Stdout,
	})
	if err != nil {
		return pipSpell{}, fmt.Errorf("bareRun, running pip command: %v", err)
	}
	_ = cache.Insert(pipModuleName, p.Name, p)
	return p, nil
}

// Save implements Module.
func (m *pipModule) Save(anySpell any) error {
	converted, err := types.To[pipSpell](anySpell)
	if err != nil {
		return fmt.Errorf("Error converting pip spell, %v", err)
	}
	if cache.Retrieve(pipModuleName, converted.Name) != nil {
		return nil
	}
	pipBinary, err := getPipBinary()
	if err != nil {
		return err
	}
	downloadPath, err := moduleDownloadPath(pipModuleName)
	if err != nil {
		return err
	}
	cmd, err := gopip.New(*pipBinary).
		Download(converted.Name, downloadPath).
		Seal()
	err = cmd.Run(&gopip.RunOptions{
		Output: os.Stdout,
	})
	if err != nil {
		return err
	}
	_ = cache.Insert(pipModuleName, converted.Name, true)
	return nil
}

// BulkSave implements Module.
func (m *pipModule) BulkSave(config *Config) error {
	for _, ps := range config.Pip {
		err := m.Save(ps)
		if err != nil {
			return err
		}
	}
	return nil
}

// cobraRun serves as an entry point to the bare run of the pipModule.
//
// Intended to be used by cobra.
func (m *pipModule) cobraRun(config *Config) func(*cobra.Command, []string) {
	return func(c *cobra.Command, args []string) {
		m.installApt(config)
		spell, err := m.bareRun(pipSpell{
			Name: args[0],
		})
		if err != nil {
			logger.Get().Fatal(err)
		}
		err = m.Commit(config, spell)
		if err != nil && err != ErrAlreadyPresent {
			logger.Get().Fatal(err)
		}
	}
}

// CliConfig implements Module.
func (m *pipModule) CliConfig(config *Config) *cobra.Command {
	return &cobra.Command{
		Args:    minArgsValidator(pipModuleName),
		Example: pipModuleExample,
		Run:     m.cobraRun(config),
		Short:   pipModuleShort,
		Use:     pipModuleName,
	}
}
