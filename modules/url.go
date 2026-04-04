package modules

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"credo/cache"
	"credo/logger"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/CREDOProject/sharedutils/types"
	"github.com/spf13/cobra"
)

const urlModuleName = "url"

const urlModuleShort = "Downloads a file from a URL and places it at a given path."

const urlModuleExample = `
Download a binary:
	credo url https://example.com/tool --path /usr/local/bin/tool

Download and extract an archive:
	credo url https://example.com/tool.tar.gz --path /usr/local/bin/
`

// Registers the urlModule.
func init() { Register(urlModuleName, func() Module { return &urlModule{} }) }

// urlModule manages downloading files (and archives) from arbitrary URLs.
type urlModule struct{}

type urlSpell struct {
	URL                  string `yaml:"url"`
	Path                 string `yaml:"path"`
	ExternalDependencies Config `yaml:"external_dependencies,omitempty"`
}

// equals returns true if u and t have the same URL and destination Path.
func (u urlSpell) equals(t equatable) bool {
	o, err := types.To[urlSpell](t)
	if err != nil {
		return false
	}
	return u.URL == o.URL && u.Path == o.Path
}

// Commit implements Module.
func (m *urlModule) Commit(config *Config, result any) error {
	newEntry, err := types.To[urlSpell](result)
	if err != nil {
		return ErrConverting
	}
	if Contains(config.Url, *newEntry) {
		return ErrAlreadyPresent
	}
	config.Url = append(config.Url, *newEntry)
	return nil
}

// Save implements Module — downloads the file into the module's download dir.
func (m *urlModule) Save(anySpell any) error {
	spell, err := types.To[urlSpell](anySpell)
	if err != nil {
		return ErrConverting
	}
	if cache.Retrieve(urlModuleName, spell.URL) != nil {
		return nil
	}
	downloadDir, err := moduleDownloadPath(urlModuleName)
	if err != nil {
		return err
	}
	dest := filepath.Join(downloadDir, path.Base(spell.URL))
	if err := downloadFile(spell.URL, dest); err != nil {
		return err
	}
	_ = cache.Insert(urlModuleName, spell.URL, true)
	return nil
}

// BulkSave implements Module.
func (m *urlModule) BulkSave(config *Config) error {
	for _, us := range config.Url {
		if err := m.Save(us); err != nil {
			return err
		}
	}
	return nil
}

// Apply implements Module — places the downloaded file at the destination path,
// decompressing zip and tar.gz archives automatically.
func (m *urlModule) Apply(anySpell any) error {
	spell, err := types.To[urlSpell](anySpell)
	if err != nil {
		return ErrConverting
	}
	downloadDir, err := moduleDownloadPath(urlModuleName)
	if err != nil {
		return err
	}
	src := filepath.Join(downloadDir, path.Base(spell.URL))
	if isArchive(spell.URL) {
		return extractArchive(src, spell.Path)
	}
	return copyExecutable(src, spell.Path)
}

// BulkApply implements Module.
func (m *urlModule) BulkApply(config *Config) error {
	for _, us := range config.Url {
		if err := m.Apply(us); err != nil {
			return err
		}
	}
	return nil
}

func (m *urlModule) bareRun(p urlSpell) (urlSpell, error) {
	if spell, ok := retrieveFromCache[urlSpell](urlModuleName, p.URL); ok {
		return *spell, nil
	}
	resp, err := http.Head(p.URL)
	if err != nil {
		return urlSpell{}, fmt.Errorf("url: HEAD %s: %v", p.URL, err)
	}
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return urlSpell{}, fmt.Errorf("url: %s returned HTTP %d", p.URL, resp.StatusCode)
	}
	_ = cache.Insert(urlModuleName, p.URL, p)
	return p, nil
}

func (m *urlModule) cobraRun(config *Config) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		destPath, _ := cmd.Flags().GetString("path")
		spell, err := m.bareRun(urlSpell{
			URL:  args[0],
			Path: destPath,
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

func (m *urlModule) cobraArgs() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s module requires at least one argument.", urlModuleName)
		}
		p, _ := cmd.Flags().GetString("path")
		if p == "" {
			return fmt.Errorf("%s module requires --path to be set.", urlModuleName)
		}
		return nil
	}
}

// CliConfig implements Module.
func (m *urlModule) CliConfig(config *Config) *cobra.Command {
	command := &cobra.Command{
		Args:    m.cobraArgs(),
		Example: urlModuleExample,
		Run:     m.cobraRun(config),
		Short:   urlModuleShort,
		Use:     urlModuleName,
	}
	command.PersistentFlags().String("path", "", "Destination path for the downloaded file or extraction directory.")
	return command
}

// --- helpers ---

// isArchive reports whether the URL points to a supported archive format.
func isArchive(rawURL string) bool {
	lower := strings.ToLower(path.Base(rawURL))
	return strings.HasSuffix(lower, ".zip") ||
		strings.HasSuffix(lower, ".tar.gz") ||
		strings.HasSuffix(lower, ".tgz")
}

// downloadFile fetches url and writes its body to dest.
func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("url: download %s: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("url: download %s returned HTTP %d", url, resp.StatusCode)
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// copyExecutable copies src to dest (creating dest if it is a directory path)
// and sets executable bits.
func copyExecutable(src, dest string) error {
	// If dest is an existing directory, place the file inside it.
	if info, err := os.Stat(dest); err == nil && info.IsDir() {
		dest = filepath.Join(dest, filepath.Base(src))
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return os.Chmod(dest, 0755)
}

// extractArchive dispatches to the appropriate extractor based on the filename.
func extractArchive(src, dest string) error {
	lower := strings.ToLower(src)
	if strings.HasSuffix(lower, ".zip") {
		return extractZip(src, dest)
	}
	return extractTarGz(src, dest)
}

// extractZip extracts a zip archive into dest.
func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(dest, DirectoryPermissions); err != nil {
		return err
	}
	cleanDest := filepath.Clean(dest) + string(os.PathSeparator)

	for _, f := range r.File {
		outPath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(outPath, cleanDest) {
			return fmt.Errorf("url: unsafe path in zip archive: %s", f.Name)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(outPath, DirectoryPermissions); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(outPath), DirectoryPermissions); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.Create(outPath)
		if err != nil {
			rc.Close()
			return err
		}
		_, copyErr := io.Copy(out, rc)
		out.Close()
		rc.Close()
		if copyErr != nil {
			return copyErr
		}
	}
	return nil
}

// extractTarGz extracts a .tar.gz or .tgz archive into dest.
func extractTarGz(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	if err := os.MkdirAll(dest, DirectoryPermissions); err != nil {
		return err
	}
	cleanDest := filepath.Clean(dest) + string(os.PathSeparator)

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		outPath := filepath.Join(dest, header.Name)
		if !strings.HasPrefix(outPath, cleanDest) {
			return fmt.Errorf("url: unsafe path in tar archive: %s", header.Name)
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(outPath, DirectoryPermissions); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(outPath), DirectoryPermissions); err != nil {
				return err
			}
			out, err := os.Create(outPath)
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(out, tr)
			out.Close()
			if copyErr != nil {
				return copyErr
			}
		}
	}
	return nil
}
