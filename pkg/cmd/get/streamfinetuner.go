package get

import (
	"github.com/spf13/cobra"
)

var (
	finetunerName         string
	finetunerOutputFormat string
)

// NewStreamFinetunerCmd creates the get streamfinetuner command
func NewStreamFinetunerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "streamfinetuner",
		Short: "Get stream finetuner information",
		Long:  `Get detailed information about a specific stream finetuner or list all stream finetuners`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement stream finetuner get logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&finetunerName, "name", "n", "", "Stream finetuner name (optional, lists all stream finetuners if not specified)")
	cmd.Flags().StringVarP(&finetunerOutputFormat, "output", "o", "table", "Output format (table, json, yaml)")

	return cmd
}
