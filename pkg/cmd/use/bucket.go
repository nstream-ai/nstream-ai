package use

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/banner"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/client"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/utils"
	authproto "github.com/nstreama-ai/nstream-ai-cli/proto/auth"
	clusterproto "github.com/nstreama-ai/nstream-ai-cli/proto/cluster"
	"github.com/spf13/cobra"
)

var (
	bucketUseName     string
	bucketRoleName    string
	bucketClusterName string
)

// NewBucketCmd creates the bucket use command
func NewBucketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bucket [bucket-name]",
		Short: "Use a specific bucket in a cluster",
		Long: `Set the current bucket context for operations within a cluster.

If bucket name is provided as an argument, it will be used directly.
Otherwise, you'll be prompted to select from available buckets.

You must have a cluster context set or provide a cluster name to use this command.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Print banner
			banner.PrintBanner()
			fmt.Println("Setting up bucket context...")
			fmt.Println()

			// Check if config exists
			configPath := filepath.Join(os.Getenv("HOME"), ".nstreamconfig")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				fmt.Println("No configuration found. You need to authenticate first.")
				fmt.Println("\nPlease choose one of the following options:")
				fmt.Println("1. Sign in to an existing account: 'nsai auth signin'")
				fmt.Println("2. Create a new account: 'nsai auth signup'")
				fmt.Println("\nAfter authentication, run 'nsai use bucket' again")
				return fmt.Errorf("authentication required")
			}

			// Load config to check user credentials
			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %v", err)
			}

			// Check if user is authenticated
			if cfg.User.AuthToken == "" {
				fmt.Println("No authentication token found. You need to sign in first.")
				fmt.Println("\nRun 'nsai auth signin' to authenticate")
				fmt.Println("After authentication, run 'nsai use bucket' again")
				return fmt.Errorf("authentication required")
			}

			// Create gRPC client
			c, err := client.NewClient("", true, "")
			if err != nil {
				return err // Return the original error message without wrapping
			}
			defer c.Close()

			// Create context
			ctx, cancel := c.WithContext(context.Background())
			defer cancel()

			// Validate user
			validateResp, err := c.AuthClient.ValidateUser(ctx, &authproto.ValidateUserRequest{
				Email: cfg.User.Email,
			})
			if err != nil {
				return fmt.Errorf("failed to validate user: %v", err)
			}

			if !validateResp.Valid {
				fmt.Println("User validation failed. Please authenticate first:")
				fmt.Println("1. Sign in: 'nsai auth signin'")
				fmt.Println("2. Sign up: 'nsai auth signup'")
				return fmt.Errorf("authentication required")
			}

			// Validate token
			tokenResp, err := c.AuthClient.ValidateToken(ctx, &authproto.ValidateTokenRequest{
				Token: cfg.User.AuthToken,
			})
			if err != nil {
				return fmt.Errorf("error validating token: %v", err)
			}

			if !tokenResp.Valid {
				fmt.Printf("Authentication token is invalid: %s\n", tokenResp.Error)
				fmt.Println("\nPlease authenticate first:")
				fmt.Println("1. Sign in: 'nsai auth signin'")
				fmt.Println("2. Sign up: 'nsai auth signup'")
				return fmt.Errorf("authentication required")
			}

			// Get cluster name from flag or config
			var clusterName string
			if bucketClusterName != "" {
				clusterName = bucketClusterName
			} else if cfg.Cluster.Name != "" {
				// Create a channel for loading animation
				done := make(chan bool)
				go utils.ShowDefaultLoading("Fetching available clusters", done)

				// List clusters
				listResp, err := c.ClusterClient.ListClusters(ctx, &clusterproto.ListClustersRequest{})
				if err != nil {
					done <- true
					return fmt.Errorf("failed to get clusters: %v", err)
				}

				done <- true

				if len(listResp.Clusters) == 0 {
					fmt.Println("\nNo clusters available.")
					fmt.Println("Please create a cluster first using 'nsai create cluster'")
					return fmt.Errorf("no clusters available")
				}

				// Display clusters in a table
				fmt.Println("\nAvailable clusters:")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "ID\tRegion\tCloud\tBucket\tIdentity")
				for i, cluster := range listResp.Clusters {
					fmt.Fprintf(w, "%d. %s\t%s\t%s\t%s\t%s\n",
						i+1,
						cluster.Id,
						cluster.Region,
						cluster.CloudProvider,
						cluster.Bucket,
						cluster.Role,
					)
				}
				w.Flush()

				fmt.Printf("\nEnter the number of the cluster to use (press Enter to use default cluster '%s'): ", cfg.Cluster.Name)
				var choice string
				fmt.Scanln(&choice)

				if choice == "" {
					// Use default cluster from config
					clusterName = cfg.Cluster.Name
				} else {
					// Convert choice to integer
					choiceInt, err := strconv.Atoi(choice)
					if err != nil || choiceInt < 1 || choiceInt > len(listResp.Clusters) {
						return fmt.Errorf("invalid cluster choice")
					}
					clusterName = listResp.Clusters[choiceInt-1].Id
				}
			} else {
				// Create a channel for loading animation
				done := make(chan bool)
				go utils.ShowDefaultLoading("Fetching available clusters", done)

				// List clusters
				listResp, err := c.ClusterClient.ListClusters(ctx, &clusterproto.ListClustersRequest{})
				if err != nil {
					done <- true
					return fmt.Errorf("failed to get clusters: %v", err)
				}

				done <- true

				if len(listResp.Clusters) == 0 {
					fmt.Println("\nNo clusters available.")
					fmt.Println("Please create a cluster first using 'nsai create cluster'")
					return fmt.Errorf("no clusters available")
				}

				// Display clusters in a table
				fmt.Println("\nAvailable clusters:")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "ID\tRegion\tCloud\tBucket\tIdentity")
				for i, cluster := range listResp.Clusters {
					fmt.Fprintf(w, "%d. %s\t%s\t%s\t%s\t%s\n",
						i+1,
						cluster.Id,
						cluster.Region,
						cluster.CloudProvider,
						cluster.Bucket,
						cluster.Role,
					)
				}
				w.Flush()

				fmt.Print("\nEnter the number of the cluster to use: ")
				var choice int
				fmt.Scanf("%d", &choice)

				if choice < 1 || choice > len(listResp.Clusters) {
					return fmt.Errorf("invalid cluster choice")
				}

				clusterName = listResp.Clusters[choice-1].Id
			}

			// Get bucket name from args or flag
			var bucketName string
			if len(args) > 0 {
				bucketName = args[0]
			} else if bucketUseName != "" {
				bucketName = bucketUseName
			}

			// If no bucket name is provided, list buckets and let user choose
			if bucketName == "" {
				// Create a channel for loading animation
				done := make(chan bool)
				go utils.ShowDefaultLoading("Fetching available buckets", done)

				// Get cluster details to check cloud provider
				detailsResp, err := c.ClusterClient.GetClusterDetails(ctx, &clusterproto.GetClusterDetailsRequest{
					ClusterName: clusterName,
				})
				if err != nil {
					done <- true
					return fmt.Errorf("failed to get cluster details: %v", err)
				}

				if detailsResp.Error != "" {
					done <- true
					return fmt.Errorf("failed to get cluster details: %s", detailsResp.Error)
				}

				// List buckets
				bucketsResp, err := c.BucketClient.ListBuckets(ctx, &clusterproto.ListBucketsRequest{
					CloudProvider: detailsResp.Config.CloudProvider,
				})
				if err != nil {
					done <- true
					return fmt.Errorf("failed to get buckets: %v", err)
				}

				done <- true

				if len(bucketsResp.Buckets) == 0 {
					fmt.Println("\nNo buckets available.")
					fmt.Println("Please create a bucket first using 'nsai create bucket'")
					return fmt.Errorf("no buckets available")
				}

				// Display buckets in a table
				fmt.Println("\nAvailable buckets:")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "ID\tName\tRegion\tProvider\tSize\tCreated At")
				for i, bucket := range bucketsResp.Buckets {
					fmt.Fprintf(w, "%d. %s\t%s\t%s\t%s\t%s\n",
						i+1,
						bucket.Name,
						bucket.Region,
						bucket.Provider,
						bucket.Size,
						bucket.CreatedAt,
					)
				}
				w.Flush()

				fmt.Print("\nEnter the number of the bucket to use: ")
				var choice int
				fmt.Scanf("%d", &choice)

				if choice < 1 || choice > len(bucketsResp.Buckets) {
					return fmt.Errorf("invalid bucket choice")
				}

				bucketName = bucketsResp.Buckets[choice-1].Name

				// Check if cloud providers match
				if detailsResp.Config.CloudProvider != bucketsResp.Buckets[choice-1].Provider {
					return fmt.Errorf("cloud provider mismatch: cluster '%s' uses '%s' but bucket '%s' uses '%s'",
						clusterName,
						detailsResp.Config.CloudProvider,
						bucketName,
						bucketsResp.Buckets[choice-1].Provider,
					)
				}
			}

			// Verify bucket access
			done := make(chan bool)
			go utils.ShowDefaultLoading("Verifying bucket access", done)

			accessResp, err := c.BucketClient.VerifyBucketAccess(ctx, &clusterproto.VerifyBucketAccessRequest{
				CloudProvider: cfg.Cluster.CloudProvider,
				Bucket:        bucketName,
				Role:          cfg.Cluster.Role,
			})
			if err != nil {
				done <- true
				return fmt.Errorf("failed to verify bucket access: %v", err)
			}

			if !accessResp.HasAccess {
				done <- true
				return fmt.Errorf("bucket access verification failed: %s", accessResp.Error)
			}

			done <- true

			// Check resource readiness
			done = make(chan bool)
			go utils.ShowDefaultLoading("Checking resource readiness", done)

			readyResp, err := c.BucketClient.CheckResourceReadiness(ctx, &clusterproto.CheckResourceReadinessRequest{
				CloudProvider: cfg.Cluster.CloudProvider,
				Bucket:        bucketName,
				Role:          cfg.Cluster.Role,
			})
			if err != nil {
				done <- true
				return fmt.Errorf("failed to check resource readiness: %v", err)
			}

			if !readyResp.Ready {
				done <- true
				return fmt.Errorf("resources not ready: %s", readyResp.Error)
			}

			done <- true

			// Update config with bucket details
			cfg.Cluster.Bucket = bucketName
			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("failed to save config: %v", err)
			}

			fmt.Printf("\r%s%s✓ Successfully set bucket context%s\n", utils.BoldColor, utils.RedColor, utils.ResetColor)
			fmt.Printf("\n%sBucket Details:%s\n", utils.BoldColor, utils.ResetColor)
			fmt.Printf("  Name: %s\n", bucketName)
			fmt.Printf("  Cluster: %s\n", clusterName)
			fmt.Printf("  Cloud Provider: %s\n", cfg.Cluster.CloudProvider)
			fmt.Printf("  Region: %s\n", cfg.Cluster.Region)
			fmt.Printf("  Identity: %s\n", cfg.Cluster.Role)
			fmt.Printf("\n%sYou can now use this bucket for operations.%s\n", utils.BoldColor, utils.ResetColor)
			return nil
		},
	}

	cmd.Flags().StringVarP(&bucketUseName, "name", "n", "", "Bucket name (optional, will prompt for selection if not provided)")
	cmd.Flags().StringVarP(&bucketClusterName, "cluster", "c", "", "Cluster name (optional, will use current cluster if not provided)")

	// Role subcommand
	roleCmd := &cobra.Command{
		Use:   "role",
		Short: "Use a specific role with the bucket",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement role context setting logic
			return nil
		},
	}
	roleCmd.Flags().StringVarP(&bucketRoleName, "name", "n", "", "Role name")
	roleCmd.MarkFlagRequired("name")

	cmd.AddCommand(roleCmd)
	return cmd
}
