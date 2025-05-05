package create

import (
	"github.com/spf13/cobra"
)

// NewCreateCmd creates the create command
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create NStream AI resources",
		Long:  `Create various NStream AI resources like clusters, principals, etc.`,
	}

	// Add subcommands
	cmd.AddCommand(
		NewClusterCmd(),
		NewBucketCmd(),
	)

	return cmd
}
