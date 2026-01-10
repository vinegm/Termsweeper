package main

import (
	"strings"
)

// Renders the menu screen.
func renderMenu(model model) string {
	var sb strings.Builder

	sb.WriteString("Termsweeper\n")

	items := []string{"Start Game", "Quit"}
	for i, it := range items {
		sb.WriteString("\n")
		if model.menuChoice == i {
			sb.WriteString(selectedMenuStyle.Render(it))
			continue
		}

		sb.WriteString(it)
	}

	return sb.String()
}
