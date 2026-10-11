package paths

import (
	"os"
	"path/filepath"
)

const appDir = "termsweeper"

// Config is $XDG_CONFIG_HOME/termsweeper/<parts...>
// defaulting to ~/.config.
func Config(parts ...string) string {
	return under(base("XDG_CONFIG_HOME", ".config"), parts)
}

func base(env string, fallback ...string) string {
	dir := os.Getenv(env)
	if filepath.IsAbs(dir) {
		return dir
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(append([]string{home}, fallback...)...)
}

func under(base string, parts []string) string {
	if base == "" {
		return filepath.Join(parts...)
	}

	return filepath.Join(append([]string{base, appDir}, parts...)...)
}
