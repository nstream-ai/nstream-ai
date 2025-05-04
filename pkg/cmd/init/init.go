package init

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	username      string
	password      string
	clusterName   string
	createCluster bool
	region        string
	cloud         string
	bucket        string
	role          string
)

// NewInitCmd creates the init command
func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the NStream AI CLI configuration",
		Long: `Initialize the NStream AI CLI by setting up your user configuration and cluster settings.

This command supports both interactive and non-interactive workflows:

Interactive Workflow:
- If no flags are provided, you'll be guided through:
  1. Sign in or Sign up process
  2. Cluster selection (create new or use existing)

Non-interactive Workflow:
- Use flags to specify your configuration:
  --user, -u: Username for authentication
  --password, -p: Password for authentication
  --cluster, -c: Cluster name to use
  --create-cluster: Create a new cluster (requires --region, --cloud, --bucket, and --role)
  --region: Region for cluster creation
  --cloud: Cloud provider for cluster creation
  --bucket: Bucket name for cluster creation
  --role: Role for cluster creation

Examples:
  # Interactive setup
  nsai init

  # Sign in with username and password
  nsai init --user myuser --password mypass

  # Create a new cluster
  nsai init --create-cluster --region us-west-2 --cloud aws --bucket mybucket --role myrole

  # Use existing cluster
  nsai init --cluster mycluster`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := filepath.Join(os.Getenv("HOME"), ".nstreamconfig")

			// Check if config file exists
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				// No config file exists, use default workflow
				fmt.Println("No configuration found. Starting with default workflow...")

				// Handle authentication first
				if err := handleAuthentication(configPath); err != nil {
					return err
				}

				// Then handle cluster operations
				if err := handleClusterOperations(configPath); err != nil {
					return err
				}

				return nil
			}

			// Config file exists, read it
			config, err := readConfig(configPath)
			if err != nil {
				return fmt.Errorf("error reading config: %v", err)
			}

			// If no auth token exists or user flag is provided, handle authentication
			if config["auth_token"] == "" || username != "" {
				if err := handleAuthentication(configPath); err != nil {
					return err
				}
			}

			// Handle cluster operations
			if err := handleClusterOperations(configPath); err != nil {
				return err
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&username, "user", "u", "", "Username for authentication")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication")
	cmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "Cluster name to use")
	cmd.Flags().BoolVar(&createCluster, "create-cluster", false, "Create a new cluster")
	cmd.Flags().StringVar(&region, "region", "", "Region for cluster creation")
	cmd.Flags().StringVar(&cloud, "cloud", "", "Cloud provider for cluster creation")
	cmd.Flags().StringVar(&bucket, "bucket", "", "Bucket name for cluster creation")
	cmd.Flags().StringVar(&role, "role", "", "Role for cluster creation")

	return cmd
}

func handleAuthentication(configPath string) error {
	// If username and password are provided, proceed with signin
	if username != "" && password != "" {
		authToken, err := authenticateWithGRPC(username, password)
		if err != nil {
			return fmt.Errorf("authentication failed: %v", err)
		}

		// Update config with auth token
		if err := updateConfig(configPath, map[string]string{
			"username":   username,
			"auth_token": authToken,
		}); err != nil {
			return fmt.Errorf("error updating config: %v", err)
		}
		return nil
	}

	// If only username is provided, ask for password
	if username != "" {
		fmt.Print("Enter password: ")
		fmt.Scanf("%s", &password)

		authToken, err := authenticateWithGRPC(username, password)
		if err != nil {
			return fmt.Errorf("authentication failed: %v", err)
		}

		// Update config with auth token
		if err := updateConfig(configPath, map[string]string{
			"username":   username,
			"auth_token": authToken,
		}); err != nil {
			return fmt.Errorf("error updating config: %v", err)
		}
		return nil
	}

	// No credentials provided, ask user what they want to do
	fmt.Println("\nWelcome to NStream AI CLI!")
	fmt.Println("Please choose an option:")
	fmt.Println("1. Sign In")
	fmt.Println("2. Sign Up")
	fmt.Print("\nEnter your choice (1 or 2): ")

	var choice int
	fmt.Scanf("%d", &choice)

	switch choice {
	case 1:
		return handleSignIn(configPath)
	case 2:
		return handleSignUp(configPath)
	default:
		return fmt.Errorf("invalid choice")
	}
}

