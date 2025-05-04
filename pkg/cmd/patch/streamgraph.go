package patch

import (
	"github.com/spf13/cobra"
)

var (
	graphName string
	graphSpec string
)

// NewStreamGraphCmd creates the patch streamgraph command
func NewStreamGraphCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "streamgraph",
		Short: "Patch a stream graph",
		Long:  `Patch a specific stream graph in the NStream AI platform`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement stream graph patching logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&graphName, "name", "n", "", "Stream graph name")
	cmd.Flags().StringVarP(&graphSpec, "spec", "s", "", "Stream graph specification in JSON format")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("spec")

	return cmd
}
