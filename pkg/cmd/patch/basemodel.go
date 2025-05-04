package patch

import (
	"github.com/spf13/cobra"
)

var (
	baseModelName string
	baseModelSpec string
)

// NewBaseModelCmd creates the patch basemodel command
func NewBaseModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "basemodel",
		Short: "Patch a base model",
		Long:  `Patch a specific base model in the NStream AI platform`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement base model patching logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&baseModelName, "name", "n", "", "Base model name")
	cmd.Flags().StringVarP(&baseModelSpec, "spec", "s", "", "Base model specification in JSON format")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("spec")

	return cmd
}
