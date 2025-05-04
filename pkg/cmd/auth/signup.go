package auth

import (
	"fmt"
	"time"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/banner"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
	"github.com/spf13/cobra"
)

func NewSignupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signup",
		Short: "Sign up for NStream AI platform",
		Long:  `Sign up for NStream AI platform using your email, organization, name, and role`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return signup()
		},
	}

	return cmd
}

// dummySignupServiceCall simulates a gRPC service call for signup
func dummySignupServiceCall(email, org, name, role, password string) (*config.Config, error) {
	// Simulate network delay for signup process
	time.Sleep(3 * time.Second)
	// Create a dummy response with just auth token
	cfg := &config.Config{
		User: config.UserConfig{
			Email:     email,
			AuthToken: "dummy-signup-token-1234567890",
			OrgName:   org,
			Role:      role,
		},
	}

	return cfg, nil
}

func signup() error {
	// Print banner and welcome message
	banner.PrintBanner()
	fmt.Println("Welcome to NStream AI CLI!")
	fmt.Println("Let's get you started with your new account...")
	fmt.Println()

	// Get user details
	fmt.Print("Enter your email: ")
	var email string
	fmt.Scanln(&email)

	fmt.Print("Enter your organization: ")
	var org string
	fmt.Scanln(&org)

	fmt.Print("Enter your name: ")
	var name string
	fmt.Scanln(&name)

	fmt.Print("Enter your role: ")
	var role string
	fmt.Scanln(&role)

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

	// Create a new channel for signup loading
	done = make(chan bool)

	// Start loading animation for signup process
	go ShowLoading("Creating your NStream AI account", done)

	// Make dummy service call for signup
	cfg, err := dummySignupServiceCall(email, org, name, role, password)
	if err != nil {
		done <- true
		return fmt.Errorf("signup failed: %v", err)
	}

	// Signal loading is complete
	done <- true

	// Create a new channel for fetching details
	done = make(chan bool)

	// Start loading animation for fetching details
	go ShowLoading("Setting up your account and cluster", done)

	// Fetch additional user and cluster details
	detailedCfg, err := DummyFetchUserDetails(cfg.User.AuthToken)
	if err != nil {
		done <- true
		return fmt.Errorf("failed to setup account: %v", err)
	}

	// Signal loading is complete
	done <- true

	// Save config
	if err := config.SaveConfig(detailedCfg); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Println("\nSuccessfully signed up!")
	fmt.Printf("Organization: %s\n", detailedCfg.User.OrgName)
	fmt.Printf("Role: %s\n", detailedCfg.User.Role)
	fmt.Printf("Current Cluster: %s (%s)\n", detailedCfg.Cluster.Name, detailedCfg.Cluster.Region)
	fmt.Println("\nYou're all set! Start using NStream AI CLI with 'nsai --help'")
	return nil
}
