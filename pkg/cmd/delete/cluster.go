package delete

import (
	"github.com/spf13/cobra"
)

var (
	clusterName string
	force       bool
)

// NewClusterCmd creates the delete cluster command
func NewClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Delete a cluster",
		Long:  `Delete a specific cluster from the NStream AI platform`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement cluster deletion logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&clusterName, "name", "n", "", "Cluster name")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force deletion without confirmation")

	cmd.MarkFlagRequired("name")

	return cmd
}
