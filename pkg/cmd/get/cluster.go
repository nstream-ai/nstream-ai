package get

import (
	"github.com/spf13/cobra"
)

var (
	clusterName  string
	outputFormat string
)

// NewClusterCmd creates the get cluster command
func NewClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Get cluster information",
		Long:  `Get detailed information about a specific cluster or list all clusters`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement cluster get logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&clusterName, "name", "n", "", "Cluster name (optional, lists all clusters if not specified)")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, yaml)")

	return cmd
}
