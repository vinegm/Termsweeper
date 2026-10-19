package main

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"

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

type flags struct {
	configPath  string
	debugFile   string
	debug       bool
	showVersion bool
	noConfig    bool
}

func main() {
	opts := parseFlags()

	if opts.showVersion {
		fmt.Println("termsweeper version", version)
		return
	}

	log.SetOutput(io.Discard)
	if opts.debug {
		file, err := tea.LogToFile(opts.debugFile, logLevel)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	cfg := loadConfig(opts.noConfig, opts.configPath)

	program := tea.NewProgram(tui.New(cfg, gameSettings(cfg)))
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}

func parseFlags() flags {
	var opts flags

	flag.BoolVar(&opts.debug, "debug", false, "enable debug mode")
	flag.StringVar(&opts.debugFile, "debug-file", logFile, "path to debug log file (ignored if --debug is not set)")
	flag.BoolVar(&opts.showVersion, "version", false, "print the version and exit")
	flag.BoolVar(&opts.noConfig, "no-config", false, "ignore config.yaml and run with the built-in defaults")
	flag.StringVar(&opts.configPath, "config", "", "config file to load instead of the default path")
	flag.Parse()

	return opts
}

// loadConfig returns the defaults with -no-config, else the file at
// path or the default path. Only an explicit path must exist.
func loadConfig(noConfig bool, path string) config.Config {
	if noConfig {
		return config.Default()
	}

	cfg, err := config.Load(cmp.Or(path, config.DefaultPath()))
	if err != nil && (!errors.Is(err, fs.ErrNotExist) || path != "") {
		warn("could not load config: %v", err)
	}

	return cfg
}

// gameSettings is the board a new game starts from.
func gameSettings(cfg config.Config) minesweeper.Settings {
	settings := settingsFrom(cfg)

	if err := settings.Validate(); err != nil {
		warn("invalid board settings, using the defaults: %v", err)

		return settingsFrom(config.Default())
	}

	return settings
}

// warn reports a problem found before the TUI starts, on
// stderr so it is seen without -debug, and in the debug log.
func warn(format string, args ...any) {
	text := fmt.Sprintf(format, args...)

	fmt.Fprintln(os.Stderr, "warning: "+text)
	log.Print(text)
}

func settingsFrom(cfg config.Config) minesweeper.Settings {
	return minesweeper.Settings{
		Rows:  cfg.Board.Rows,
		Cols:  cfg.Board.Cols,
		Mines: cfg.Board.Mines,
	}
}
