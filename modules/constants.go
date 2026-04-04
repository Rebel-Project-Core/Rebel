package modules

// DirectoryPermissions is the standard permission bits for module download directories.
const DirectoryPermissions = 0755

// PipVenvDirectoryName is the subdirectory name for the Python virtual environment.
const PipVenvDirectoryName = "venv"

// CranLibraryDirectoryName is the R library directory name.
const CranLibraryDirectoryName = "R-Library"

// PipAptDependencies lists the apt packages required before pip can run.
var PipAptDependencies = []string{"python3", "python3-pip", "python3-venv"}

// CranAptDependencies lists the apt packages required before CRAN can run.
var CranAptDependencies = []string{"r-base", "r-base-dev"}

// Cache key constants for the cran module.
const (
	cranCacheKeyApply = cranModuleName + "apply"
	cranCacheKeySave  = cranModuleName + "save"
	cranCacheKeyBare  = cranModuleName + "bare"
	cranCacheKeyBioc  = cranModuleName + "bioc"
)
