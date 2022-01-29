package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"

	"github.com/vikin91/wurdle/pkg/model"
	"github.com/vikin91/wurdle/pkg/ui"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	var (
		daemonMode bool
		showHelp   bool
		opts       []tea.ProgramOption
	)

	flag.BoolVar(&daemonMode, "d", false, "run as a daemon")
	flag.BoolVar(&showHelp, "h", false, "show help")
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if daemonMode || !isatty.IsTerminal(os.Stdout.Fd()) {
		// If we're in daemon mode don't render the TUI
		opts = []tea.ProgramOption{tea.WithoutRenderer()}
	} else {
		// If we're in TUI mode, discard log output
		log.SetOutput(ioutil.Discard)
		opts = []tea.ProgramOption{
			tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
			tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
		}
	}

	help := help.New()
	state := &State{
		model: model.NewModel(6),
		ui:    ui.NewUI(),
		help:  &help,
	}
	state.model.AddGames("en", "de", "pl")

	p := tea.NewProgram(state,

		opts...)
	if err := p.Start(); err != nil {
		fmt.Println("Error starting Bubble Tea program:", err)
		os.Exit(1)
	}
}
