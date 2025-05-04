package get

import (
	"github.com/spf13/cobra"
)

// NewGetCmd creates the root get command
func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get NStream AI platform resource information",
		Long:  `Get detailed information about various NStream AI platform resources`,
	}

	// Add all subcommands
	cmd.AddCommand(
		NewClusterCmd(),
		NewBucketCmd(),
		NewStreamGraphCmd(),
		NewStreamFinetunerCmd(),
		NewBaseModelCmd(),
		NewEmbeddingModelCmd(),
		NewMegaModelCmd(),
		NewKnowledgeBaseCmd(),
		NewStreamConnectorCmd(),
	)

	return cmd
}
