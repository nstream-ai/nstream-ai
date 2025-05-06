package auth

import (
	"context"
	"fmt"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/banner"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/client"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
	authproto "github.com/nstreama-ai/nstream-ai-cli/proto/auth"
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

	// Create gRPC client
	c, err := client.NewClient("", true, "")
	if err != nil {
		done <- true
		return fmt.Errorf("failed to create client: %v", err)
	}
	defer c.Close()

	// Call SignUp service
	ctx, cancel := c.WithContext(context.Background())
	defer cancel()

	signUpResp, err := c.AuthClient.SignUp(ctx, &authproto.SignUpRequest{
		Email:        email,
		Name:         name,
		Organization: org,
		Role:         role,
	})
	if err != nil {
		done <- true
		return fmt.Errorf("failed to send sign up request: %v", err)
	}

	if !signUpResp.Success {
		done <- true
		return fmt.Errorf("sign up failed: %s", signUpResp.Error)
	}

	// Signal loading is complete
	done <- true

	// Get OTP from user
	fmt.Print("Enter the OTP received via email: ")
	var otp string
	fmt.Scanln(&otp)

	// Create a new channel for signup loading
	done = make(chan bool)

	// Start loading animation for signup process
	go ShowLoading("Creating your NStream AI account", done)

	// Verify sign up with OTP
	verifyResp, err := c.AuthClient.VerifySignUp(ctx, &authproto.VerifySignUpRequest{
		Email: email,
		Otp:   otp,
	})
	if err != nil {
		done <- true
		return fmt.Errorf("signup failed: %v", err)
	}

	if verifyResp.Error != "" {
		done <- true
		return fmt.Errorf("signup failed: %s", verifyResp.Error)
	}

	// Create config with auth token
	cfg := &config.Config{
		User: config.UserConfig{
			Email:     email,
			AuthToken: verifyResp.AuthToken,
			OrgName:   verifyResp.UserInfo.Organization,
			Role:      verifyResp.UserInfo.Role,
		},
		Cluster: config.ClusterConfig{},
	}

	// Signal loading is complete
	done <- true

	// Save config
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Println("\nSuccessfully signed up!")
	fmt.Printf("Organization: %s\n", cfg.User.OrgName)
	fmt.Printf("Role: %s\n", cfg.User.Role)
	fmt.Println("\nYou're all set! Start using NStream AI CLI with 'nsai --help'")
	return nil
}
