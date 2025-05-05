package create

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/banner"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
	"github.com/spf13/cobra"
)

var (
	bucketName     string
	bucketProvider string
	bucketRegion   string
)

// MockBucket represents a bucket in the system
type MockBucket struct {
	Name      string
	Region    string
	Provider  string
	Size      string
	CreatedAt string
}

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

			// Get cluster details to check cloud provider
			var clusterCloudProvider string
			if cfg.Cluster.Name != "" {
				// Create a channel for loading animation
				done := make(chan bool)
				go ShowLoading("Fetching cluster details", done)

				// Mock gRPC call to get cluster details
				time.Sleep(1 * time.Second)
				clusterDetails, err := mockGetClusterDetails(cfg.Cluster.Name)
				if err != nil {
					done <- true
					return fmt.Errorf("failed to get cluster details: %v", err)
				}
				done <- true
				clusterCloudProvider = clusterDetails.CloudProvider
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

			// Mock gRPC call to get buckets
			time.Sleep(1 * time.Second)
			buckets, err := mockListBuckets("")
			if err != nil {
				done <- true
				return fmt.Errorf("failed to get buckets: %v", err)
			}
			done <- true

			// Filter buckets by cloud provider
			var compatibleBuckets []MockBucket
			for _, bucket := range buckets {
				if bucket.Provider == clusterCloudProvider {
					compatibleBuckets = append(compatibleBuckets, bucket)
				}
			}

			// If there are compatible buckets, ask if user wants to use one
			if len(compatibleBuckets) > 0 {
				fmt.Printf("\nFound %d existing bucket(s) compatible with %s cloud provider:\n", len(compatibleBuckets), clusterCloudProvider)
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "ID\tName\tRegion\tProvider\tSize\tCreated At")
				for i, bucket := range compatibleBuckets {
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
					if err != nil || choiceInt < 1 || choiceInt > len(compatibleBuckets) {
						return fmt.Errorf("invalid bucket choice")
					}

					selectedBucket := compatibleBuckets[choiceInt-1]

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

			// Mock gRPC call to create bucket
			time.Sleep(2 * time.Second)
			bucketDetails, err := mockCreateBucket(name, region, clusterCloudProvider)
			if err != nil {
				done <- true
				return fmt.Errorf("failed to create bucket: %v", err)
			}
			done <- true

			// Update config with bucket details
			cfg.Cluster.Bucket = name
			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("failed to save config: %v", err)
			}

			fmt.Printf("\n%s✓ Successfully created bucket%s\n", boldColor, resetColor)
			fmt.Printf("\n%sBucket Details:%s\n", boldColor, resetColor)
			fmt.Printf("  Name: %s\n", bucketDetails.Name)
			fmt.Printf("  Region: %s\n", bucketDetails.Region)
			fmt.Printf("  Provider: %s\n", bucketDetails.Provider)
			fmt.Printf("  Size: %s\n", bucketDetails.Size)
			fmt.Printf("  Created At: %s\n", bucketDetails.CreatedAt)
			return nil
		},
	}

	cmd.Flags().StringVarP(&bucketName, "name", "n", "", "Bucket name (optional, will prompt if not provided)")
	cmd.Flags().StringVarP(&bucketProvider, "provider", "p", "", "Cloud provider (aws/gcp/azure)")
	cmd.Flags().StringVarP(&bucketRegion, "region", "r", "", "Region for the bucket")

	return cmd
}

// MockCreateBucket simulates creating a new bucket
func mockCreateBucket(name, region, provider string) (*MockBucket, error) {
	// Simulate gRPC call delay
	time.Sleep(2 * time.Second)

	// Return mock bucket details
	return &MockBucket{
		Name:      name,
		Region:    region,
		Provider:  provider,
		Size:      "0 B",
		CreatedAt: time.Now().Format("2006-01-02"),
	}, nil
}

// MockListBuckets simulates listing available buckets
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

// MockGetClusterDetails simulates getting cluster details
func mockGetClusterDetails(clusterName string) (*MockCluster, error) {
	// Simulate gRPC call delay
	time.Sleep(1 * time.Second)

	// Return mock cluster details
	return &MockCluster{
		ID:            clusterName,
		Region:        "us-west-2",
		CloudProvider: "aws",
		Bucket:        "",
		Role:          "",
	}, nil
}

// MockCluster represents a cluster in the system
type MockCluster struct {
	ID            string
	Region        string
	CloudProvider string
	Bucket        string
	Role          string
}
