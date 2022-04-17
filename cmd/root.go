package cmd

import (
	"fmt"
	"os"

	"github.com/nao1215/mkgoprj/v2/internal/completion"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mkgoprj",
	Short: `mkgoprj make go project for command line application or library`,
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
