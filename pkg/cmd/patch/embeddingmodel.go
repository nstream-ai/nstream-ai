package patch

import (
	"github.com/spf13/cobra"
)

var (
	embeddingModelName string
	embeddingModelSpec string
)

// NewEmbeddingModelCmd creates the patch embeddingmodel command
func NewEmbeddingModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "embeddingmodel",
		Short: "Patch an embedding model",
		Long:  `Patch a specific embedding model in the NStream AI platform`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement embedding model patching logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&embeddingModelName, "name", "n", "", "Embedding model name")
	cmd.Flags().StringVarP(&embeddingModelSpec, "spec", "s", "", "Embedding model specification in JSON format")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("spec")

	return cmd
}
