package main

import (
	"flag"
	"io"
	"log"
	"os"

	"termsweeper/src"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	logFile  = "/tmp/termsweeper.log"
	logLevel = "debug"
)

func main() {
	debug := flag.Bool("debug", false, "enable debug mode")
	debugFile := flag.String("debug-file", logFile, "path to debug log file (ignored if --debug is not set)")

	flag.Parse()

	log.SetOutput(io.Discard)
	if *debug {
		file, err := tea.LogToFile(*debugFile, logLevel)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	if err := src.LoadConfig(); err != nil && !os.IsNotExist(err) {
		log.Printf("Could not load config: %v", err)
	}

	program := tea.NewProgram(src.InitialModel(), tea.WithAltScreen())
	_, err := program.Run()
	if err != nil {
		log.Fatal(err)
	}
}
