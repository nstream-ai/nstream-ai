package get

import (
	"github.com/spf13/cobra"
)

var (
	embeddingModelName         string
	embeddingModelOutputFormat string
)

// NewEmbeddingModelCmd creates the get embeddingmodel command
func NewEmbeddingModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "embeddingmodel",
		Short: "Get embedding model information",
		Long:  `Get detailed information about a specific embedding model or list all embedding models`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement embedding model get logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&embeddingModelName, "name", "n", "", "Embedding model name (optional, lists all embedding models if not specified)")
	cmd.Flags().StringVarP(&embeddingModelOutputFormat, "output", "o", "table", "Output format (table, json, yaml)")

	return cmd
}
