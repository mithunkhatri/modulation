package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mithunkhatri/modulation/internal/audio"
	"github.com/mithunkhatri/modulation/internal/tui"
)

var version = "dev"

func main() {
	player, err := audio.NewPlayer()
	if err != nil {
		fmt.Printf("Error creating player: %v\n", err)
		os.Exit(1)
	}

	m := tui.NewModel(player)
	tui.Version = version

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
