package init

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/api"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
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
			// First ensure authentication is set up
			if err := ensureAuthentication(); err != nil {
				return err
			}

			// Then handle cluster operations
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
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error reading config: %v", err)
	}

	// Check if auth token exists
	if cfg.User.AuthToken == "" {
		fmt.Println("\nNo authentication token found. Authentication required.")
		return promptForAuth()
	}

	// Check if user exists
	valid, err := api.MockValidateUser(cfg.User.Email)
	if err != nil || !valid {
		fmt.Println("\nUser validation failed. Authentication required.")
		return promptForAuth()
	}

	// Check if token is valid
	resp, err := api.MockValidateToken(cfg.User.AuthToken)
	if err != nil {
		return fmt.Errorf("error validating token: %v", err)
	}

	if !resp.Valid {
		fmt.Printf("\nAuthentication token is invalid: %s\n", resp.Error)
		fmt.Println("Authentication required.")
		return promptForAuth()
	}

	// If cluster token exists, validate it too
	if cfg.Cluster.ClusterToken != "" {
		clusterResp, err := api.MockValidateClusterToken(cfg.Cluster.ClusterToken)
		if err != nil {
			return fmt.Errorf("error validating cluster token: %v", err)
		}

		if !clusterResp.Valid {
			fmt.Printf("\nCluster token is invalid: %s\n", clusterResp.Error)
			// Don't require re-authentication for invalid cluster token
			// Just clear it from config
			cfg.Cluster.ClusterToken = ""
			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("error saving config: %v", err)
			}
		}
	}

	return nil
}

func promptForAuth() error {
	fmt.Println("\nPlease choose an option:")
	fmt.Println("1. Sign In")
	fmt.Println("2. Sign Up")
	fmt.Print("\nEnter your choice (1 or 2): ")

	var choice int
	fmt.Scanf("%d", &choice)

	switch choice {
	case 1:
		return runSignin()
	case 2:
		return runSignup()
	default:
		fmt.Println("\nInvalid choice. Please try again.")
		return promptForAuth()
	}
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
	err := cmd.Run()
	if err != nil {
		return err
	}

	// After successful signin, continue with cluster operations
	return handleClusterOperations()
}

func runSignup() error {
	// Execute signup command
	cmd := exec.Command(os.Args[0], "auth", "signup")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		return err
	}

	// After successful signup, continue with cluster operations
	return handleClusterOperations()
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
