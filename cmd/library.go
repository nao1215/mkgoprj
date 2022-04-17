package cmd

import (
	"os"

	"github.com/nao1215/mkgoprj/internal/ioutils"
	"github.com/nao1215/mkgoprj/internal/project"
	"github.com/spf13/cobra"
)

var libraryCmd = &cobra.Command{
	Use:   "library",
	Short: "Make golang library  project",
	Long: `Make golang project for command line interface with cobra.
You need to specify IMPORT_PATH as argument. ※ IMPORT_PATH is same as $ go mod init IMPORT_PATH`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(library(cmd, args))
	},
}

func init() {
	libraryCmd.Flags().BoolP("no-root", "n", false, "Create files in the current directory without creating the project root directory")
	rootCmd.AddCommand(libraryCmd)
}

func library(cmd *cobra.Command, args []string) int {
	if len(args) == 0 {
		ioutils.Die("need import path or project name")
	}

	noRoot, err := cmd.Flags().GetBool("no-root")
	if err != nil {
		ioutils.Die("can not parse command line argument (--no-root)")
	}

	prj := project.NewProject(args[0], true, false, noRoot)
	prj.Make()

	return 0
}
