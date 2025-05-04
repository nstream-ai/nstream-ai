package get

import (
	"github.com/spf13/cobra"
)

var (
	bucketName         string
	bucketOutputFormat string
)

// NewBucketCmd creates the get bucket command
func NewBucketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bucket",
		Short: "Get bucket information",
		Long:  `Get detailed information about a specific bucket or list all buckets`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement bucket get logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&bucketName, "name", "n", "", "Bucket name (optional, lists all buckets if not specified)")
	cmd.Flags().StringVarP(&bucketOutputFormat, "output", "o", "table", "Output format (table, json, yaml)")

	return cmd
}
