package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// writeConfig writes a config file into a temporary directory
// and returns its path.
func writeConfig(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("writing the test config: %v", err)
	}

	return path
}

func TestLoadMissingFileKeepsDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.yaml")

	config, err := Load(path)

	if !os.IsNotExist(err) {
		t.Fatalf("Load = %v, want a not-exist error", err)
	}

	if !reflect.DeepEqual(config, Default()) {
		t.Error("a missing config file did not fall back to the defaults")
	}
}

func TestLoadLayersFileOverDefaults(t *testing.T) {
	path := writeConfig(t, "board_rows: 20\nmine_count: 30\nflag_char: \"⚑\"\n")

	config, err := Load(path)
	if err != nil {
		t.Fatalf("Load = %v, want nil", err)
	}

	if got, want := config.BoardRows, 20; got != want {
		t.Errorf("BoardRows = %d, want %d", got, want)
	}

	if got, want := config.MineCount, 30; got != want {
		t.Errorf("MineCount = %d, want %d", got, want)
	}

	if got, want := config.FlagChar, "⚑"; got != want {
		t.Errorf("FlagChar = %q, want %q", got, want)
	}

	// Untouched settings keep their default.
	if got, want := config.BoardCols, Default().BoardCols; got != want {
		t.Errorf("BoardCols = %d, want the default %d", got, want)
	}
}

func TestLoadRejectsUnusableConfigs(t *testing.T) {
	tests := []struct {
		name     string
		contents string
	}{
		{"malformed yaml", "board_rows: [1, 2\n"},
		{"wrong type", "board_rows: \"twenty\"\n"},
		{"negative scrolloff", "scrolloff: -1\n"},
		{"negative sidescrolloff", "sidescrolloff: -1\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, err := Load(writeConfig(t, test.contents))

			if err == nil {
				t.Fatal("Load = nil, want an error")
			}

			if !reflect.DeepEqual(config, Default()) {
				t.Error("an invalid config was not replaced by the defaults")
			}
		})
	}
}

func TestDefaultConfigValidates(t *testing.T) {
	if err := Default().validate(); err != nil {
		t.Fatalf("the built-in defaults do not validate: %v", err)
	}
}

func TestDefaultPathIsAbsolute(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home directory in this environment")
	}

	if got, want := DefaultPath(), filepath.Join(home, ".config", "termsweeper", "config.yaml"); got != want {
		t.Errorf("DefaultPath = %q, want %q", got, want)
	}
}
