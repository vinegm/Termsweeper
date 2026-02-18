package main

import (
	"flag"
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"termsweeper/internal/config"
	"termsweeper/internal/tui"
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

	cfg, err := config.Load(config.DefaultPath())
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Could not load config: %v", err)
	}

	program := tea.NewProgram(tui.New(cfg), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
