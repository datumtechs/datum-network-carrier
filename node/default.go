package node

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// DefaultConfig contains reasonable default settings.
var DefaultConfig = Config{
	DataDir: DefaultDataDir(),
	Name:    "carrier",
}

// DefaultDataDir is the default data directory to use for the databases and other
// persistence requirements.
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Carrier")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "Carrier")
		} else {
			return filepath.Join(home, ".carrier")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
