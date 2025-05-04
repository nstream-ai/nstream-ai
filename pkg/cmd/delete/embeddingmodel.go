package delete

import (
	"github.com/spf13/cobra"
)

var (
	embeddingModelName        string
	embeddingModelForceDelete bool
)

// NewEmbeddingModelCmd creates the delete embeddingmodel command
func NewEmbeddingModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "embeddingmodel",
		Short: "Delete an embedding model",
		Long:  `Delete a specific embedding model from the NStream AI platform`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement embedding model deletion logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&embeddingModelName, "name", "n", "", "Embedding model name")
	cmd.Flags().BoolVarP(&embeddingModelForceDelete, "force", "f", false, "Force deletion without confirmation")

	cmd.MarkFlagRequired("name")

	return cmd
}
