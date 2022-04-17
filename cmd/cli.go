package cmd

import (
	"os"

	"github.com/nao1215/ubume/internal/ioutils"
	"github.com/nao1215/ubume/internal/project"
	"github.com/spf13/cobra"
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Make golang project for command line interface with cobra",
	Long: `Make golang project for command line interface with cobra.
You need to specify IMPORT_PATH as argument. ※ IMPORT_PATH is same as $ go mod init IMPORT_PATH`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(cli(cmd, args))
	},
}

func init() {
	cliCmd.Flags().BoolP("no-root", "n", false, "Create files in the current directory without creating the project root directory")
	rootCmd.AddCommand(cliCmd)
}

func cli(cmd *cobra.Command, args []string) int {
	if len(args) == 0 {
		ioutils.Die("need import path or project name")
	}

	noRoot, err := cmd.Flags().GetBool("no-root")
	if err != nil {
		ioutils.Die("can not parse command line argument (--no-root)")
	}

	prj := project.NewProject(args[0], false, true, noRoot)
	prj.Make()

	return 0
}
