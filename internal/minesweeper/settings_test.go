package minesweeper

import "testing"

func TestSettingsValidate(t *testing.T) {
	tests := []struct {
		name     string
		settings Settings
		wantErr  bool
	}{
		{"defaults", Settings{Rows: 16, Cols: 16, Mines: 40}, false},
		{"smallest board", Settings{Rows: MinDim, Cols: MinDim, Mines: MinMines}, false},
		{"full board", Settings{Rows: 5, Cols: 5, Mines: 24}, false},
		{"rows below minimum", Settings{Rows: MinDim - 1, Cols: 16, Mines: 10}, true},
		{"cols below minimum", Settings{Rows: 16, Cols: MinDim - 1, Mines: 10}, true},
		{"rows above maximum", Settings{Rows: MaxDim + 1, Cols: 16, Mines: 10}, true},
		{"cols above maximum", Settings{Rows: 16, Cols: MaxDim + 1, Mines: 10}, true},
		{"no mines", Settings{Rows: 16, Cols: 16, Mines: 0}, true},
		{"no room for the first move", Settings{Rows: 5, Cols: 5, Mines: 25}, true},
		{"more mines than cells", Settings{Rows: 16, Cols: 16, Mines: 1000}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.settings.Validate()

			if test.wantErr && err == nil {
				t.Fatalf("Validate(%+v) = nil, want an error", test.settings)
			}

			if !test.wantErr && err != nil {
				t.Fatalf("Validate(%+v) = %v, want nil", test.settings, err)
			}
		})
	}
}
