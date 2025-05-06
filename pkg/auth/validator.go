package auth

import (
	"context"
	"fmt"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/client"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
	authproto "github.com/nstreama-ai/nstream-ai-cli/proto/auth"
)

// Validator handles authentication validation
type Validator struct {
	client *client.Client
	config *config.Config
}

// NewValidator creates a new Validator instance
func NewValidator() (*Validator, error) {
	c, err := client.NewClient("", true, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	return &Validator{
		client: c,
		config: cfg,
	}, nil
}

// ValidateUser checks if the user is valid
func (v *Validator) ValidateUser(ctx context.Context) error {
	validateResp, err := v.client.AuthClient.ValidateUser(ctx, &authproto.ValidateUserRequest{
		Email: v.config.User.Email,
	})
	if err != nil {
		return fmt.Errorf("failed to validate user: %v", err)
	}

	if !validateResp.Valid {
		return fmt.Errorf("user validation failed")
	}

	return nil
}

// ValidateToken checks if the auth token is valid
func (v *Validator) ValidateToken(ctx context.Context) error {
	tokenResp, err := v.client.AuthClient.ValidateToken(ctx, &authproto.ValidateTokenRequest{
		Token: v.config.User.AuthToken,
	})
	if err != nil {
		return fmt.Errorf("error validating token: %v", err)
	}

	if !tokenResp.Valid {
		return fmt.Errorf("authentication token is invalid: %s", tokenResp.Error)
	}

	return nil
}

// ValidateClusterToken checks if the cluster token is valid
func (v *Validator) ValidateClusterToken(ctx context.Context) error {
	if v.config.Cluster.ClusterToken == "" {
		return nil // No cluster token to validate
	}

	clusterResp, err := v.client.AuthClient.ValidateClusterToken(ctx, &authproto.ValidateClusterTokenRequest{
		Token: v.config.Cluster.ClusterToken,
	})
	if err != nil {
		return fmt.Errorf("error validating cluster token: %v", err)
	}

	if !clusterResp.Valid {
		return fmt.Errorf("cluster token is invalid: %s", clusterResp.Error)
	}

	return nil
}

// ValidateAll performs all validation checks
func (v *Validator) ValidateAll(ctx context.Context) error {
	if err := v.ValidateUser(ctx); err != nil {
		return err
	}

	if err := v.ValidateToken(ctx); err != nil {
		return err
	}

	if err := v.ValidateClusterToken(ctx); err != nil {
		return err
	}

	return nil
}

// Close closes the client connection
func (v *Validator) Close() {
	if v.client != nil {
		v.client.Close()
	}
}
