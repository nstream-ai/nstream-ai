package create

import (
	"github.com/spf13/cobra"
)

var (
	knowledgeBaseName       string
	knowledgeBaseType       string
	knowledgeBaseConfigPath string
)

// NewKnowledgeBaseCmd creates the knowledgebase command
func NewKnowledgeBaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "knowledgebase",
		Short: "Create a new knowledge base",
		Long:  `Create a new knowledge base with specified configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement knowledge base creation logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&knowledgeBaseName, "name", "n", "", "Knowledge base name")
	cmd.Flags().StringVarP(&knowledgeBaseType, "type", "t", "", "Knowledge base type")
	cmd.Flags().StringVarP(&knowledgeBaseConfigPath, "config", "c", "", "Path to knowledge base configuration file")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("config")

	return cmd
}
