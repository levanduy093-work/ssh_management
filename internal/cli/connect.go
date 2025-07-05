package cli

import (
	"fmt"
	"strconv"

	"github.com/levanduy/ssh_management/internal/domain"
	"github.com/levanduy/ssh_management/pkg/ssh"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:     "connect [name_or_id]",
	Aliases: []string{"c", "ssh"},
	Short:   "Connect to an SSH host",
	Long:    `Connect to an SSH host by name or ID. This will start an SSH session.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]

		// Try to parse as ID first
		if id, err := strconv.Atoi(identifier); err == nil {
			connectByID(id)
		} else {
			connectByName(identifier)
		}
	},
}

var showCmd = &cobra.Command{
	Use:   "show [name_or_id]",
	Short: "Show detailed information about a host",
	Long:  `Display detailed information about an SSH host including connection command.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]

		// Try to parse as ID first
		if id, err := strconv.Atoi(identifier); err == nil {
			showHostByID(id)
		} else {
			showHostByName(identifier)
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(showCmd)
}

func connectByID(id int) {
	host, err := hostService.GetHostByID(id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	connectToHost(host)
}

func connectByName(name string) {
	host, err := hostService.GetHostByName(name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	connectToHost(host)
}

func connectToHost(host *domain.Host) {
	fmt.Printf("🔗 Connecting to %s (%s@%s:%d)...\n", host.Name, host.Username, host.Hostname, host.Port)

	// Increment use count
	if err := hostService.ConnectToHost(host.ID); err != nil {
		fmt.Printf("Warning: Failed to update usage statistics: %v\n", err)
	}

	// Connect via SSH
	if err := ssh.ConnectToHost(host); err != nil {
		fmt.Printf("SSH connection failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Connection to %s closed.\n", host.Name)
}

func showHostByID(id int) {
	host, err := hostService.GetHostByID(id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	showHost(host)
}

func showHostByName(name string) {
	host, err := hostService.GetHostByName(name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	showHost(host)
}

func showHost(host *domain.Host) {
	fmt.Printf("📋 Host Details\n")
	fmt.Printf("═══════════════\n")
	fmt.Printf("ID:          %d\n", host.ID)
	fmt.Printf("Name:        %s\n", host.Name)
	fmt.Printf("Hostname:    %s\n", host.Hostname)
	fmt.Printf("Username:    %s\n", host.Username)
	fmt.Printf("Port:        %d\n", host.Port)

	if host.KeyPath != "" {
		fmt.Printf("SSH Key:     %s\n", host.KeyPath)
	}

	if host.Description != "" {
		fmt.Printf("Description: %s\n", host.Description)
	}

	if host.Tags != "" {
		fmt.Printf("Tags:        %s\n", host.Tags)
	}

	fmt.Printf("Use Count:   %d\n", host.UseCount)

	if !host.LastUsed.IsZero() {
		fmt.Printf("Last Used:   %s\n", formatTime(host.LastUsed))
	}

	fmt.Printf("Created:     %s\n", host.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated:     %s\n", host.UpdatedAt.Format("2006-01-02 15:04:05"))

	fmt.Printf("\n💻 SSH Command\n")
	fmt.Printf("══════════════\n")
	fmt.Printf("%s\n", ssh.BuildSSHCommand(host))
}
