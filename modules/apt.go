package modules

import (
	"credo/cache"
	"credo/logger"
	"credo/suggest"
	"fmt"
	"regexp"
	"sync"

	"github.com/CREDOProject/go-apt-client"
	goosinfo "github.com/CREDOProject/go-osinfo"
	"github.com/CREDOProject/sharedutils/types"
	"github.com/spf13/cobra"
)

const aptModuleName = "apt"

var aptMutex sync.Mutex

const aptModuleShort = "Retrieves an apt package and its dependencies."

const aptModuleExample = `
Install a apt package:
	credo apt python3
`

var isAptOptional = regexp.MustCompile(`\<(?P<name>..*)\>`)

// Registers the aptModule.
func init() {
	osinfo, err := goosinfo.Retrieve()
	if err != nil {
		logger.Get().Fatal(err)
	}
	supportedDistributions := map[string]struct{}{
		"ubuntu": {},
		"debian": {},
	}
	for _, distribution := range osinfo.Like {
		if _, ok := supportedDistributions[distribution]; ok {
			Register(aptModuleName, func() Module { return &aptModule{} })
			return
		}
	}
	if _, ok := supportedDistributions[osinfo.Distribution]; ok {
		Register(aptModuleName, func() Module { return &aptModule{} })
		return
	}
}

// aptModule is used to manage the apt scope in the credospell configuration.
type aptModule struct{}

type aptSpell struct {
	Name                 string     `yaml:"name"`
	Optional             bool       `yaml:"optional,omitempty"`
	Dependencies         []aptSpell `yaml:"dependencies,omitempty"`
	ExternalDependencies Config     `yaml:"external_dependencies,omitempty"`
}

// equals returns true if a and t have the same Name, Optional flag, and Dependencies.
func (a aptSpell) equals(t equatable) bool {
	o, err := types.To[aptSpell](t)
	if err != nil {
		return false
	}
	equality := len(o.Dependencies) == len(a.Dependencies)
	if !equality {
		return false
	}
	for i := range o.Dependencies {
		equality = equality &&
			o.Dependencies[i].equals(a.Dependencies[i])
	}
	return equality
}

// BulkSave implements Module.
func (m *aptModule) BulkSave(config *Config) error {
	for _, as := range config.Apt {
		for _, dep := range as.Dependencies {
			if dep.Optional {
				continue
			}
			err := m.Save(dep)
			if err != nil {
				return err
			}
		}
		err := m.Save(as)
		if err != nil {
			return err
		}
	}
	return nil
}

// CliConfig implements Module.
func (m *aptModule) CliConfig(config *Config) *cobra.Command {
	return &cobra.Command{
		Args:    minArgsValidator(aptModuleName),
		Example: aptModuleExample,
		Run:     m.cobraRun(config),
		Short:   aptModuleShort,
		Use:     aptModuleName,
	}
}

// Function used to run the module from the command line.
// It serves as an entry point to the bare run of the aptModule.
//
// Intended to be used by cobra.
func (m *aptModule) cobraRun(config *Config) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		name := args[0]
		spell, err := m.bareRun(aptSpell{
			Name: name,
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

func (*aptModule) bareRun(s aptSpell) (aptSpell, error) {
	if spell, ok := retrieveFromCache[aptSpell](aptModuleName, s.Name); ok {
		return *spell, nil
	}
	aptPack := &apt.Package{
		Name: s.Name,
	}

	aptMutex.Lock()
	_, err := apt.CheckForUpdates()
	if err != nil {
		return aptSpell{}, fmt.Errorf("While running: %s, failed to check for updates: %w", s.Name, err)
	}
	aptMutex.Unlock()

	output, err := apt.InstallDry(aptPack)
	logger.Get().Print(string(output))
	if err != nil {
		return aptSpell{}, err
	}
	depList, err := apt.GetDependencies(aptPack)
	if err != nil {
		return aptSpell{}, err
	}
	for _, dependency := range depList {
		isOptional := isAptOptional.MatchString(dependency)
		cleanDependency := dependency
		matches := isAptOptional.FindStringSubmatch(dependency)
		nameIndex := isAptOptional.SubexpIndex("name")
		if nameIndex != -1 && isOptional {
			cleanDependency = matches[nameIndex]
			suggest.Register(suggest.Suggestion{
				Module:    aptModuleName,
				From:      aptPack.Name,
				Suggested: cleanDependency,
			})
		}
		s.Dependencies = append(s.Dependencies, aptSpell{
			Name:     cleanDependency,
			Optional: isOptional,
		})
	}
	_ = cache.Insert(aptModuleName, s.Name, s)
	return s, nil
}

// Commit implements Module.
func (*aptModule) Commit(config *Config, result any) error {
	newEntry, err := types.To[aptSpell](result)
	if err != nil {
		return ErrConverting
	}
	if Contains(config.Apt, *newEntry) {
		return ErrAlreadyPresent
	}
	config.Apt = append(config.Apt, *newEntry)
	return nil
}

// Save implements Module.
func (*aptModule) Save(anySpell any) error {
	spell, err := types.To[aptSpell](anySpell)
	if err != nil {
		return ErrConverting
	}
	if cache.Retrieve(aptModuleName, spell.Name) != nil {
		return nil
	}
	downloadPath, err := moduleDownloadPath(aptModuleName)
	if err != nil {
		return err
	}
	aptPack := &apt.Package{
		Name: spell.Name,
	}
	out, err := apt.Download(aptPack, downloadPath)
	logger.Get().Print(string(out))
	if err == nil {
		_ = cache.Insert(aptModuleName, spell.Name, true)
	}
	return err
}

// Apply implements Module.
func (m *aptModule) Apply(anySpell any) error {
	spell, err := types.To[aptSpell](anySpell)
	if err != nil {
		return ErrConverting
	}
	downloadPath, err := moduleDownloadPath(aptModuleName)
	if err != nil {
		return err
	}
	aptPack := &apt.Package{
		Name: spell.Name,
	}
	out, err := apt.Install(downloadPath, aptPack)
	logger.Get().Print(string(out))
	return err
}

// BulkApply implements Module.
func (m *aptModule) BulkApply(config *Config) error {
	for _, as := range config.Apt {
		for _, dep := range as.Dependencies {
			if dep.Optional {
				continue
			}
			err := m.Apply(dep)
			if err != nil {
				return err
			}
		}
		err := m.Apply(as)
		if err != nil {
			return err
		}
	}
	return nil
}
