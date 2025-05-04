package create

import (
	"github.com/spf13/cobra"
)

var (
	connectorName string
	connectorType string
	configPath    string
)

// NewStreamConnectorCmd creates the streamconnector command
func NewStreamConnectorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "streamconnector",
		Short: "Create a new stream connector",
		Long:  `Create a new stream connector with specified configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement stream connector creation logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&connectorName, "name", "n", "", "Connector name")
	cmd.Flags().StringVarP(&connectorType, "type", "t", "", "Connector type")
	cmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to connector configuration file")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("config")

	return cmd
}
