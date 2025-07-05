package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List all SSH hosts",
	Long:    `List all SSH hosts in the database with their details.`,
	Run: func(cmd *cobra.Command, args []string) {
		hosts, err := hostService.GetAllHosts()
		if err != nil {
			fmt.Printf("Error retrieving hosts: %v\n", err)
			return
		}

		if len(hosts) == 0 {
			fmt.Println("No hosts found. Use 'sshm add' to add a new host.")
			return
		}

		// Print header
		fmt.Printf("%-4s %-20s %-25s %-15s %-10s %-20s\n",
			"ID", "NAME", "HOST", "USER", "PORT", "LAST USED")
		fmt.Println(strings.Repeat("-", 100))

		// Print hosts
		for _, host := range hosts {
			lastUsed := "Never"
			if !host.LastUsed.IsZero() {
				lastUsed = formatTime(host.LastUsed)
			}

			fmt.Printf("%-4d %-20s %-25s %-15s %-10d %-20s\n",
				host.ID,
				truncate(host.Name, 20),
				truncate(host.Hostname, 25),
				truncate(host.Username, 15),
				host.Port,
				lastUsed)

			// Show description and tags if present
			if host.Description != "" || host.Tags != "" {
				fmt.Printf("     ")
				if host.Description != "" {
					fmt.Printf("📝 %s", truncate(host.Description, 50))
				}
				if host.Tags != "" {
					if host.Description != "" {
						fmt.Printf(" | ")
					}
					fmt.Printf("🏷️  %s", host.Tags)
				}
				fmt.Println()
			}
		}

		fmt.Printf("\nTotal: %d host(s)\n", len(hosts))
	},
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search hosts by name, hostname, description, or tags",
	Long:  `Search for SSH hosts using a query string that matches name, hostname, description, or tags.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		hosts, err := hostService.SearchHosts(query)
		if err != nil {
			fmt.Printf("Error searching hosts: %v\n", err)
			return
		}

		if len(hosts) == 0 {
			fmt.Printf("No hosts found matching '%s'\n", query)
			return
		}

		fmt.Printf("Found %d host(s) matching '%s':\n\n", len(hosts), query)

		// Print header
		fmt.Printf("%-4s %-20s %-25s %-15s %-10s %-20s\n",
			"ID", "NAME", "HOST", "USER", "PORT", "LAST USED")
		fmt.Println(strings.Repeat("-", 100))

		// Print hosts
		for _, host := range hosts {
			lastUsed := "Never"
			if !host.LastUsed.IsZero() {
				lastUsed = formatTime(host.LastUsed)
			}

			fmt.Printf("%-4d %-20s %-25s %-15s %-10d %-20s\n",
				host.ID,
				truncate(host.Name, 20),
				truncate(host.Hostname, 25),
				truncate(host.Username, 15),
				host.Port,
				lastUsed)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Hour {
		return fmt.Sprintf("%dm ago", int(diff.Minutes()))
	}
	if diff < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	}
	if diff < 7*24*time.Hour {
		return fmt.Sprintf("%dd ago", int(diff.Hours()/24))
	}
	return t.Format("2006-01-02")
}
