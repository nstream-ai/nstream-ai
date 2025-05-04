package get

import (
	"github.com/spf13/cobra"
)

var (
	connectorName         string
	connectorOutputFormat string
)

// NewStreamConnectorCmd creates the get streamconnector command
func NewStreamConnectorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "streamconnector",
		Short: "Get stream connector information",
		Long:  `Get detailed information about a specific stream connector or list all stream connectors`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement stream connector get logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&connectorName, "name", "n", "", "Stream connector name (optional, lists all stream connectors if not specified)")
	cmd.Flags().StringVarP(&connectorOutputFormat, "output", "o", "table", "Output format (table, json, yaml)")

	return cmd
}
