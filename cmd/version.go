package cmd

import (
	"fmt"

	"github.com/nao1215/ubume/internal/cmdinfo"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show " + cmdinfo.Name() + " command version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmdinfo.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
