package get

import (
	"github.com/spf13/cobra"
)

var (
	baseModelName         string
	baseModelOutputFormat string
)

// NewBaseModelCmd creates the get basemodel command
func NewBaseModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "basemodel",
		Short: "Get base model information",
		Long:  `Get detailed information about a specific base model or list all base models`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement base model get logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&baseModelName, "name", "n", "", "Base model name (optional, lists all base models if not specified)")
	cmd.Flags().StringVarP(&baseModelOutputFormat, "output", "o", "table", "Output format (table, json, yaml)")

	return cmd
}
