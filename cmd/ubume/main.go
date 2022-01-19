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
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/jessevdk/go-flags"
)

type options struct {
	Libary  bool `short:"l" long:"library" description:"Generate library project template (default: application project)"`
	NoRoot  bool `short:"n" long:"no-root" description:"Create files in the current directory without creating the project root directory"`
	Version bool `short:"v" long:"version" description:"Show ubume command version"`
}

var osExit = os.Exit

const cmdName string = "ubume"
const version string = "1.2.0"

const (
	exitSuccess int = iota // 0
	exitFailure
)

// project have project information to be generated.
type project struct {
	importPath string // same as "$ git mod init <importPath>"
	name       string // project (command) name
	version    string // project version
	library    bool   // flag that make library project
	noRoot     bool   // flag that don't make project root directory
	makes      func()
}

// main is entry point of ubume command.
func main() {
	var opts options
	args := parseArgs(&opts)
	prj := newProject(args, opts)
	prj.canMake()
	prj.makes()
}

// makeAppPrj generate application project directory and files.
func (p *project) makeAppPrj() {
	p.makeProjectDirs()
	p.makeAppProjectFiles()
	p.goModInit()
}

// makeLibPrj generate library project directory and files.
func (p *project) makeLibPrj() {
	p.makeProjectDirs()
	p.makeLibProjectFiles()
	p.goModInit()
}

// makeProjectDirs create all directories in project template.
// If it can not make directories, exit command.
func (p *project) makeProjectDirs() {
	dirs := []string{}
	if p.noRoot {
		dirs = append(dirs, filepath.Join(".github", "workflows"))
	} else {
		dirs = append(dirs, filepath.Join(p.name, ".github", "workflows"))
	}

	if p.library {
		if !p.noRoot {
			dirs = append(dirs, p.name)
		}
	} else {
		if p.noRoot {
			dirs = append(dirs, filepath.Join("cmd", p.name))
		} else {
			dirs = append(dirs, filepath.Join(p.name, "cmd", p.name))
		}
	}
	mkDirs(dirs)
}

// mkDirs create multiple specified directories.
// If the parent directory does not exist, create the parent directory as well.
// If an error occurs, exit command.
func mkDirs(paths []string) {
	for _, path := range paths {
		target := os.ExpandEnv(path)
		err := os.MkdirAll(target, 0755)
		if err != nil {
			die(err.Error())
		}
	}
}

// makeProjectDirs create all files in application project template.
func (p *project) makeAppProjectFiles() {
	p.makeAppMainSourecCodeFile()
	p.makeTestFile()
	p.makeDocGoFile()
	p.makeMakefileForApp()
	p.makeGitHubActionsFile()
	p.makeChangelogFile()
}

// makeLibProjectFiles create all files in library project template.
func (p *project) makeLibProjectFiles() {
	p.makeLibSourceCodeFile()
	p.makeTestFile()
	p.makeDocGoFile()
	p.makeMakefileForLib()
	p.makeGitHubActionsFile()
	p.makeChangelogFile()
}

// makeAppMainSourecCodeFile create file that is source code for command.
func (p *project) makeAppMainSourecCodeFile() {
	var path string
	if p.noRoot {
		path = filepath.Join("cmd", p.name, "main.go")
	} else {
		path = filepath.Join(p.name, "cmd", p.name, "main.go")
	}

	code := `package main

import "fmt"

func main() {
	fmt.Println(HelloWorld())
}

func HelloWorld() string {
	return "Hello, World"
}
`
	writeFile(code, path)
}

// makeLibSourceCodeFile create file that is source code for library.
func (p *project) makeLibSourceCodeFile() {
	var path string
	if p.noRoot {
		path = filepath.Join(p.name + ".go")
	} else {
		path = filepath.Join(p.name, p.name+".go")
	}

	code := `package XXX_PKG_XXX

func HelloWorld() string {
	return "Hello, World"
}
`
	writeFile(strings.ReplaceAll(code, "XXX_PKG_XXX", p.name), path)
}

