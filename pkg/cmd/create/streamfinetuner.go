package create

import (
	"github.com/spf13/cobra"
)

var (
	finetunerName       string
	finetunerType       string
	finetunerConfigPath string
)

// NewStreamFinetunerCmd creates the streamfinetuner command
func NewStreamFinetunerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "streamfinetuner",
		Short: "Create a new stream finetuner",
		Long:  `Create a new stream finetuner with specified configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement stream finetuner creation logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&finetunerName, "name", "n", "", "Finetuner name")
	cmd.Flags().StringVarP(&finetunerType, "type", "t", "", "Finetuner type")
	cmd.Flags().StringVarP(&finetunerConfigPath, "config", "c", "", "Path to finetuner configuration file")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("config")

	return cmd
}