func handleSignIn(configPath string) error {
	fmt.Print("Enter username: ")
	fmt.Scanf("%s", &username)

	fmt.Print("Enter password: ")
	fmt.Scanf("%s", &password)

	authToken, err := authenticateWithGRPC(username, password)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Update config with auth token
	if err := updateConfig(configPath, map[string]string{
		"username":   username,
		"auth_token": authToken,
	}); err != nil {
		return fmt.Errorf("error updating config: %v", err)
	}

	return nil
}

func handleSignUp(configPath string) error {
	var email, orgName, role string

	fmt.Print("Enter email: ")
	fmt.Scanf("%s", &email)

	fmt.Print("Enter organization name: ")
	fmt.Scanf("%s", &orgName)

	fmt.Print("Enter role: ")
	fmt.Scanf("%s", &role)

	// TODO: Implement signup via gRPC
	// After successful signup, proceed with signin
	return handleSignIn(configPath)
}

func handleClusterOperations(configPath string) error {
	// If create-cluster flag is set, all required flags must be provided
	if createCluster {
		if region == "" || cloud == "" || bucket == "" || role == "" {
			return fmt.Errorf("when --create-cluster is set, all of --region, --cloud, --bucket, and --role are required")
		}

		createdClusterName, err := createNewCluster(region, cloud, bucket, role)
		if err != nil {
			return fmt.Errorf("error creating cluster: %v", err)
		}
		clusterName = createdClusterName
	} else if clusterName != "" {
		// If cluster name is provided (but not creating), verify it exists
		exists, err := verifyClusterExists(clusterName)
		if err != nil {
			return fmt.Errorf("error verifying cluster: %v", err)
		}
		if !exists {
			return fmt.Errorf("cluster %s does not exist", clusterName)
		}
	} else {
		// Interactive cluster selection
		fmt.Println("\nCluster Setup")
		fmt.Println("1. Create new cluster")
		fmt.Println("2. Use existing cluster")
		fmt.Print("\nEnter your choice (1 or 2): ")

		var choice int
		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			// Get cluster creation prerequisites
			fmt.Print("Enter region: ")
			fmt.Scanf("%s", &region)

			fmt.Print("Enter cloud provider: ")
			fmt.Scanf("%s", &cloud)

			fmt.Print("Enter bucket name: ")
			fmt.Scanf("%s", &bucket)

			fmt.Print("Enter role: ")
			fmt.Scanf("%s", &role)

			createdClusterName, err := createNewCluster(region, cloud, bucket, role)
			if err != nil {
				return fmt.Errorf("error creating cluster: %v", err)
			}
			clusterName = createdClusterName
		case 2:
			// List available clusters and let user choose
			clusters, err := listClusters()
			if err != nil {
				return fmt.Errorf("error listing clusters: %v", err)
			}

			if len(clusters) == 0 {
				return fmt.Errorf("no clusters available")
			}

			fmt.Println("\nAvailable clusters:")
			for i, cluster := range clusters {
				fmt.Printf("%d. %s\n", i+1, cluster)
			}

			fmt.Print("\nEnter the number of the cluster to use: ")
			var choice int
			fmt.Scanf("%d", &choice)

			if choice < 1 || choice > len(clusters) {
				return fmt.Errorf("invalid cluster choice")
			}

			clusterName = clusters[choice-1]
		default:
			return fmt.Errorf("invalid choice")
		}
	}

	// Update config with cluster information
	if err := updateConfig(configPath, map[string]string{
		"cluster": clusterName,
	}); err != nil {
		return fmt.Errorf("error updating config: %v", err)
	}

	return nil
}

// Helper functions
func listClusters() ([]string, error) {
	// TODO: Implement cluster listing via gRPC
	return []string{"cluster1", "cluster2"}, nil
}

func readConfig(configPath string) (map[string]string, error) {
	// TODO: Implement config reading
	return make(map[string]string), nil
}

func authenticateWithGRPC(username, password string) (string, error) {
	// TODO: Implement authentication via gRPC
	return "dummy-token", nil
}

func verifyClusterExists(clusterName string) (bool, error) {
	// TODO: Implement cluster verification via gRPC
	return true, nil
}

func createNewCluster(region, cloud, bucket, role string) (string, error) {
	// TODO: Implement cluster creation via gRPC
	return "new-cluster", nil
}

func updateConfig(configPath string, updates map[string]string) error {
	// TODO: Implement config updating
	return nil
}