// makeTestFile create file that is test source code for command/library.
func (p *project) makeTestFile() {
	code := `package XXX_PKG_XXX

import "testing"

func TestHelloWorld(t *testing.T) {
	if HelloWorld() != "Hello, World" {
		t.Errorf("HelloWorlf = %s, want \"Hello, World\"", HelloWorld())
	}
}
`
	var path string
	if p.library {
		if p.noRoot {
			path = filepath.Join(p.name + "_test.go")
		} else {
			path = filepath.Join(p.name, p.name+"_test.go")
		}
		code = strings.ReplaceAll(code, "XXX_PKG_XXX", p.name)
	} else {
		if p.noRoot {
			path = filepath.Join("cmd", p.name, "main_test.go")
		} else {
			path = filepath.Join(p.name, "cmd", p.name, "main_test.go")
		}
		code = strings.ReplaceAll(code, "XXX_PKG_XXX", "main")
	}
	writeFile(code, path)
}

// makeDocGoFile create doc.go for module description.
func (p *project) makeDocGoFile() {
	code := `// This package is generated by ubume command.
//
// If you publish this module on GitHub etc., please write down the 
// application description in this file. When your package is published
// to "https://pkg.go.dev/", the contents of doc.go will be automatically
// listed as an overview on the pkg.go.dev.
package XXX_PKG_XXX
`
	var path string
	if p.library {
		if p.noRoot {
			path = filepath.Join("doc.go")
		} else {
			path = filepath.Join(p.name, "doc.go")
		}
		code = strings.ReplaceAll(code, "XXX_PKG_XXX", p.name)
	} else {
		if p.noRoot {
			path = filepath.Join("cmd", p.name, "doc.go")
		} else {
			path = filepath.Join(p.name, "cmd", p.name, "doc.go")
		}
		code = strings.ReplaceAll(code, "XXX_PKG_XXX", "main")
	}
	writeFile(code, path)
}

// makeMakefileForApp create Makefile at application project template root directory.
func (p *project) makeMakefileForApp() {
	var path string
	if p.noRoot {
		path = "Makefile"
	} else {
		path = filepath.Join(p.name, "Makefile")
	}

	data := `.PHONY: build test clean vet fmt chkfmt

APP         = XXX_APP_XXX
GO          = go
GO_BUILD    = $(GO) build
GO_FORMAT   = $(GO) fmt
GOFMT       = gofmt
GO_LIST     = $(GO) list
GO_TEST     = $(GO) test -v
GO_TOOL     = $(GO) tool
GO_VET      = $(GO) vet
GO_DEP      = $(GO) mod
GOOS        = XXX_OS_XXX
GO_PKGROOT  = ./...
GO_PACKAGES = $(shell $(GO_LIST) $(GO_PKGROOT))

build:  ## Build binary 
	env GO111MODULE=on GOOS=$(GOOS) $(GO_BUILD) $(GO_LDFLAGS) -o $(APP) XXX_CODE_XXX

clean: ## Clean project
	-rm -rf $(APP) cover.out cover.html

test: ## Start test
	env GOOS=$(GOOS) $(GO_TEST) -cover $(GO_PKGROOT) -coverprofile=cover.out
	$(GO_TOOL) cover -html=cover.out -o cover.html

vet: ## Start go vet
	$(GO_VET) $(GO_PACKAGES)

fmt: ## Format go source code 
	$(GO_FORMAT) $(GO_PKGROOT)

.DEFAULT_GOAL := help
help:  
	@grep -E '^[0-9a-zA-Z_-]+[[:blank:]]*:.*?## .*$$' $(MAKEFILE_LIST) | sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1;32m%-15s\033[0m %s\n", $$1, $$2}'
`
	data = strings.Replace(data, "XXX_APP_XXX", p.name, 1)
	data = strings.Replace(data, "XXX_OS_XXX", runtime.GOOS, 1)
	data = strings.Replace(data, "XXX_CODE_XXX", filepath.Join("cmd", p.name, "main.go"), 1)
	writeFile(data, path)
}

