package use

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/banner"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
	"github.com/spf13/cobra"
)

const (
	redColor     = "\033[31m"
	boldColor    = "\033[1m"
	resetColor   = "\033[0m"
	blinkColor   = "\033[5m"
	reverseColor = "\033[7m"
)

var (
	useClusterName string
	useBucketName  string
	useRoleName    string
)

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

			// Check if config exists
			configPath := filepath.Join(os.Getenv("HOME"), ".nstreamconfig")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				fmt.Println("No configuration found. You need to authenticate first.")
				fmt.Println("\nPlease choose one of the following options:")
				fmt.Println("1. Sign in to an existing account: 'nsai auth signin'")
				fmt.Println("2. Create a new account: 'nsai auth signup'")
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
				return fmt.Errorf("authentication required")
			}

			// Create a channel for loading animation
			done := make(chan bool)

			// Get cluster name from args or flag
			var clusterName string
			if len(args) > 0 {
				clusterName = args[0]
			} else if useClusterName != "" {
				clusterName = useClusterName
			}

			// If cluster name is provided, use it directly
			if clusterName != "" {
				// Start loading animation
				go ShowLoading("Verifying cluster access", done)

				// Mock gRPC call to verify cluster exists
				time.Sleep(2 * time.Second)
				exists, err := mockVerifyClusterExists(clusterName)
				if err != nil {
					done <- true
					return fmt.Errorf("failed to verify cluster: %v", err)
				}

				done <- true

				if !exists {
					return fmt.Errorf("cluster %s does not exist", clusterName)
				}

				// Start loading animation for fetching cluster details
				done = make(chan bool)
				go ShowLoading("Fetching cluster details", done)

				// Mock gRPC call to get cluster details
				time.Sleep(2 * time.Second)
				clusterDetails, err := mockGetClusterDetails(clusterName)
				if err != nil {
					done <- true
					return fmt.Errorf("failed to get cluster details: %v", err)
				}

				done <- true

				// Update config with cluster details
				cfg.Cluster = *clusterDetails
				if err := config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("failed to save config: %v", err)
				}

				fmt.Printf("\r%s%s✓ Successfully set cluster context%s\n", boldColor, redColor, resetColor)
				fmt.Printf("\n%sCluster Details:%s\n", boldColor, resetColor)
				fmt.Printf("  Name: %s\n", clusterName)
				fmt.Printf("  Region: %s\n", clusterDetails.Region)
				fmt.Printf("  Cloud Provider: %s\n", clusterDetails.CloudProvider)
				fmt.Printf("  Bucket: %s\n", clusterDetails.Bucket)
				fmt.Printf("  Identity: %s\n", clusterDetails.Role)
				fmt.Printf("  Auth Token: %s\n", clusterDetails.ClusterToken)
				fmt.Printf("\n%sYou can now use this cluster for operations.%s\n", boldColor, resetColor)
				return nil
			}

			// Start loading animation for fetching clusters
			go ShowLoading("Fetching available clusters", done)

			// Mock gRPC call to list clusters
			time.Sleep(2 * time.Second)
			clusters, err := mockListClusters()
			if err != nil {
				done <- true
				return fmt.Errorf("failed to list clusters: %v", err)
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

			// Start loading animation for fetching cluster details
			done = make(chan bool)
			go ShowLoading("Fetching cluster details", done)

			// Mock gRPC call to get cluster details
			time.Sleep(2 * time.Second)
			selectedCluster := clusters[choice-1]
			clusterDetails, err := mockGetClusterDetails(selectedCluster.ID)
			if err != nil {
				done <- true
				return fmt.Errorf("failed to get cluster details: %v", err)
			}

			done <- true

			// Update config with cluster details
			cfg.Cluster = *clusterDetails
			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("failed to save config: %v", err)
			}

			fmt.Printf("\r%s%s✓ Successfully set cluster context%s\n", boldColor, redColor, resetColor)
			fmt.Printf("\n%sCluster Details:%s\n", boldColor, resetColor)
			fmt.Printf("  Name: %s\n", selectedCluster.ID)
			fmt.Printf("  Region: %s\n", clusterDetails.Region)
			fmt.Printf("  Cloud Provider: %s\n", clusterDetails.CloudProvider)
			fmt.Printf("  Bucket: %s\n", clusterDetails.Bucket)
			fmt.Printf("  Identity: %s\n", clusterDetails.Role)
			fmt.Printf("  Auth Token: %s\n", clusterDetails.ClusterToken)
			fmt.Printf("\n%sYou can now use this cluster for operations.%s\n", boldColor, resetColor)
			return nil
		},
	}

	cmd.Flags().StringVarP(&useClusterName, "name", "n", "", "Cluster name (optional, will prompt for selection if not provided)")

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

