package main

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"

	tea "charm.land/bubbletea/v2"

	"github.com/vinegm/termsweeper/internal/config"
	"github.com/vinegm/termsweeper/internal/minesweeper"
	"github.com/vinegm/termsweeper/internal/tui"
)

const (
	logFile  = "/tmp/termsweeper.log"
	logLevel = "debug"
)

var version = "dev"

func main() {
	debug := flag.Bool("debug", false, "enable debug mode")
	debugFile := flag.String("debug-file", logFile, "path to debug log file (ignored if --debug is not set)")
	showVersion := flag.Bool("version", false, "print the version and exit")
	noConfig := flag.Bool("no-config", false, "ignore config.yaml and run with the built-in defaults")
	configPath := flag.String("config", "", "config file to load instead of the default path")

	flag.Parse()

	if *showVersion {
		fmt.Println("termsweeper version", version)
		return
	}

	log.SetOutput(io.Discard)
	if *debug {
		file, err := tea.LogToFile(*debugFile, logLevel)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	cfg := loadConfig(*noConfig, *configPath)

	program := tea.NewProgram(tui.New(cfg, gameSettings(cfg)))
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}

// loadConfig returns the defaults with -no-config, else the file at
// path or the default path. Only an explicit path must exist.
func loadConfig(noConfig bool, path string) config.Config {
	if noConfig {
		return config.Default()
	}

	cfg, err := config.Load(cmp.Or(path, config.DefaultPath()))
	if err != nil && (!errors.Is(err, fs.ErrNotExist) || path != "") {
		log.Printf("Could not load config: %v", err)
	}

	return cfg
}

// gameSettings is the board a new game starts from.
func gameSettings(cfg config.Config) minesweeper.Settings {
	settings := settingsFrom(cfg)

	if err := settings.Validate(); err != nil {
		log.Printf("Invalid board settings, using the defaults: %v", err)

		return settingsFrom(config.Default())
	}

	return settings
}

func settingsFrom(cfg config.Config) minesweeper.Settings {
	return minesweeper.Settings{
		Rows:  cfg.Board.Rows,
		Cols:  cfg.Board.Cols,
		Mines: cfg.Board.Mines,
	}
}