// makeMakefileForLib create Makefile at application project template root directory.
func (p *project) makeMakefileForLib() {
	var path string
	if p.noRoot {
		path = "Makefile"
	} else {
		path = filepath.Join(p.name, "Makefile")
	}
	data := `.PHONY: build test clean vet fmt chkfmt

APP         = XXX_APP_XXX
GO          = go
GO_BUILD    = $(GO) build
GO_FORMAT   = $(GO) fmt
GOFMT       = gofmt
GO_LIST     = $(GO) list
GO_TEST     = $(GO) test -v
GO_TOOL     = $(GO) tool
GO_VET      = $(GO) vet
GO_DEP      = $(GO) mod
GO_LDFLAGS  = -ldflags="-s -w"
GOOS        = XXX_OS_XXX
GO_PKGROOT  = ./...
GO_PACKAGES = $(shell $(GO_LIST) $(GO_PKGROOT))

test: ## Start test
	env GOOS=$(GOOS) $(GO_TEST) -cover $(GO_PKGROOT) -coverprofile=cover.out
	$(GO_TOOL) cover -html=cover.out -o cover.html

vet: ## Start go vet
	$(GO_VET) $(GO_PACKAGES)

fmt: ## Format go source code 
	$(GO_FORMAT) $(GO_PKGROOT)

.DEFAULT_GOAL := help
help:  
	@grep -E '^[0-9a-zA-Z_-]+[[:blank:]]*:.*?## .*$$' $(MAKEFILE_LIST) | sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1;32m%-15s\033[0m %s\n", $$1, $$2}'
`
	data = strings.Replace(data, "XXX_APP_XXX", p.name, 1)
	data = strings.Replace(data, "XXX_OS_XXX", runtime.GOOS, 1)
	data = strings.Replace(data, "XXX_CODE_XXX", filepath.Join("cmd", p.name, "main.go"), 1)
	writeFile(data, path)
}

// makeChangelogFile create CHAGELOG.md at project template root directory.
func (p *project) makeChangelogFile() {
	var path string
	if p.noRoot {
		path = "Changelog.md"
	} else {
		path = filepath.Join(p.name, "Changelog.md")
	}

	data := `# Changelog
All notable changes to this project will be documented in this file.  
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).   
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
`
	writeFile(data, path)
}

// makeGitHubActions make build.yml and test.yml for GitHub Actions workflows.
func (p *project) makeGitHubActionsFile() {
	if !p.library {
		p.makeBuildYml()
	}
	p.makeUnitTestYml()
}

func (p *project) makeBuildYml() {
	var path string
	if p.noRoot {
		path = filepath.Join(".github", "workflows", "build.yml")
	} else {
		path = filepath.Join(p.name, ".github", "workflows", "build.yml")
	}
	data := `name: Build

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:

    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2

    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: XXX_VER_XXX

    - name: Build
      run: make build
`
	data = strings.Replace(data, "XXX_VER_XXX", goVersion(), 1)
	writeFile(data, path)
}

func (p *project) makeUnitTestYml() {
	var path string
	if p.noRoot {
		path = filepath.Join(".github", "workflows", "unit_test.yml")
	} else {
		path = filepath.Join(p.name, ".github", "workflows", "unit_test.yml")
	}
	data := `name: UnitTest

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  unit-test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2

    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.17

    - name: UnitTest
      run: make test
`
	data = strings.Replace(data, "XXX_VER_XXX", goVersion(), 1)
	writeFile(data, path)
}

// goVersion return runtime golang version(only number, not include "go" or "cpu name")
func goVersion() string {
	runtimeVer := runtime.Version()

	rex := regexp.MustCompile("[0-9]\\.[0-9]*")
	return rex.FindString(runtimeVer)
}

