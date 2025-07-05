package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "v1.0.0" // Will be set during build

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display the current version of SSH Manager.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("SSH Manager (sshm) %s\n", version)
		fmt.Println("🚀 Easy SSH host management for terminal")
		fmt.Println("")
		fmt.Println("Built with ❤️  using Go, Cobra, and Bubble Tea")
		fmt.Println("Repository: https://github.com/levanduy/ssh_management")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
