package project

import (
	"os"
	"path"
)

// ProjectDirectoryName is the name of the local project directory.
const ProjectDirectoryName = "credoenv"

// Global project path cache.
var gPath *string

// ProjectPath returns the project directory path, creating it if necessary.
// The returned string is a copy and safe to use after subsequent calls.
func ProjectPath() (string, error) {
	if gPath != nil {
		return *gPath, nil
	}
	basePath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	projectPath := path.Join(basePath, ProjectDirectoryName)

	// Create project path.
	err = os.MkdirAll(projectPath, 0755)
	if err != nil {
		return "", err
	}
	gPath = &projectPath

	return projectPath, nil
}
