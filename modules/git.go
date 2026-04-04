package modules

import (
	"credo/logger"
	"credo/project"
	"fmt"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/cobra"

	goisgiturl "github.com/CREDOProject/go-isgiturl"
	"github.com/CREDOProject/sharedutils/types"
)

const gitModuleName = "git"

const gitModuleShort = "Retrieves a remote git repository."

const gitModuleExample = `
Clone a git repository:
	credo git https://github.com/kendomaniac/rCASC

Clone a git repository at a specific version tag:
	credo git https://github.com/kendomaniac/docker4seq 2.1.2
`

// Registers the gitModule.
func init() { Register(gitModuleName, func() Module { return &gitModule{} }) }

// gitModule is used to manage the git scope in the credospell configuration.
type gitModule struct{}

// Apply implements Module.
func (m *gitModule) Apply(any) error {
	return nil
}

// BulkApply implements Module.
func (m *gitModule) BulkApply(config *Config) error {
	return nil
}

func (m *gitModule) Commit(config *Config, result any) error {
	newEntry, err := types.To[gitSpell](result)
	if err != nil {
		return ErrConverting
	}
	if Contains(config.Git, *newEntry) {
		return ErrAlreadyPresent
	}
	config.Git = append(config.Git, *newEntry)
	return nil
}

func (m *gitModule) bareRun(p gitSpell) (gitSpell, error) {
	version := p.Version

	var ref plumbing.ReferenceName
	if len(version) > 0 && version != "HEAD" {
		ref = plumbing.ReferenceName("refs/tags/" + version)
	}

	_, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:               p.URL,
		Depth:             1,
		SingleBranch:      true,
		RecurseSubmodules: 1,
		ReferenceName:     ref,
	})

	if err != nil {
		return gitSpell{URL: p.URL, Version: version}, err
	}
	return gitSpell{URL: p.URL, Version: version}, nil
}

func (m *gitModule) Save(anySpell any) error {
	spell, err := types.To[gitSpell](anySpell)
	if err != nil {
		return ErrConverting
	}

	projectPathStr, err := project.ProjectPath()
	if err != nil {
		return err
	}

	_, _, _, repoPath := goisgiturl.FindScpLikeComponents(spell.URL)
	joinedPath := path.Join(strings.Split(repoPath, "/")...)

	cloneOptions := &git.CloneOptions{
		URL:               spell.URL,
		Depth:             1,
		SingleBranch:      true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	if len(spell.Version) > 0 && spell.Version != "HEAD" {
		cloneOptions.ReferenceName = plumbing.ReferenceName("refs/tags/" + spell.Version)
	}

	_, err = git.PlainClone(path.Join(projectPathStr, gitModuleName, joinedPath), false, cloneOptions)
	return err
}

func (m *gitModule) BulkSave(config *Config) error {
	for _, gs := range config.Git {
		err := m.Save(gs)
		if err != nil {
			return err
		}
	}
	return nil
}

// Struct containing a Spell Entry for a Git repo.
type gitSpell struct {
	URL                  string `yaml:"url"`
	Version              string `yaml:"version"`
	ExternalDependencies Config `yaml:"external_dependencies,omitempty"`
}

// equals returns true if s and t have the same URL and Version.
func (s gitSpell) equals(t equatable) bool {
	o, err := types.To[gitSpell](t)
	if err != nil {
		return false
	}
	return strings.Compare(s.URL, o.URL) == 0 &&
		strings.Compare(s.Version, o.Version) == 0
}

// Function used to run the module from the command line.
// It serves as an entry point to the bare run of the gitModule.
//
// Intended to be used by cobra.
func (m *gitModule) cobraRun(config *Config) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		version := ""
		if len(args) > 1 {
			version = args[1]
		}
		spell, err := m.bareRun(gitSpell{
			URL:     args[0],
			Version: version,
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

// Function used to validate the arguments passed to the git command.
// If no arguments are passed, it returns an error.
// Otherwise it returns nil.
//
// Intended to be used by cobra.
func (m *gitModule) cobraArgs() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s module requires at least one argument.",
				gitModuleName)
		}
		url := args[0]
		if !goisgiturl.IsGitUrl(url) {
			return fmt.Errorf("\"%s\" doesn't look like a git uri.", url)
		}
		return nil
	}
}

// CliConfig implements Module.
func (m *gitModule) CliConfig(config *Config) *cobra.Command {
	return &cobra.Command{
		Use:     gitModuleName,
		Short:   gitModuleShort,
		Example: gitModuleExample,
		Args:    m.cobraArgs(),
		Run:     m.cobraRun(config),
	}
}
