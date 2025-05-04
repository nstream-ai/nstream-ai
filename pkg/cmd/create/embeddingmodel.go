package create

import (
	"github.com/spf13/cobra"
)

var (
	embeddingModelName       string
	embeddingModelType       string
	embeddingModelConfigPath string
)

// NewEmbeddingModelCmd creates the embeddingmodel command
func NewEmbeddingModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "embeddingmodel",
		Short: "Create a new embedding model",
		Long:  `Create a new embedding model with specified configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement embedding model creation logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&embeddingModelName, "name", "n", "", "Embedding model name")
	cmd.Flags().StringVarP(&embeddingModelType, "type", "t", "", "Embedding model type")
	cmd.Flags().StringVarP(&embeddingModelConfigPath, "config", "c", "", "Path to embedding model configuration file")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("config")

	return cmd
}
