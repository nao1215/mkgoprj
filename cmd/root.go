package cmd

import (
	"fmt"
	"os"

	"github.com/nao1215/ubume/internal/completion"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ubume",
	Short: `ubume make go project for command line application or library`,
}

func exitError(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

// Execute start command.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	completion.DeployShellCompletionFileIfNeeded(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		exitError(err)
	}
}
