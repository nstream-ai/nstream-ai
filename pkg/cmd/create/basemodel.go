package create

import (
	"github.com/spf13/cobra"
)

var (
	baseModelName       string
	baseModelType       string
	baseModelConfigPath string
)

// NewBaseModelCmd creates the basemodel command
func NewBaseModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "basemodel",
		Short: "Create a new base model",
		Long:  `Create a new base model with specified configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement base model creation logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&baseModelName, "name", "n", "", "Base model name")
	cmd.Flags().StringVarP(&baseModelType, "type", "t", "", "Base model type")
	cmd.Flags().StringVarP(&baseModelConfigPath, "config", "c", "", "Path to base model configuration file")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("config")

	return cmd
}
