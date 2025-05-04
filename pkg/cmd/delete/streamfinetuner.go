package delete

import (
	"github.com/spf13/cobra"
)

var (
	finetunerName        string
	finetunerForceDelete bool
)

// NewStreamFinetunerCmd creates the delete streamfinetuner command
func NewStreamFinetunerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "streamfinetuner",
		Short: "Delete a stream finetuner",
		Long:  `Delete a specific stream finetuner from the NStream AI platform`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement stream finetuner deletion logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&finetunerName, "name", "n", "", "Stream finetuner name")
	cmd.Flags().BoolVarP(&finetunerForceDelete, "force", "f", false, "Force deletion without confirmation")

	cmd.MarkFlagRequired("name")

	return cmd
}
