package modules

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

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

func downloadFile(url string, outputPath string) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: received status code %d", resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// BulkSave implements Module.
func (u *urlModule) BulkSave(config *Config) error {
	panic("unimplemented")
}

func (u *urlModule) cobraArgs() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s module requires at least one argument: the URL to download", urlModuleName)
		}
		return nil
	}
}

func (u *urlModule) cobraRun(config *Config) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		// Placeholder for actual implementation
		fmt.Printf("Downloading from URL: %s\n", args[0])
	}
}

// CliConfig implements Module.
func (u *urlModule) CliConfig(config *Config) *cobra.Command {
	return &cobra.Command{
		Args:    u.cobraArgs(),
		Example: urlModuleExample,
		Run:     u.cobraRun(config),
		Short:   urlModuleShort,
		Use:     urlModuleName,
	}
}

// Commit implements Module.
func (u *urlModule) Commit(config *Config, result any) error {
	panic("unimplemented")
}

// Save implements Module.
func (u *urlModule) Save(any) error {
	panic("unimplemented")
}
