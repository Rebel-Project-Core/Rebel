package modules

import "github.com/spf13/cobra"

const urlModuleName = "url"

const urlModuleShort = "Downloads a file from a URL."

const urlModuleExample = `
Download a file from a URL:
	credo url https://example.com/file.txt

Download a file from a URL and save it to a specific path:
	credo url https://example.com/file.txt --output /path/to/save/file.txt
`

// Registers the urlModule.
func init() { Register(urlModuleName, func() Module { return &urlModule{} }) }

// urlModule is used to manage the url scope in the credospell configuration.
type urlModule struct{}

// Apply implements Module.
func (u *urlModule) Apply(any) error {
	panic("unimplemented")
}

// BulkApply implements Module.
func (u *urlModule) BulkApply(config *Config) error {
	panic("unimplemented")
}

// BulkSave implements Module.
func (u *urlModule) BulkSave(config *Config) error {
	panic("unimplemented")
}

// CliConfig implements Module.
func (u *urlModule) CliConfig(config *Config) *cobra.Command {
	panic("unimplemented")
}

// Commit implements Module.
func (u *urlModule) Commit(config *Config, result any) error {
	panic("unimplemented")
}

// Save implements Module.
func (u *urlModule) Save(any) error {
	panic("unimplemented")
}
