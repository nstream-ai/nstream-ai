package delete

import (
	"github.com/spf13/cobra"
)

var (
	bucketName  string
	forceDelete bool
)

// NewBucketCmd creates the delete bucket command
func NewBucketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bucket",
		Short: "Delete a bucket",
		Long:  `Delete a specific bucket from the NStream AI platform`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement bucket deletion logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&bucketName, "name", "n", "", "Bucket name")
	cmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Force deletion without confirmation")

	cmd.MarkFlagRequired("name")

	return cmd
}
