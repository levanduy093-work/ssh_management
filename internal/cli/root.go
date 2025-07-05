package cli

import (
	"fmt"
	"os"

	"github.com/levanduy/ssh_management/internal/repo"
	"github.com/levanduy/ssh_management/internal/service"
	"github.com/spf13/cobra"
)

var (
	dbPath      string
	hostService *service.HostService
)

var rootCmd = &cobra.Command{
	Use:   "sshm",
	Short: "SSH Manager - Easy SSH host management",
	Long: `SSH Manager (sshm) is a terminal-based tool for managing SSH connections.
It allows you to store, organize, and quickly connect to SSH hosts.

Features:
- Add, edit, and remove SSH hosts
- Interactive TUI for browsing hosts
- Search and filter hosts by name, hostname, or tags
- Quick connection with usage tracking`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initializeService()
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Default to TUI mode when no command specified
		tuiCmd.Run(cmd, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", service.GetDefaultDatabasePath(), "Database file path")
}

func initializeService() {
	repo, err := repo.NewSQLiteRepo(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	hostService = service.NewHostService(repo)
}
