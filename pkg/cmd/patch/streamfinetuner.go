package patch

import (
	"github.com/spf13/cobra"
)

var (
	finetunerName string
	finetunerSpec string
)

// NewStreamFinetunerCmd creates the patch streamfinetuner command
func NewStreamFinetunerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "streamfinetuner",
		Short: "Patch a stream finetuner",
		Long:  `Patch a specific stream finetuner in the NStream AI platform`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement stream finetuner patching logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&finetunerName, "name", "n", "", "Stream finetuner name")
	cmd.Flags().StringVarP(&finetunerSpec, "spec", "s", "", "Stream finetuner specification in JSON format")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("spec")

	return cmd
}
