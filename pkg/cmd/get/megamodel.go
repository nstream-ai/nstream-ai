package get

import (
	"github.com/spf13/cobra"
)

var (
	megaModelName         string
	megaModelOutputFormat string
)

// NewMegaModelCmd creates the get megamodel command
func NewMegaModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "megamodel",
		Short: "Get mega model information",
		Long:  `Get detailed information about a specific mega model or list all mega models`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement mega model get logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&megaModelName, "name", "n", "", "Mega model name (optional, lists all mega models if not specified)")
	cmd.Flags().StringVarP(&megaModelOutputFormat, "output", "o", "table", "Output format (table, json, yaml)")

	return cmd
}
