package create

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/api"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/banner"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
	"github.com/spf13/cobra"
)

// NewClusterCmd creates the cluster command
func NewClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster [cluster-name]",
		Short: "Create a new NStream AI cluster",
		Long:  `Create a new NStream AI cluster with specified configuration`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createCluster(args[0])
		},
	}

	cmd.Flags().StringP("type", "t", "basic", "Cluster type (basic/standard/enterprise)")
	cmd.Flags().StringP("cloud", "c", "gcp", "Cloud provider (aws/gcp/azure)")
	cmd.Flags().StringP("region", "r", "", "Region for the cluster")
	cmd.Flags().StringP("bucket", "b", "", "Bucket name for storage")
	cmd.Flags().StringP("role", "p", "", "Role/principal to assume for bucket access")

	return cmd
}

func createCluster(name string) error {
	// Print banner
	banner.PrintBanner()
	fmt.Println("Creating a new NStream AI cluster...")
	fmt.Println()

	// Check if config exists
	if !config.ConfigExists() {
		fmt.Println("No configuration found. Please authenticate first:")
		fmt.Println("1. Sign in: 'nsai auth signin'")
		fmt.Println("2. Sign up: 'nsai auth signup'")
		return fmt.Errorf("authentication required")
	}

	// Load config to check user credentials
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Check if user is authenticated
	if cfg.User.AuthToken == "" {
		fmt.Println("No authentication token found. Please authenticate first:")
		fmt.Println("1. Sign in: 'nsai auth signin'")
		fmt.Println("2. Sign up: 'nsai auth signup'")
		return fmt.Errorf("authentication required")
	}

	// Check if user exists
	valid, err := api.MockValidateUser(cfg.User.Email)
	if err != nil || !valid {
		fmt.Println("User validation failed. Please authenticate first:")
		fmt.Println("1. Sign in: 'nsai auth signin'")
		fmt.Println("2. Sign up: 'nsai auth signup'")
		return fmt.Errorf("authentication required")
	}

	// Check if token is valid
	resp, err := api.MockValidateToken(cfg.User.AuthToken)
	if err != nil {
		return fmt.Errorf("error validating token: %v", err)
	}

	if !resp.Valid {
		fmt.Printf("Authentication token is invalid: %s\n", resp.Error)
		fmt.Println("\nPlease authenticate first:")
		fmt.Println("1. Sign in: 'nsai auth signin'")
		fmt.Println("2. Sign up: 'nsai auth signup'")
		return fmt.Errorf("authentication required")
	}

	// Get cluster type
	clusterType, err := getClusterType()
	if err != nil {
		return err
	}

	// Get cloud provider
	cloudProvider, err := getCloudProvider()
	if err != nil {
		return err
	}

	// Get region
	region, err := getRegion(cloudProvider)
	if err != nil {
		return err
	}

	// Check for existing buckets with matching cloud provider
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
		if bucket.Provider == cloudProvider {
			compatibleBuckets = append(compatibleBuckets, bucket)
		}
	}

	var bucket string
	// If there are compatible buckets, ask if user wants to use one
	if len(compatibleBuckets) > 0 {
		fmt.Printf("\nFound %d existing bucket(s) compatible with %s cloud provider:\n", len(compatibleBuckets), cloudProvider)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tName\tRegion\tProvider\tSize\tCreated At")
		for i, b := range compatibleBuckets {
			fmt.Fprintf(w, "%d. %s\t%s\t%s\t%s\t%s\n",
				i+1,
				b.Name,
				b.Region,
				b.Provider,
				b.Size,
				b.CreatedAt,
			)
		}
		w.Flush()

		fmt.Print("\nWould you like to use one of these buckets? (y/n): ")
		var useExisting string
		fmt.Scanln(&useExisting)

		if useExisting == "y" || useExisting == "yes" {
			fmt.Print("\nEnter the number of the bucket to use: ")
			var choice int
			fmt.Scanln(&choice)

			if choice < 1 || choice > len(compatibleBuckets) {
				return fmt.Errorf("invalid bucket choice")
			}

			bucket = compatibleBuckets[choice-1].Name
		} else {
			// Get new bucket name
			fmt.Print("\nEnter your bucket name: ")
			fmt.Scanln(&bucket)
		}
	} else {
		// Get new bucket name
		fmt.Print("\nEnter your bucket name: ")
		fmt.Scanln(&bucket)
	}

	// Get role/principal for bucket access
	fmt.Printf("\nEnter the name for your bucket access %s: ", getRoleType(cloudProvider))
	var userRole string
	fmt.Scanln(&userRole)

	// Get NStream service role
	done = make(chan bool)
	go ShowLoading("Fetching NStream service role", done)
	serviceRole, err := DummyGetServiceRole(cloudProvider)
	if err != nil {
		done <- true
		return fmt.Errorf("failed to get service role: %v", err)
	}
	done <- true

	// Show cloud-specific setup instructions
	fmt.Println("\nPlease follow these steps to set up bucket access:")
	fmt.Println(GetCloudSetupInstructions(cloudProvider, serviceRole, userRole))
	fmt.Print("\nPress Enter when you have completed the setup...")
	fmt.Scanln()

	// Verify bucket access
	done = make(chan bool)
	go ShowLoading("Verifying bucket access", done)
	if err := VerifyBucketAccess(cloudProvider, bucket, userRole); err != nil {
		done <- true
		return fmt.Errorf("failed to verify bucket access: %v", err)
	}
	done <- true

	// Check resource readiness
	done = make(chan bool)
	go ShowLoading("Checking resource readiness", done)
	if err := CheckResourceReadiness(cloudProvider, bucket, userRole); err != nil {
		done <- true
		return fmt.Errorf("resources not ready: %v", err)
	}
	done <- true

	// Create cluster
	done = make(chan bool)
	go ShowLoading("Creating your NStream AI cluster", done)

	clusterConfig, err := DummyCreateCluster(name, clusterType, cloudProvider, region, bucket, userRole)
	if err != nil {
		done <- true
		return fmt.Errorf("failed to create cluster: %v", err)
	}
	done <- true

	// Update config
	cfg.Cluster = *clusterConfig
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Println("\nSuccessfully created cluster!")
	fmt.Printf("Name: %s\n", clusterConfig.Name)
	fmt.Printf("Type: %s\n", clusterType)
	fmt.Printf("Cloud Provider: %s\n", cloudProvider)
	fmt.Printf("Region: %s\n", region)
	fmt.Printf("Bucket: %s\n", bucket)
	fmt.Printf("Bucket Access %s: %s\n", getRoleType(cloudProvider), userRole)
	fmt.Printf("NStream Service Role: %s\n", serviceRole)
	fmt.Println("\nYou can now use this cluster with 'nsai cluster use " + name + "'")

	return nil
}

