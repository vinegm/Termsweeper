package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	logFile  = "bubbletea.log"
	logLevel = "debug"
)

// Tea initialization
func (model model) Init() tea.Cmd { return nil }

func main() {
	file, err := tea.LogToFile(logFile, logLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	program := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err = program.Run()
	if err != nil {
		log.Fatal(err)
	}
}
