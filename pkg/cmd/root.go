package cmd

import (
	"github.com/nstreama-ai/nstream-ai-cli/pkg/banner"
	authcmd "github.com/nstreama-ai/nstream-ai-cli/pkg/cmd/auth"
	createcmd "github.com/nstreama-ai/nstream-ai-cli/pkg/cmd/create"
	initcmd "github.com/nstreama-ai/nstream-ai-cli/pkg/cmd/init"
	usecmd "github.com/nstreama-ai/nstream-ai-cli/pkg/cmd/use"
	"github.com/spf13/cobra"
)

var version = "dev" // This will be set during build

var rootCmd = &cobra.Command{
	Use:   "nsai",
	Short: "NStream AI CLI",
	Long:  `A command line interface for NStream AI platform`,
}

func init() {
	rootCmd.Version = version

	// Create a custom help template that includes the banner
	helpTemplate := banner.GetBanner() + "\n\n" + rootCmd.HelpTemplate()
	rootCmd.SetHelpTemplate(helpTemplate)

	// Add init command
	rootCmd.AddCommand(initcmd.NewInitCmd())

	// Add auth commands
	authcmd.RegisterCommands(rootCmd)

	// Add create command
	rootCmd.AddCommand(createcmd.NewCreateCmd())

	// Add use command
	rootCmd.AddCommand(usecmd.NewUseCmd())
}

func Execute() error {
	return rootCmd.Execute()
}
