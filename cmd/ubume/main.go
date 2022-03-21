//
// ubume/cmd/ubume/main.go
//
// Copyright 2022 Naohiro CHIKAMATSU
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"github.com/nao1215/ubume/internal/interactive"
	"github.com/nao1215/ubume/internal/ioutils"
	"github.com/nao1215/ubume/internal/project"
)

// options is ubume command options.
type options struct {
	CLI         bool `short:"c" long:"cli" description:"Generate cli project with cobra (default: application project)"`
	Interactive bool `short:"i" long:"interactive" description:"Generate cli project with interactive mode"`
	Libary      bool `short:"l" long:"library" description:"Generate library project template (default: application project)"`
	NoRoot      bool `short:"n" long:"no-root" description:"Create files in the current directory without creating the project root directory"`
	Version     bool `short:"v" long:"version" description:"Show ubume command version"`
}

// osExit is funtion pointer that prepare a function pointer for unit testing
var osExit = os.Exit

const version string = "1.5.2"

// main is entry point of ubume command.
func main() {
	var opts options
	args := parseArgs(&opts)
	prj := project.NewProject(args[0], opts.Libary, opts.CLI, opts.NoRoot)
	prj.Make()
}

// parseArgs parse command line arguments.
// In this method, process for version option, help option, and lack of arguments.
func parseArgs(opts *options) []string {
	p := newParser(opts)

	args, err := p.Parse()
	if err != nil {
		// If user specify --help option, help message already showed in p.Parse().
		// Moreover, help messages are stored inside err.
		if !strings.Contains(err.Error(), ioutils.CmdName) {
			showHelp(p)
			showHelpFooter()
		} else {
			showHelpFooter()
		}
		osExit(1)
	}

	if opts.Version {
		showVersion(ioutils.CmdName, version)
		osExit(0)
	}

	if opts.Interactive {
		args, err := interact(opts)
		if err != nil {
			ioutils.Die(err.Error())
		}
		return args
	}

	if len(args) != 1 {
		showHelp(p)
		showHelpFooter()
		osExit(1)
	}

	if opts.CLI && opts.Libary {
		ioutils.Die("can not specify --cli and --library at same time")
		osExit(1)
	}
	return args
}

// showHelp print help messages.
func showHelp(p *flags.Parser) {
	p.WriteHelp(os.Stdout)
}

// showHelpFooter print author contact information.
func showHelpFooter() {
	fmt.Println("")
	fmt.Println("Contact:")
	fmt.Println("  If you find the bugs, please report the content of the error.")
	fmt.Println("  [GitHub Issue] https://github.com/nao1215/ubume/issues")
}

// newParser return initialized flags.Parser.
func newParser(opts *options) *flags.Parser {
	parser := flags.NewParser(opts, flags.Default)
	parser.Name = ioutils.CmdName
	parser.Usage = "[OPTIONS] IMPORT_PATH  ※ IMPORT_PATH is same as $ go mod init IMPORT_PATH"
	return parser
}

// showVersion show ubume command version information.
func showVersion(cmdName string, version string) {
	description := cmdName + " version " + version + " (under Apache License version 2.0)"
	fmt.Fprintln(os.Stdout, description)
}

func interact(opts *options) ([]string, error) {
	ip, err := interactive.ImportPath()
	if err != nil {
		return nil, err
	}
	projectName := filepath.Base(ip)

	kind, err := interactive.ProjectKind()
	if err != nil {
		return nil, err
	}

	rootInfo, err := interactive.ProjectRootDir(projectName)
	if err != nil {
		return nil, err
	}

	if !interactive.FinalConfirm(ip, projectName, kind, rootInfo) {
		fmt.Println("")
		fmt.Println(ioutils.CmdName + ": Please try again:" + color.YellowString("'ubume --interactive'"))
		os.Exit(0)
	}

	opts.CLI = interactive.IsCLI(kind)
	opts.Libary = interactive.IsLibrary(kind)
	opts.NoRoot = interactive.IsNoRoot(rootInfo)

	return []string{ip}, nil
}
