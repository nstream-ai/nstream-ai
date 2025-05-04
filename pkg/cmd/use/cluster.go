package use

import (
	"github.com/spf13/cobra"
)

var (
	useClusterName string
	useBucketName  string
	useRoleName    string
)

// NewClusterCmd creates the cluster use command
func NewClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Use a specific cluster",
		Long:  `Set the current cluster context for operations`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement cluster context setting logic
			return nil
		},
	}

	cmd.Flags().StringVarP(&useClusterName, "name", "n", "", "Cluster name")
	cmd.MarkFlagRequired("name")

	// Bucket subcommand
	bucketCmd := &cobra.Command{
		Use:   "bucket",
		Short: "Use a specific bucket in the cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement bucket context setting logic
			return nil
		},
	}
	bucketCmd.Flags().StringVarP(&useBucketName, "name", "n", "", "Bucket name")
	bucketCmd.MarkFlagRequired("name")

	// Role subcommand
	roleCmd := &cobra.Command{
		Use:   "role",
		Short: "Use a specific role in the cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement role context setting logic
			return nil
		},
	}
	roleCmd.Flags().StringVarP(&useRoleName, "name", "n", "", "Role name")
	roleCmd.MarkFlagRequired("name")

	cmd.AddCommand(bucketCmd, roleCmd)
	return cmd
} 