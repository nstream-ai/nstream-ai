package init

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
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
		Long: `Initialize the NStream AI CLI by setting up your configuration.

This command supports both interactive and non-interactive workflows:

Interactive Workflow:
- If no flags are provided, you'll be guided through:
  1. Authentication (using auth commands)
  2. Cluster setup (using create/use commands)

Non-interactive Workflow:
- Use flags to specify your configuration:
  --cluster, -c: Cluster name to use
  --create-cluster: Create a new cluster (requires --region, --cloud, --bucket, and --role)
  --region: Region for cluster creation
  --cloud: Cloud provider for cluster creation
  --bucket: Bucket name for cluster creation
  --role: Role for cluster creation

Examples:
  # Interactive setup
  nsai init

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

				// First ensure authentication is set up
				if err := ensureAuthentication(); err != nil {
					return err
				}

				// Then handle cluster operations
				if err := handleClusterOperations(); err != nil {
					return err
				}

				return nil
			}

			// Config file exists, handle cluster operations
			if err := handleClusterOperations(); err != nil {
				return err
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "Cluster name to use")
	cmd.Flags().BoolVar(&createCluster, "create-cluster", false, "Create a new cluster")
	cmd.Flags().StringVar(&region, "region", "", "Region for cluster creation")
	cmd.Flags().StringVar(&cloud, "cloud", "", "Cloud provider for cluster creation")
	cmd.Flags().StringVar(&bucket, "bucket", "", "Bucket name for cluster creation")
	cmd.Flags().StringVar(&role, "role", "", "Role for cluster creation")

	return cmd
}

func ensureAuthentication() error {
	// Check if user is authenticated
	configPath := filepath.Join(os.Getenv("HOME"), ".nstreamconfig")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("\nNo configuration file found. Authentication required.")
		return promptForAuth()
	}

	// Read config file
	config, err := readConfig(configPath)
	if err != nil {
		return fmt.Errorf("error reading config: %v", err)
	}

	// Check if auth token exists and is valid
	token, exists := config["auth_token"]
	if !exists || token == "" {
		fmt.Println("\nNo authentication token found. Authentication required.")
		return promptForAuth()
	}

	// Check if token is expired
	if isTokenExpired(token) {
		fmt.Println("\nAuthentication token has expired. Re-authentication required.")
		return promptForAuth()
	}

	return nil
}

func promptForAuth() error {
	fmt.Println("\nPlease choose an option:")
	fmt.Println("1. Sign In (nsai auth signin)")
	fmt.Println("2. Sign Up (nsai auth signup)")
	fmt.Print("\nEnter your choice (1 or 2): ")

	var choice int
	fmt.Scanf("%d", &choice)

	switch choice {
	case 1:
		return runSignin()
	case 2:
		return runSignup()
	default:
		return fmt.Errorf("invalid choice")
	}
}

func isTokenExpired(token string) bool {
	// Execute auth check command
	cmd := exec.Command(os.Args[0], "auth", "check")
	cmd.Env = append(os.Environ(), fmt.Sprintf("NSTREAM_AUTH_TOKEN=%s", token))

	// Run silently (don't show output)
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Run()
	return err != nil // If error, token is expired/invalid
}

func handleClusterOperations() error {
	if createCluster {
		// Use create cluster command
		return runCreateCluster()
	} else if clusterName != "" {
		// Use existing cluster
		return runUseCluster(clusterName)
	} else {
		// Interactive mode
		fmt.Println("\nPlease choose an option:")
		fmt.Println("1. Create new cluster (nsai create cluster)")
		fmt.Println("2. Use existing cluster (nsai use cluster)")
		fmt.Print("\nEnter your choice (1 or 2): ")

		var choice int
		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			return runCreateCluster()
		case 2:
			return runUseCluster("")
		default:
			return fmt.Errorf("invalid choice")
		}
	}
}

func runSignin() error {
	// Execute signin command
	cmd := exec.Command(os.Args[0], "auth", "signin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runSignup() error {
	// Execute signup command
	cmd := exec.Command(os.Args[0], "auth", "signup")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runCreateCluster() error {
	// Execute create cluster command with appropriate flags
	args := []string{"create", "cluster"}

	// Prompt for cluster name if not provided
	var clusterName string
	fmt.Print("\nEnter cluster name: ")
	fmt.Scanf("%s", &clusterName)
	args = append(args, clusterName)

	if region != "" {
		args = append(args, "--region", region)
	}
	if cloud != "" {
		args = append(args, "--cloud", cloud)
	}
	if bucket != "" {
		args = append(args, "--bucket", bucket)
	}
	if role != "" {
		args = append(args, "--role", role)
	}

	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runUseCluster(name string) error {
	// Execute use cluster command
	args := []string{"use", "cluster"}
	if name != "" {
		args = append(args, name)
	}

	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func readConfig(configPath string) (map[string]string, error) {
	// Read and parse config file
	config := make(map[string]string)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse config file (implement your config parsing logic here)
	// This is a placeholder - implement actual config parsing
	_ = data // TODO: Parse data into config map
	return config, nil
}
