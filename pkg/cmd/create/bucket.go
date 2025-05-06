package create

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/banner"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/client"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
	clusterproto "github.com/nstreama-ai/nstream-ai-cli/proto/cluster"
	"github.com/spf13/cobra"
)

var (
	bucketName     string
	bucketProvider string
	bucketRegion   string
)

// NewBucketCmd creates the bucket command
func NewBucketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bucket [bucket-name]",
		Short: "Create a new bucket or select an existing one",
		Long: `Create a new bucket or select an existing one for use with a cluster.

If creating a new bucket, you'll be prompted for:
- Bucket name
- Cloud provider
- Region

If selecting an existing bucket, you'll be shown a list of available buckets
that are compatible with your cluster's cloud provider.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Print banner
			banner.PrintBanner()
			fmt.Println("Setting up bucket...")
			fmt.Println()

			// Check if config exists
			configPath := filepath.Join(os.Getenv("HOME"), ".nstreamconfig")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				fmt.Println("No configuration found. You need to authenticate first.")
				fmt.Println("\nPlease choose one of the following options:")
				fmt.Println("1. Sign in to an existing account: 'nsai auth signin'")
				fmt.Println("2. Create a new account: 'nsai auth signup'")
				fmt.Println("\nAfter authentication, run 'nsai create bucket' again")
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
				fmt.Println("After authentication, run 'nsai create bucket' again")
				return fmt.Errorf("authentication required")
			}

			// Check if user exists
			if cfg.User.Email == "" {
				fmt.Println("No user found. You need to sign in first.")
				fmt.Println("\nRun 'nsai auth signin' to authenticate")
				fmt.Println("After authentication, run 'nsai create bucket' again")
				return fmt.Errorf("authentication required")
			}

			// Initialize gRPC client
			c, err := client.NewClient("localhost:8080", false, "")
			if err != nil {
				return fmt.Errorf("failed to create client: %v", err)
			}
			defer c.Close()

			ctx := context.Background()

			// Get cluster details to check cloud provider
			var clusterCloudProvider string
			if cfg.Cluster.Name != "" {
				// Create a channel for loading animation
				done := make(chan bool)
				go ShowLoading("Fetching cluster details", done)

				// Get cluster details
				detailsResp, err := c.ClusterClient.GetClusterDetails(ctx, &clusterproto.GetClusterDetailsRequest{
					ClusterName: cfg.Cluster.Name,
				})
				if err != nil {
					done <- true
					return fmt.Errorf("failed to get cluster details: %v", err)
				}
				if detailsResp.Error != "" {
					done <- true
					return fmt.Errorf("failed to get cluster details: %s", detailsResp.Error)
				}
				done <- true
				clusterCloudProvider = detailsResp.Config.CloudProvider
			} else {
				// If no cluster is set, ask for cloud provider
				fmt.Println("\nNo cluster context found. Please select a cloud provider:")
				fmt.Println("1. AWS")
				fmt.Println("2. GCP")
				fmt.Println("3. Azure")
				fmt.Print("\nEnter your choice (1-3): ")

				reader := bufio.NewReader(os.Stdin)
				choice, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read input: %v", err)
				}
				choice = strings.TrimSpace(choice)

				switch choice {
				case "1":
					clusterCloudProvider = "aws"
				case "2":
					clusterCloudProvider = "gcp"
				case "3":
					clusterCloudProvider = "azure"
				default:
					return fmt.Errorf("invalid cloud provider choice")
				}
			}

			// Check if there are existing buckets
			done := make(chan bool)
			go ShowLoading("Checking existing buckets", done)

			// List buckets
			bucketsResp, err := c.BucketClient.ListBuckets(ctx, &clusterproto.ListBucketsRequest{
				CloudProvider: clusterCloudProvider,
			})
			if err != nil {
				done <- true
				return fmt.Errorf("failed to get buckets: %v", err)
			}
			done <- true

			// If there are compatible buckets, ask if user wants to use one
			if len(bucketsResp.Buckets) > 0 {
				fmt.Printf("\nFound %d existing bucket(s) compatible with %s cloud provider:\n", len(bucketsResp.Buckets), clusterCloudProvider)
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

				fmt.Print("\nWould you like to use one of these buckets? (y/n): ")
				reader := bufio.NewReader(os.Stdin)
				useExisting, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read input: %v", err)
				}
				useExisting = strings.TrimSpace(strings.ToLower(useExisting))

				if useExisting == "y" || useExisting == "yes" {
					fmt.Print("\nEnter the number of the bucket to use: ")
					choice, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("failed to read input: %v", err)
					}
					choice = strings.TrimSpace(choice)
					choiceInt, err := strconv.Atoi(choice)
					if err != nil || choiceInt < 1 || choiceInt > len(bucketsResp.Buckets) {
						return fmt.Errorf("invalid bucket choice")
					}

					selectedBucket := bucketsResp.Buckets[choiceInt-1]

					// Update config with bucket details
					cfg.Cluster.Bucket = selectedBucket.Name
					if err := config.SaveConfig(cfg); err != nil {
						return fmt.Errorf("failed to save config: %v", err)
					}

					fmt.Printf("\n%s✓ Successfully set bucket context%s\n", boldColor, resetColor)
					fmt.Printf("\n%sBucket Details:%s\n", boldColor, resetColor)
					fmt.Printf("  Name: %s\n", selectedBucket.Name)
					fmt.Printf("  Region: %s\n", selectedBucket.Region)
					fmt.Printf("  Provider: %s\n", selectedBucket.Provider)
					fmt.Printf("  Size: %s\n", selectedBucket.Size)
					fmt.Printf("  Created At: %s\n", selectedBucket.CreatedAt)
					return nil
				}
			}

			// If no compatible buckets or user wants to create new, proceed with bucket creation
			fmt.Println("\nCreating a new bucket...")

			// Get bucket name
			var name string
			if len(args) > 0 {
				name = args[0]
			} else {
				fmt.Print("Enter bucket name: ")
				reader := bufio.NewReader(os.Stdin)
				name, err = reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read input: %v", err)
				}
				name = strings.TrimSpace(name)
			}

			// Get region
			region, err := getRegion(clusterCloudProvider)
			if err != nil {
				return err
			}

			// Create bucket
			done = make(chan bool)
			go ShowLoading("Creating bucket", done)

			// Create bucket using gRPC
			createResp, err := c.ClusterClient.CreateCluster(ctx, &clusterproto.CreateClusterRequest{
				Name:          name,
				Type:          "basic",
				CloudProvider: clusterCloudProvider,
				Region:        region,
				Bucket:        name,
				Role:          "",
			})
			if err != nil {
				done <- true
				return fmt.Errorf("failed to create bucket: %v", err)
			}
			if createResp.Error != "" {
				done <- true
				return fmt.Errorf("failed to create bucket: %s", createResp.Error)
			}
			done <- true

			// Update config with bucket details
			cfg.Cluster.Bucket = name
			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("failed to save config: %v", err)
			}

			fmt.Printf("\n%s✓ Successfully created bucket%s\n", boldColor, resetColor)
			fmt.Printf("\n%sBucket Details:%s\n", boldColor, resetColor)
			fmt.Printf("  Name: %s\n", createResp.Config.Bucket)
			fmt.Printf("  Region: %s\n", createResp.Config.Region)
			fmt.Printf("  Provider: %s\n", createResp.Config.CloudProvider)
			return nil
		},
	}

	cmd.Flags().StringVarP(&bucketName, "name", "n", "", "Bucket name (optional, will prompt if not provided)")
	cmd.Flags().StringVarP(&bucketProvider, "provider", "p", "", "Cloud provider (aws/gcp/azure)")
	cmd.Flags().StringVarP(&bucketRegion, "region", "r", "", "Region for the bucket")

	return cmd
}
