package main

import (
	"testing"

	"termsweeper/internal/config"
	"termsweeper/internal/minesweeper"
)

func TestGameSettingsFromConfig(t *testing.T) {
	cfg := config.Default()
	cfg.BoardRows, cfg.BoardCols, cfg.MineCount = 20, 30, 50

	want := minesweeper.Settings{Rows: 20, Cols: 30, Mines: 50}

	if got := gameSettings(cfg); got != want {
		t.Errorf("gameSettings = %+v, want %+v", got, want)
	}
}

func TestGameSettingsFallsBackToDefaults(t *testing.T) {
	defaults := settingsFrom(config.Default())

	tests := []struct {
		name              string
		rows, cols, mines int
	}{
		{"board too small", 1, 16, 10},
		{"board too big", 16, 100_000, 10},
		{"no mines", 16, 16, 0},
		{"more mines than cells", 8, 8, 64},
		{"negative values", -5, -5, -5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := config.Default()
			cfg.BoardRows, cfg.BoardCols, cfg.MineCount = test.rows, test.cols, test.mines

			if got := gameSettings(cfg); got != defaults {
				t.Errorf("gameSettings = %+v, want the defaults %+v", got, defaults)
			}
		})
	}
}

func TestDefaultConfigIsPlayable(t *testing.T) {
	if err := settingsFrom(config.Default()).Validate(); err != nil {
		t.Fatalf("the built-in board does not validate: %v", err)
	}
}
