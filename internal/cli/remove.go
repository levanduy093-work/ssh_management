package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove [name_or_id]",
	Aliases: []string{"rm", "delete", "del"},
	Short:   "Remove an SSH host",
	Long:    `Remove an SSH host from the database. This action cannot be undone.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]

		// Try to parse as ID first
		if id, err := strconv.Atoi(identifier); err == nil {
			removeByID(id)
		} else {
			removeByName(identifier)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func removeByID(id int) {
	host, err := hostService.GetHostByID(id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	removeHost(host.ID, host.Name)
}

func removeByName(name string) {
	host, err := hostService.GetHostByName(name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	removeHost(host.ID, host.Name)
}

func removeHost(id int, name string) {
	fmt.Printf("⚠️  Are you sure you want to remove host '%s' (ID: %d)? [y/N]: ", name, id)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := strings.ToLower(strings.TrimSpace(scanner.Text()))

	if response != "y" && response != "yes" {
		fmt.Println("❌ Remove operation cancelled.")
		return
	}

	if err := hostService.DeleteHost(id); err != nil {
		fmt.Printf("Error removing host: %v\n", err)
		return
	}

	fmt.Printf("✅ Host '%s' has been removed successfully.\n", name)
}
