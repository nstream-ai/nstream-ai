package use

import (
	"github.com/spf13/cobra"
)

// NewUseCmd creates the root use command
func NewUseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use",
		Short: "Set context for NStream AI platform resources",
		Long:  `Set the current context for various NStream AI platform resources including clusters, buckets, and roles`,
	}

	// Add all subcommands
	cmd.AddCommand(
		NewClusterCmd(),
		NewBucketCmd(),
	)

	return cmd
}
