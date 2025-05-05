package use

import (
	"github.com/spf13/cobra"
)

// NewUseCmd creates the use command
func NewUseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use",
		Short: "Set context for NStream AI resources",
		Long:  `Set the current context for various NStream AI resources like clusters, buckets, etc.`,
	}

	// Add subcommands
	cmd.AddCommand(
		NewClusterCmd(),
		NewBucketCmd(),
	)

	return cmd
}