// goModInit execute "$ go mod init <importPath>"
// If it can not execute "$ go mod", exit command.
func (p *project) goModInit() {
	if !p.noRoot {
		err := os.Chdir(p.name)
		if err != nil {
			die(err.Error())
		}
	}

	if err := exec.Command("go", "mod", "init", p.importPath).Run(); err != nil {
		die(err.Error())
	}
}

// canMake check whether can create project template or not.
// If it can't create the project, exit command.
func (p *project) canMake() {
	p.canUseGoCmd()
	p.canMakePrjDir()
	if p.noRoot {
		p.canMakePrjFile()
	}
}

// canUseGoCmd check whether go command install in the system.
// If not install, exit command.
func (p *project) canUseGoCmd() {
	_, err := exec.LookPath("go")
	if err != nil {
		die("this system does not install go cmd. Please download golang")
	}
}

// canMakePrjDir exit the command if there is a directory with the same name or
// if the project name is an empty string.
func (p *project) canMakePrjDir() {
	if exists(p.name) {
		die("same name project already exists at current directory")
	}
	if p.name == "" {
		die("project name is empty (import path end with \"/\"?)")
	}
}

// canMakePrjFile check whether the file ubume is trying to generate already exists.
func (p *project) canMakePrjFile() {
	files := []string{
		// for library
		p.name + ".go",
		p.name + "_test.go",
		// for app
		filepath.Join("cmd", p.name, "main.go"),
		filepath.Join("cmd", p.name, "main_test.go"),
		filepath.Join("cmd", p.name, "doc.go"),
		"Makefile",
		"Changelog.md",
		"go.mod",
		".github",
	}
	for _, v := range files {
		if exists(v) {
			die("same name file (" + v + ") already exists at current directory")
		}
	}
}

// die exit program with message.
func die(msg string) {
	fmt.Fprintln(os.Stderr, cmdName+": "+msg)
	osExit(exitFailure)
}

// exists check whether file or directory exists.
func exists(path string) bool {
	_, err := os.Stat(path)
	return (err == nil)
}

// writeFile write string to file.
// If it can not create file, exit command.
func writeFile(text string, path string) {
	file, err := os.Create(path)
	if err != nil {
		die(err.Error())
	}
	defer file.Close()
	file.Write(([]byte)(text))
}

// newProject return initialized project struct.
func newProject(arg []string, opts options) *project {
	var prj project
	prj.importPath = arg[0]

	if strings.Contains(prj.importPath, "/") {
		arr := strings.Split(prj.importPath, "/")
		prj.name = arr[len(arr)-1]
	} else {
		prj.name = prj.importPath
	}
	prj.version = "0.0.1"

	if opts.Libary {
		prj.library = true
		prj.makes = prj.makeLibPrj
	} else {
		prj.library = false
		prj.makes = prj.makeAppPrj
	}

	prj.noRoot = opts.NoRoot
	return &prj
}

// parseArgs parse command line arguments.
// In this method, process for version option, help option, and lack of arguments.
func parseArgs(opts *options) []string {
	p := newParser(opts)

	args, err := p.Parse()
	if err != nil {
		// If user specify --help option, help message already showed in p.Parse().
		// Moreover, help messages are stored inside err.
		if !strings.Contains(err.Error(), cmdName) {
			showHelp(p)
			showHelpFooter()
		} else {
			showHelpFooter()
		}
		osExit(exitFailure)
	}

	if opts.Version {
		showVersion(cmdName, version)
		osExit(exitSuccess)
	}

	if len(args) != 1 {
		showHelp(p)
		showHelpFooter()
		osExit(exitFailure)
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
	parser.Name = cmdName
	parser.Usage = "[OPTIONS] IMPORT_PATH  ※ IMPORT_PATH is same as $ go mod init IMPORT_PATH"
	return parser
}

// showVersion show ubume command version information.
func showVersion(cmdName string, version string) {
	description := cmdName + " version " + version + " (under Apache License verison 2.0)"
	fmt.Fprintln(os.Stdout, description)
}
