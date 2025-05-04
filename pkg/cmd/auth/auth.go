package auth

import (
	"github.com/spf13/cobra"
)

// RegisterCommands registers all auth-related commands
func RegisterCommands(rootCmd *cobra.Command) {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication commands",
		Long:  `Commands for authentication and authorization`,
	}

	authCmd.AddCommand(NewSigninCmd())
	authCmd.AddCommand(NewSignupCmd())

	rootCmd.AddCommand(authCmd)
}
