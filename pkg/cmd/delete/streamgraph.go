package delete

import (
	"github.com/spf13/cobra"
)

var (
	graphName        string
	graphForceDelete bool
)

// NewStreamGraphCmd creates the delete streamgraph command
func NewStreamGraphCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "streamgraph",
		Short: "Delete a stream graph",
		Long:  `Delete a specific stream graph from the NStream AI platform`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement stream graph deletion logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&graphName, "name", "n", "", "Stream graph name")
	cmd.Flags().BoolVarP(&graphForceDelete, "force", "f", false, "Force deletion without confirmation")

	cmd.MarkFlagRequired("name")

	return cmd
}
