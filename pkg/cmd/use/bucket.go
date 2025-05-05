package use

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/banner"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
	"github.com/spf13/cobra"
)

var (
	bucketUseName     string
	bucketRoleName    string
	bucketClusterName string
)

// MockBucket represents a bucket in the system
type MockBucket struct {
	Name      string
	Region    string
	Provider  string
	Size      string
	CreatedAt string
}

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

			// Check if user exists
			if cfg.User.Email == "" {
				fmt.Println("No user found. You need to sign in first.")
				fmt.Println("\nRun 'nsai auth signin' to authenticate")
				fmt.Println("After authentication, run 'nsai use bucket' again")
				return fmt.Errorf("authentication required")
			}

			// Get cluster name from flag or config
			var clusterName string
			if bucketClusterName != "" {
				clusterName = bucketClusterName
			} else if cfg.Cluster.Name != "" {
				// Create a channel for loading animation
				done := make(chan bool)
				go ShowLoading("Fetching available clusters", done)

				// Mock gRPC call to get clusters
				time.Sleep(1 * time.Second)
				clusters, err := mockListClusters()
				if err != nil {
					done <- true
					return fmt.Errorf("failed to get clusters: %v", err)
				}

				done <- true

				if len(clusters) == 0 {
					fmt.Println("\nNo clusters available.")
					fmt.Println("Please create a cluster first using 'nsai create cluster'")
					return fmt.Errorf("no clusters available")
				}

				// Display clusters in a table
				fmt.Println("\nAvailable clusters:")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "ID\tRegion\tCloud\tBucket\tIdentity")
				for i, cluster := range clusters {
					fmt.Fprintf(w, "%d. %s\t%s\t%s\t%s\t%s\n",
						i+1,
						cluster.ID,
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
					if err != nil || choiceInt < 1 || choiceInt > len(clusters) {
						return fmt.Errorf("invalid cluster choice")
					}
					clusterName = clusters[choiceInt-1].ID
				}
			} else {
				// Create a channel for loading animation
				done := make(chan bool)
				go ShowLoading("Fetching available clusters", done)

				// Mock gRPC call to get clusters
				time.Sleep(1 * time.Second)
				clusters, err := mockListClusters()
				if err != nil {
					done <- true
					return fmt.Errorf("failed to get clusters: %v", err)
				}

				done <- true

				if len(clusters) == 0 {
					fmt.Println("\nNo clusters available.")
					fmt.Println("Please create a cluster first using 'nsai create cluster'")
					return fmt.Errorf("no clusters available")
				}

				// Display clusters in a table
				fmt.Println("\nAvailable clusters:")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "ID\tRegion\tCloud\tBucket\tIdentity")
				for i, cluster := range clusters {
					fmt.Fprintf(w, "%d. %s\t%s\t%s\t%s\t%s\n",
						i+1,
						cluster.ID,
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

				if choice < 1 || choice > len(clusters) {
					return fmt.Errorf("invalid cluster choice")
				}

				clusterName = clusters[choice-1].ID
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
				go ShowLoading("Fetching available buckets", done)

				// Mock gRPC call to get buckets
				time.Sleep(2 * time.Second)
				buckets, err := mockListBuckets(clusterName)
				if err != nil {
					done <- true
					return fmt.Errorf("failed to get buckets: %v", err)
				}

				done <- true

				if len(buckets) == 0 {
					fmt.Println("\nNo buckets available.")
					fmt.Println("Please create a bucket first using 'nsai create bucket'")
					return fmt.Errorf("no buckets available")
				}

				// Display buckets in a table
				fmt.Println("\nAvailable buckets:")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "ID\tName\tRegion\tProvider\tSize\tCreated At")
				for i, bucket := range buckets {
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

				if choice < 1 || choice > len(buckets) {
					return fmt.Errorf("invalid bucket choice")
				}

				bucketName = buckets[choice-1].Name

				// Get cluster details to check cloud provider
				clusterDetails, err := mockGetClusterDetails(clusterName)
				if err != nil {
					return fmt.Errorf("failed to get cluster details: %v", err)
				}

				// Check if cloud providers match
				if clusterDetails.CloudProvider != buckets[choice-1].Provider {
					return fmt.Errorf("cloud provider mismatch: cluster '%s' uses '%s' but bucket '%s' uses '%s'",
						clusterName, clusterDetails.CloudProvider, bucketName, buckets[choice-1].Provider)
				}
			} else {
				// Get bucket details to check cloud provider
				bucketDetails, err := mockGetBucketDetails(clusterName, bucketName)
				if err != nil {
					return fmt.Errorf("failed to get bucket details: %v", err)
				}

				// Get cluster details to check cloud provider
				clusterDetails, err := mockGetClusterDetails(clusterName)
				if err != nil {
					return fmt.Errorf("failed to get cluster details: %v", err)
				}

				// Check if cloud providers match
				if clusterDetails.CloudProvider != bucketDetails.Provider {
					return fmt.Errorf("cloud provider mismatch: cluster '%s' uses '%s' but bucket '%s' uses '%s'",
						clusterName, clusterDetails.CloudProvider, bucketName, bucketDetails.Provider)
				}
			}

			// Start loading animation for fetching bucket details
			done := make(chan bool)
			go ShowLoading("Fetching bucket details", done)

			// Mock gRPC call to get bucket details
			time.Sleep(2 * time.Second)
			bucketDetails, err := mockGetBucketDetails(clusterName, bucketName)
			if err != nil {
				done <- true
				return fmt.Errorf("failed to get bucket details: %v", err)
			}

			done <- true

			// Update config with bucket details
			cfg.Cluster.Bucket = bucketName
			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("failed to save config: %v", err)
			}

			fmt.Printf("\r%s%s✓ Successfully set bucket context%s\n", boldColor, redColor, resetColor)
			fmt.Printf("\n%sBucket Details:%s\n", boldColor, resetColor)
			fmt.Printf("  Name: %s\n", bucketDetails.Name)
			fmt.Printf("  Region: %s\n", bucketDetails.Region)
			fmt.Printf("  Provider: %s\n", bucketDetails.Provider)
			fmt.Printf("  Size: %s\n", bucketDetails.Size)
			fmt.Printf("  Created At: %s\n", bucketDetails.CreatedAt)
			fmt.Printf("\n%sYou can now use this bucket for operations in cluster %s.%s\n", boldColor, clusterName, resetColor)
			return nil
		},
	}

	cmd.Flags().StringVarP(&bucketUseName, "name", "n", "", "Bucket name (optional, will prompt for selection if not provided)")
	cmd.Flags().StringVarP(&bucketClusterName, "cluster", "c", "", "Cluster name (optional, will use current cluster or prompt for selection if not provided)")

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

// Mock gRPC service calls

func mockListBuckets(clusterName string) ([]MockBucket, error) {
	// Simulate gRPC call delay
	time.Sleep(1 * time.Second)

	// Return mock buckets
	return []MockBucket{
		{
			Name:      "bucket-1",
			Region:    "us-west-2",
			Provider:  "aws",
			Size:      "1.2 TB",
			CreatedAt: "2024-01-01",
		},
		{
			Name:      "bucket-2",
			Region:    "us-central1",
			Provider:  "gcp",
			Size:      "2.5 TB",
			CreatedAt: "2024-02-01",
		},
		{
			Name:      "bucket-3",
			Region:    "eastus",
			Provider:  "azure",
			Size:      "3.1 TB",
			CreatedAt: "2024-03-01",
		},
	}, nil
}

func mockGetBucketDetails(clusterName, bucketName string) (*MockBucket, error) {
	// Simulate gRPC call delay
	time.Sleep(1 * time.Second)

	// Return mock bucket details
	return &MockBucket{
		Name:      bucketName,
		Region:    "us-west-2",
		Provider:  "aws",
		Size:      "1.2 TB",
		CreatedAt: "2024-01-01",
	}, nil
}
