package auth

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/banner"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/client"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
	authproto "github.com/nstreama-ai/nstream-ai-cli/proto/auth"
	clusterproto "github.com/nstreama-ai/nstream-ai-cli/proto/cluster"
	"github.com/spf13/cobra"
)

func NewSigninCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signin",
		Short: "Sign in to NStream AI platform",
		Long:  `Sign in to NStream AI platform using your email and password`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return signin()
		},
	}

	return cmd
}

func signin() error {
	// Print banner and welcome message
	banner.PrintBanner()
	fmt.Println("Welcome to NStream AI CLI!")
	fmt.Println("Please sign in to continue...")
	fmt.Println()

	// Get email from user
	fmt.Print("Enter your email: ")
	var email string
	fmt.Scanln(&email)

	// Create a channel to signal when loading is done
	done := make(chan bool)

	// Start loading animation for sending password email
	go ShowLoading("Sending password to your email", done)

	// Create gRPC client
	c, err := client.NewClient("", true, "")
	if err != nil {
		done <- true
		return fmt.Errorf("failed to create client: %v", err)
	}
	defer c.Close()

	// Call SignIn service
	ctx, cancel := c.WithContext(context.Background())
	defer cancel()

	signInResp, err := c.AuthClient.SignIn(ctx, &authproto.SignInRequest{
		Email: email,
	})
	if err != nil {
		done <- true
		return fmt.Errorf("failed to send sign in request: %v", err)
	}

	if !signInResp.Success {
		done <- true
		return fmt.Errorf("sign in failed: %s", signInResp.Error)
	}

	// Signal loading is complete
	done <- true

	// Get OTP from user
	fmt.Print("Enter the OTP received via email: ")
	var otp string
	fmt.Scanln(&otp)

	// Create a new channel for authentication loading
	done = make(chan bool)

	// Start loading animation for authentication
	go ShowLoading("Authenticating with NStream AI service", done)

	// Verify sign in with OTP
	verifyResp, err := c.AuthClient.VerifySignIn(ctx, &authproto.VerifySignInRequest{
		Email: email,
		Otp:   otp,
	})
	if err != nil {
		done <- true
		return fmt.Errorf("authentication failed: %v", err)
	}

	if verifyResp.Error != "" {
		done <- true
		return fmt.Errorf("authentication failed: %s", verifyResp.Error)
	}

	// Create config with auth token
	cfg := &config.Config{
		User: config.UserConfig{
			Email:     email,
			AuthToken: verifyResp.AuthToken,
			OrgName:   verifyResp.UserInfo.Organization,
			Role:      verifyResp.UserInfo.Role,
		},
	}

	// Signal loading is complete
	done <- true

	// Create a new channel for fetching cluster details
	done = make(chan bool)

	// Start loading animation for fetching cluster details
	go ShowLoading("Fetching your cluster details", done)

	// Get cluster details
	listClustersResp, err := c.ClusterClient.ListClusters(ctx, &clusterproto.ListClustersRequest{
		AuthToken: verifyResp.AuthToken,
	})
	if err != nil {
		done <- true
		return fmt.Errorf("failed to fetch cluster details: %v", err)
	}

	done <- true

	if len(listClustersResp.Clusters) > 0 {
		fmt.Println("\nAvailable clusters:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tRegion\tCloud\tBucket\tIdentity")
		for i, cluster := range listClustersResp.Clusters {
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

		fmt.Print("\nWould you like to select a cluster? (y/n): ")
		var choice string
		fmt.Scanln(&choice)

		if choice == "y" || choice == "yes" {
			fmt.Print("\nEnter the number of the cluster to use: ")
			var clusterChoice int
			fmt.Scanln(&clusterChoice)

			if clusterChoice < 1 || clusterChoice > len(listClustersResp.Clusters) {
				return fmt.Errorf("invalid cluster choice")
			}

			selectedCluster := listClustersResp.Clusters[clusterChoice-1]
			cfg.Cluster = config.ClusterConfig{
				Name:          selectedCluster.Id,
				Region:        selectedCluster.Region,
				CloudProvider: selectedCluster.CloudProvider,
				Bucket:        selectedCluster.Bucket,
				Role:          selectedCluster.Role,
			}
		}
	}

	// Save config
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Println("\nSuccessfully signed in!")
	fmt.Printf("Organization: %s\n", cfg.User.OrgName)
	fmt.Printf("Role: %s\n", cfg.User.Role)
	if cfg.Cluster.Name != "" {
		fmt.Printf("Selected Cluster: %s (%s)\n", cfg.Cluster.Name, cfg.Cluster.Region)
	}
	fmt.Println("\nYou're all set! Start using NStream AI CLI with 'nsai --help'")
	return nil
}
