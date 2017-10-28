package appconfig

import (
	"os"
	"path/filepath"
)

// GetAppPath returns the absolute application root path
func GetAppPath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return BuildPath(false, dir, "..")
}

// BuildPath builds a path from an array of path parts (used to include proper PathSeparator)
func BuildPath(isFile bool, parts ...string) string {
	result := ""
	for _, part := range parts {
		if part == "" {
			continue
		}

		result += part

		// Add a `os.PathSeparator` after if it's missing
		if len(result) >= 1 && string(result[len(result)-1]) != string(os.PathSeparator) {
			result += string(os.PathSeparator)
		}
	}

	if len(result) > 0 && isFile && string(result[len(result)-1]) == string(os.PathSeparator) {
		result = result[0 : len(result)-1]
	}

	return result
}
