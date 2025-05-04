package auth

import (
	"fmt"
	"time"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
)

const (
	redColor     = "\033[31m"
	boldColor    = "\033[1m"
	resetColor   = "\033[0m"
	blinkColor   = "\033[5m"
	reverseColor = "\033[7m"
)

// ShowLoading displays a colorful and interactive loading spinner with the given message
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
			// Show completion with a colorful checkmark
			fmt.Printf("\r%s%s%s %s✓%s\n", boldColor, redColor, message, blinkColor, resetColor)
			return
		default:
			now := time.Now()
			elapsed += int(now.Sub(lastTime).Milliseconds())
			lastTime = now

			// Change color every 500ms
			if elapsed >= 500 {
				j = (j + 1) % len(colors)
				elapsed = 0
			}

			// Create a dynamic message with color
			dynamicMessage := fmt.Sprintf("%s%s%s%s", colors[j], boldColor, message, resetColor)

			// Add some interactive elements based on time
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

// DummySendPasswordEmail simulates sending password generation email
func DummySendPasswordEmail(email string) error {
	// Simulate network delay for email sending
	time.Sleep(1 * time.Second)
	return nil
}

// DummyFetchUserDetails simulates fetching user and cluster details
func DummyFetchUserDetails(authToken string) (*config.Config, error) {
	// Simulate network delay
	time.Sleep(1 * time.Second)

	// Create a dummy response with user and cluster details
	cfg := &config.Config{
		User: config.UserConfig{
			Email:     "user@example.com",
			OrgName:   "nstream-ai",
			Role:      "developer",
			AuthToken: authToken,
		},
		Cluster: config.ClusterConfig{
			Name:          "default-cluster",
			CloudProvider: "aws",
			Region:        "us-west-2",
			Bucket:        "nstream-ai-bucket",
			ClusterToken:  "dummy-cluster-token-1234567890",
		},
	}

	return cfg, nil
}
