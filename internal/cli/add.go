package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new SSH host",
	Long:  `Add a new SSH host to the database with interactive prompts.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)

		var name string
		if len(args) > 0 {
			name = args[0]
		} else {
			name = promptString(scanner, "Host name")
		}

		hostname := promptString(scanner, "Hostname/IP")
		username := promptString(scanner, "Username")

		portStr := promptStringWithDefault(scanner, "Port", "22")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Printf("Invalid port number: %s\n", portStr)
			return
		}

		keyPath := promptStringWithDefault(scanner, "SSH key path (optional)", "")
		description := promptStringWithDefault(scanner, "Description (optional)", "")
		tags := promptStringWithDefault(scanner, "Tags (comma-separated, optional)", "")

		host, err := hostService.CreateHost(name, hostname, username, port, keyPath, description, tags)
		if err != nil {
			fmt.Printf("Error creating host: %v\n", err)
			return
		}

		fmt.Printf("✅ Host '%s' added successfully (ID: %d)\n", host.Name, host.ID)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func promptString(scanner *bufio.Scanner, prompt string) string {
	for {
		fmt.Printf("%s: ", prompt)
		scanner.Scan()
		value := strings.TrimSpace(scanner.Text())
		if value != "" {
			return value
		}
		fmt.Printf("Please enter a value for %s\n", prompt)
	}
}

func promptStringWithDefault(scanner *bufio.Scanner, prompt, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	scanner.Scan()
	value := strings.TrimSpace(scanner.Text())
	if value == "" {
		return defaultValue
	}
	return value
}
