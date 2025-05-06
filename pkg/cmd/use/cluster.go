package use

import (
	"context"
	"fmt"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/auth"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/banner"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/cluster"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var useClusterName string

// NewClusterCmd creates the cluster use command
func NewClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster [cluster-name]",
		Short: "Use a specific cluster",
		Long: `Set the current cluster context for operations.

If cluster name is provided as an argument, it will be used directly.
Otherwise, you'll be prompted to select from available clusters.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Print banner
			banner.PrintBanner()
			fmt.Println("Setting up cluster context...")
			fmt.Println()

			// Create validator and validate authentication
			validator, err := auth.NewValidator()
			if err != nil {
				return fmt.Errorf("failed to create validator: %v", err)
			}
			defer validator.Close()

			ctx := context.Background()
			if err := validator.ValidateAll(ctx); err != nil {
				fmt.Println("\nAuthentication failed:")
				fmt.Printf("Error: %v\n", err)
				fmt.Println("\nPlease try the following:")
				fmt.Println("1. Sign in again: 'nsai auth signin'")
				fmt.Println("2. If that doesn't work, sign up: 'nsai auth signup'")
				return fmt.Errorf(utils.ErrAuthRequired)
			}

			// Create cluster operations
			ops, err := cluster.NewOperations()
			if err != nil {
				return fmt.Errorf("failed to create cluster operations: %v", err)
			}
			defer ops.Close()

			// Get cluster name from args or flag
			var clusterName string
			if len(args) > 0 {
				clusterName = args[0]
			} else if useClusterName != "" {
				clusterName = useClusterName
			}

			// If cluster name is provided, use it directly
			if clusterName != "" {
				// Create a channel for loading animation
				done := make(chan bool)
				go utils.ShowDefaultLoading("Fetching cluster details", done)

				// Get cluster details
				details, err := ops.GetClusterDetails(ctx, clusterName)
				if err != nil {
					done <- true
					return err
				}

				done <- true

				// Update config with cluster details
				if err := ops.UpdateConfig(clusterName, details); err != nil {
					return fmt.Errorf("failed to save config: %v", err)
				}

				fmt.Printf("\r%s%s✓ Successfully set cluster context%s\n", utils.BoldColor, utils.RedColor, utils.ResetColor)
				fmt.Printf("\n%sCluster Details:%s\n", utils.BoldColor, utils.ResetColor)
				fmt.Printf("  Name: %s\n", clusterName)
				fmt.Printf("  Region: %s\n", details.Region)
				fmt.Printf("  Cloud Provider: %s\n", details.CloudProvider)
				fmt.Printf("  Bucket: %s\n", details.Bucket)
				fmt.Printf("  Identity: %s\n", details.Role)
				fmt.Printf("  Auth Token: %s\n", details.ClusterToken)
				fmt.Printf("\n%sYou can now use this cluster for operations.%s\n", utils.BoldColor, utils.ResetColor)
				return nil
			}

			// Create a channel for loading animation
			done := make(chan bool)
			go utils.ShowDefaultLoading("Fetching available clusters", done)

			// List clusters
			clusters, err := ops.ListClusters(ctx)
			if err != nil {
				done <- true
				return err
			}

			done <- true

			if len(clusters) == 0 {
				fmt.Println("\nNo clusters available.")
				fmt.Println("Please create a cluster first using 'nsai create cluster'")
				return fmt.Errorf(utils.ErrNoClusters)
			}

			// Display clusters
			fmt.Println("\nAvailable clusters:")
			ops.DisplayClusters(clusters)

			fmt.Print("\nEnter the number of the cluster to use: ")
			var choice int
			fmt.Scanf("%d", &choice)

			if choice < 1 || choice > len(clusters) {
				return fmt.Errorf(utils.ErrInvalidChoice)
			}

			// Get details for selected cluster
			selectedCluster := clusters[choice-1]
			done = make(chan bool)
			go utils.ShowDefaultLoading("Fetching cluster details", done)

			details, err := ops.GetClusterDetails(ctx, selectedCluster.Id)
			if err != nil {
				done <- true
				return err
			}

			done <- true

			// Update config with cluster details
			if err := ops.UpdateConfig(selectedCluster.Id, details); err != nil {
				return fmt.Errorf("failed to save config: %v", err)
			}

			fmt.Printf("\r%s%s✓ Successfully set cluster context%s\n", utils.BoldColor, utils.RedColor, utils.ResetColor)
			fmt.Printf("\n%sCluster Details:%s\n", utils.BoldColor, utils.ResetColor)
			fmt.Printf("  Name: %s\n", selectedCluster.Id)
			fmt.Printf("  Region: %s\n", details.Region)
			fmt.Printf("  Cloud Provider: %s\n", details.CloudProvider)
			fmt.Printf("  Bucket: %s\n", details.Bucket)
			fmt.Printf("  Identity: %s\n", details.Role)
			fmt.Printf("  Auth Token: %s\n", details.ClusterToken)
			fmt.Printf("\n%sYou can now use this cluster for operations.%s\n", utils.BoldColor, utils.ResetColor)
			return nil
		},
	}

	cmd.Flags().StringVarP(&useClusterName, "name", "n", "", "Cluster name (optional, will prompt for selection if not provided)")
	return cmd
}
