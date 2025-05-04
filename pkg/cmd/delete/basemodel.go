package delete

import (
	"github.com/spf13/cobra"
)

var (
	baseModelName        string
	baseModelForceDelete bool
)

// NewBaseModelCmd creates the delete basemodel command
func NewBaseModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "basemodel",
		Short: "Delete a base model",
		Long:  `Delete a specific base model from the NStream AI platform`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement base model deletion logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&baseModelName, "name", "n", "", "Base model name")
	cmd.Flags().BoolVarP(&baseModelForceDelete, "force", "f", false, "Force deletion without confirmation")

	cmd.MarkFlagRequired("name")

	return cmd
}
