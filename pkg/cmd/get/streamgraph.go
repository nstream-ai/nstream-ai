package get

import (
	"github.com/spf13/cobra"
)

var (
	graphName         string
	graphOutputFormat string
)

// NewStreamGraphCmd creates the get streamgraph command
func NewStreamGraphCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "streamgraph",
		Short: "Get stream graph information",
		Long:  `Get detailed information about a specific stream graph or list all stream graphs`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement stream graph get logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&graphName, "name", "n", "", "Stream graph name (optional, lists all stream graphs if not specified)")
	cmd.Flags().StringVarP(&graphOutputFormat, "output", "o", "table", "Output format (table, json, yaml)")

	return cmd
}