func getClusterType() (string, error) {
	fmt.Println("\nAvailable cluster types:")
	fmt.Println("1. Basic (Free)")
	fmt.Println("2. Standard (Requires credits)")
	fmt.Println("3. Enterprise (Requires credits)")
	fmt.Print("\nSelect cluster type (1-3): ")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		return "basic", nil
	case 2:
		return "standard", nil
	case 3:
		return "enterprise", nil
	default:
		return "", fmt.Errorf("invalid cluster type selection")
	}
}

func getCloudProvider() (string, error) {
	fmt.Println("\nAvailable cloud providers:")
	fmt.Println("1. AWS")
	fmt.Println("2. GCP")
	fmt.Println("3. Azure")
	fmt.Print("\nSelect cloud provider (1-3): ")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		return "aws", nil
	case 2:
		return "gcp", nil
	case 3:
		return "azure", nil
	default:
		return "", fmt.Errorf("invalid cloud provider selection")
	}
}

func getRegion(provider string) (string, error) {
	regions := GetCloudRegions(provider)
	if len(regions) == 0 {
		return "", fmt.Errorf("no regions available for provider %s", provider)
	}

	fmt.Printf("\nAvailable regions for %s:\n", provider)
	for i, region := range regions {
		fmt.Printf("%d. %s\n", i+1, region)
	}
	fmt.Print("\nSelect region (1-", len(regions), "): ")

	var choice int
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(regions) {
		return "", fmt.Errorf("invalid region selection")
	}

	return regions[choice-1], nil
}

func getRoleType(provider string) string {
	switch provider {
	case "aws":
		return "IAM Role"
	case "gcp":
		return "Service Account"
	case "azure":
		return "Service Principal"
	default:
		return "Role"
	}
}
