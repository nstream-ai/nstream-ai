package create

import (
	"github.com/spf13/cobra"
)

var (
	graphName       string
	graphType       string
	graphConfigPath string
)

// NewStreamGraphCmd creates the streamgraph command
func NewStreamGraphCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "streamgraph",
		Short: "Create a new stream graph",
		Long:  `Create a new stream graph with specified configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement stream graph creation logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&graphName, "name", "n", "", "Graph name")
	cmd.Flags().StringVarP(&graphType, "type", "t", "", "Graph type")
	cmd.Flags().StringVarP(&graphConfigPath, "config", "c", "", "Path to graph configuration file")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("config")

	return cmd
}
