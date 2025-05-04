package auth

import (
	"fmt"
	"time"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/banner"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
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

// dummyServiceCall simulates a gRPC service call
func dummyServiceCall(email, password string) (*config.Config, error) {
	// Simulate network delay
	time.Sleep(2 * time.Second)

	// Create a dummy response with just auth token
	cfg := &config.Config{
		User: config.UserConfig{
			Email:     email,
			AuthToken: "dummy-auth-token-1234567890",
		},
	}

	return cfg, nil
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

	// Make dummy service call to send password email
	if err := DummySendPasswordEmail(email); err != nil {
		done <- true
		return fmt.Errorf("failed to send password email: %v", err)
	}

	// Signal loading is complete
	done <- true

	// Get password from user
	fmt.Print("Enter the password received via email: ")
	var password string
	fmt.Scanln(&password)

	// Create a new channel for authentication loading
	done = make(chan bool)

	// Start loading animation for authentication
	go ShowLoading("Authenticating with NStream AI service", done)

	// Make dummy service call
	cfg, err := dummyServiceCall(email, password)
	if err != nil {
		done <- true
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Signal loading is complete
	done <- true

	// Create a new channel for fetching details
	done = make(chan bool)

	// Start loading animation for fetching details
	go ShowLoading("Fetching your account details", done)

	// Fetch additional user and cluster details
	detailedCfg, err := DummyFetchUserDetails(cfg.User.AuthToken)
	if err != nil {
		done <- true
		return fmt.Errorf("failed to fetch account details: %v", err)
	}

	// Signal loading is complete
	done <- true

	// Save config
	if err := config.SaveConfig(detailedCfg); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Println("\nSuccessfully signed in!")
	fmt.Printf("Organization: %s\n", detailedCfg.User.OrgName)
	fmt.Printf("Role: %s\n", detailedCfg.User.Role)
	fmt.Printf("Current Cluster: %s (%s)\n", detailedCfg.Cluster.Name, detailedCfg.Cluster.Region)
	fmt.Println("\nYou're all set! Start using NStream AI CLI with 'nsai --help'")
	return nil
}
