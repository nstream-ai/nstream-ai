package get

import (
	"github.com/spf13/cobra"
)

var (
	knowledgeBaseName         string
	knowledgeBaseOutputFormat string
)

// NewKnowledgeBaseCmd creates the get knowledgebase command
func NewKnowledgeBaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "knowledgebase",
		Short: "Get knowledge base information",
		Long:  `Get detailed information about a specific knowledge base or list all knowledge bases`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement knowledge base get logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&knowledgeBaseName, "name", "n", "", "Knowledge base name (optional, lists all knowledge bases if not specified)")
	cmd.Flags().StringVarP(&knowledgeBaseOutputFormat, "output", "o", "table", "Output format (table, json, yaml)")

	return cmd
}
