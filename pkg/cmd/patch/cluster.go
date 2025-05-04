package patch

import (
	"github.com/spf13/cobra"
)

var (
	clusterName string
	clusterSpec string
)

// NewClusterCmd creates the patch cluster command
func NewClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Patch a cluster",
		Long:  `Patch a specific cluster in the NStream AI platform`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement cluster patching logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&clusterName, "name", "n", "", "Cluster name")
	cmd.Flags().StringVarP(&clusterSpec, "spec", "s", "", "Cluster specification in JSON format")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("spec")

	return cmd
}
