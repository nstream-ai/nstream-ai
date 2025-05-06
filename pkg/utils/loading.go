package utils

import (
	"fmt"
	"time"
)

const (
	redColor     = "\033[31m"
	boldColor    = "\033[1m"
	resetColor   = "\033[0m"
	blinkColor   = "\033[5m"
	reverseColor = "\033[7m"
)

// LoadingConfig holds configuration for the loading animation
type LoadingConfig struct {
	Message     string
	Spinner     []string
	Colors      []string
	Interval    time.Duration
	Emoji       []string
	EmojiTiming []int
}

// DefaultConfig returns the default loading animation configuration
func DefaultConfig(message string) LoadingConfig {
	return LoadingConfig{
		Message:     message,
		Spinner:     []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		Colors:      []string{redColor, boldColor, blinkColor, reverseColor},
		Interval:    100 * time.Millisecond,
		Emoji:       []string{" 🔥", " ⚡", " 💫", " ✨"},
		EmojiTiming: []int{100, 200, 300, 400},
	}
}

// ShowLoading displays a loading animation with the given configuration
func ShowLoading(config LoadingConfig, done chan bool) {
	i := 0
	j := 0
	lastTime := time.Now()
	elapsed := 0

	for {
		select {
		case <-done:
			fmt.Printf("\r%s%s%s %s✓%s\n", boldColor, redColor, config.Message, blinkColor, resetColor)
			return
		default:
			now := time.Now()
			elapsed += int(now.Sub(lastTime).Milliseconds())
			lastTime = now

			if elapsed >= 500 {
				j = (j + 1) % len(config.Colors)
				elapsed = 0
			}

			dynamicMessage := fmt.Sprintf("%s%s%s%s", config.Colors[j], boldColor, config.Message, resetColor)

			extra := ""
			if len(config.EmojiTiming) > 0 && len(config.Emoji) > 0 {
				for k, timing := range config.EmojiTiming {
					if k < len(config.Emoji) && elapsed < timing {
						extra = config.Emoji[k]
						break
					}
				}
			}

			fmt.Printf("\r%s %s%s", dynamicMessage, config.Spinner[i], extra)
			time.Sleep(config.Interval)
			i = (i + 1) % len(config.Spinner)
		}
	}
}

// ShowDefaultLoading is a convenience function that uses the default configuration
func ShowDefaultLoading(message string, done chan bool) {
	ShowLoading(DefaultConfig(message), done)
}
