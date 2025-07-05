package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/levanduy/ssh_management/internal/ui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:     "tui",
	Aliases: []string{"ui", "interactive"},
	Short:   "Launch interactive TUI mode",
	Long:    `Launch the interactive Terminal User Interface for browsing and managing SSH hosts.`,
	Run: func(cmd *cobra.Command, args []string) {
		model := ui.NewModel(hostService)
		p := tea.NewProgram(model, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
