package main

import (
	"log"

	"termsweeper/src"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	logFile  = "bubbletea.log"
	logLevel = "debug"
)

func main() {
	file, err := tea.LogToFile(logFile, logLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	program := tea.NewProgram(src.InitialModel(), tea.WithAltScreen())
	_, err = program.Run()
	if err != nil {
		log.Fatal(err)
	}
}