// ShowLoading displays a loading animation
func ShowLoading(message string, done chan bool) {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	colors := []string{redColor, boldColor, blinkColor, reverseColor}
	i := 0
	j := 0
	lastTime := time.Now()
	elapsed := 0

	for {
		select {
		case <-done:
			fmt.Printf("\r%s%s%s %s✓%s\n", boldColor, redColor, message, blinkColor, resetColor)
			return
		default:
			now := time.Now()
			elapsed += int(now.Sub(lastTime).Milliseconds())
			lastTime = now

			if elapsed >= 500 {
				j = (j + 1) % len(colors)
				elapsed = 0
			}

			dynamicMessage := fmt.Sprintf("%s%s%s%s", colors[j], boldColor, message, resetColor)

			extra := ""
			if elapsed < 100 {
				extra = " 🔥"
			} else if elapsed < 200 {
				extra = " ⚡"
			} else if elapsed < 300 {
				extra = " 💫"
			} else if elapsed < 400 {
				extra = " ✨"
			}

			fmt.Printf("\r%s %s%s", dynamicMessage, spinner[i], extra)
			time.Sleep(100 * time.Millisecond)
			i = (i + 1) % len(spinner)
		}
	}
}

// Mock gRPC service calls

type MockCluster struct {
	ID            string
	Region        string
	CloudProvider string
	Bucket        string
	Role          string
}

func mockListClusters() ([]MockCluster, error) {
	// Simulate gRPC call delay
	time.Sleep(1 * time.Second)

	// Return mock clusters with more detailed data
	return []MockCluster{
		{
			ID:            "cluster-1",
			Region:        "us-west-2",
			CloudProvider: "aws",
			Bucket:        "nstream-bucket-1",
			Role:          "nstream-role-1",
		},
		{
			ID:            "cluster-2",
			Region:        "us-central1",
			CloudProvider: "gcp",
			Bucket:        "nstream-bucket-2",
			Role:          "nstream-sa-2",
		},
		{
			ID:            "cluster-3",
			Region:        "eastus",
			CloudProvider: "azure",
			Bucket:        "nstream-bucket-3",
			Role:          "nstream-sp-3",
		},
	}, nil
}

func mockVerifyClusterExists(clusterName string) (bool, error) {
	// Simulate gRPC call delay
	time.Sleep(1 * time.Second)

	// Check if cluster exists in mock data
	clusters, _ := mockListClusters()
	for _, cluster := range clusters {
		if cluster.ID == clusterName {
			return true, nil
		}
	}
	return false, nil
}

func mockGetClusterDetails(clusterName string) (*config.ClusterConfig, error) {
	// Simulate gRPC call delay
	time.Sleep(1 * time.Second)

	// Get cluster details from mock data
	clusters, _ := mockListClusters()
	for _, cluster := range clusters {
		if cluster.ID == clusterName {
			return &config.ClusterConfig{
				Name:          cluster.ID,
				Region:        cluster.Region,
				CloudProvider: cluster.CloudProvider,
				Bucket:        cluster.Bucket,
				Role:          cluster.Role,
				ClusterToken:  "dummy-cluster-token-1234567890",
			}, nil
		}
	}
	return nil, fmt.Errorf("cluster %s not found", clusterName)
}
