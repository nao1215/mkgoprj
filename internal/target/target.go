// Package target handles information about directories and files to be generated.
package target

import (
	"path/filepath"
	"strings"

	"github.com/nao1215/mkgoprj/v2/internal/gotool"
)

// Dirs returns the directory to be created.
// name   : Project name
// lib    : Whether to create library project
// noRoot : Whether to create the project root directory (project name directory)
func Dirs(name string, lib, cli, noRoot bool) []string {
	dirs := []string{}
	if noRoot {
		dirs = append(dirs, filepath.Join(".github", "workflows"))
		dirs = append(dirs, filepath.Join(".github", "ISSUE_TEMPLATE"))
	} else {
		dirs = append(dirs, filepath.Join(name, ".github", "workflows"))
		dirs = append(dirs, filepath.Join(name, ".github", "ISSUE_TEMPLATE"))
	}

	if lib {
		if !noRoot {
			dirs = append(dirs, name)
		}
	} else if cli {
		if noRoot {
			dirs = append(dirs, "cmd")
			dirs = append(dirs, filepath.Join("internal", "print"))
		} else {
			dirs = append(dirs, filepath.Join(name, "cmd"))
			dirs = append(dirs, filepath.Join(name, "internal", "print"))
		}
	}
	return dirs
}

// Files returns the directory to be created.
func Files(name, importPath string, lib, cli, noRoot bool) map[string]string {
	files := map[string]string{}

	if lib {
		path, code := librarySourceCodeFile(name, noRoot)
		files[path] = code
	} else if cli {
		path, code := cliMainSourceCodeFile(name, importPath, noRoot)
		files[path] = code
	}

	if !cli {
		path, code := mainTestFile(name, lib, noRoot)
		files[path] = code
	}

	if cli {
		path, code := rootFile(name, importPath, noRoot)
		files[path] = code
		path, code = versionFile(name, noRoot)
		files[path] = code
		path, code = printFile(name, importPath, noRoot)
		files[path] = code
		path, code = printTestFile(name, noRoot)
		files[path] = code
	}

	path, code := makefile(name, importPath, lib, cli, noRoot)
	files[path] = code

	path, code = changelogFile(name, noRoot)
	files[path] = code

	if !lib {
		path, code = githubBuildYml(name, noRoot)
		files[path] = code
		path, code = githubRelease(name, noRoot)
		files[path] = code
		path, code = goreleaser(name, importPath, noRoot, cli)
		files[path] = code
	}

	// contributor command has many bugs. Not use it.
	//path, code = githubContributors(name, noRoot)
	//files[path] = code

	path, code = githubUnitTestYml(name, noRoot)
	files[path] = code

	path, code = githubReviewDog(name, noRoot)
	files[path] = code

	path, code = codeOfConduct(name, noRoot)
	files[path] = code

	path, code = dependabot(name, noRoot)
	files[path] = code

	path, code = bugReportTemplate(name, noRoot)
	files[path] = code

	path, code = issueTemplate(name, noRoot)
	files[path] = code

	return files
}

func cliMainSourceCodeFile(name, importPath string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = "main.go"
	} else {
		path = filepath.Join(name, "main.go")
	}

	code := `package main

import "XXX_IMPORT_PATH_XXX"

func main() {
	cmd.Execute()
}
`
	return path, strings.ReplaceAll(code, "XXX_IMPORT_PATH_XXX", filepath.Join(importPath, "cmd"))
}

func librarySourceCodeFile(name string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join(name + ".go")
	} else {
		path = filepath.Join(name, name+".go")
	}

	code := `package XXX_PKG_XXX

func HelloWorld() string {
	return "Hello, World"
}
`
	return path, strings.ReplaceAll(code, "XXX_PKG_XXX", name)
}

func mainTestFile(name string, libProject, noRoot bool) (string, string) {
	code := `package XXX_PKG_XXX

import "testing"
	
func TestHelloWorld(t *testing.T) {
	if HelloWorld() != "Hello, World" {
		t.Errorf("HelloWorlf = %s, want \"Hello, World\"", HelloWorld())
	}
}
	`
	var path string
	if libProject {
		if noRoot {
			path = filepath.Join(name + "_test.go")
		} else {
			path = filepath.Join(name, name+"_test.go")
		}
		code = strings.ReplaceAll(code, "XXX_PKG_XXX", name)
	}
	return path, code
}

func makefile(name, importPath string, libProject, cli, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = "Makefile"
	} else {
		path = filepath.Join(name, "Makefile")
	}

	code := `.PHONY: build test clean vet fmt chkfmt

APP         = XXX_APP_XXX
VERSION     = $(shell git describe --tags --abbrev=0)
GO          = go
GO_BUILD    = $(GO) build
GO_FORMAT   = $(GO) fmt
GOFMT       = gofmt
GO_LIST     = $(GO) list
GO_TEST     = $(GO) test -v
GO_TOOL     = $(GO) tool
GO_VET      = $(GO) vet
GO_DEP      = $(GO) mod
GOOS        = ""
GOARCH      = ""
GO_PKGROOT  = ./...
GO_PACKAGES = $(shell $(GO_LIST) $(GO_PKGROOT))
GO_LDFLAGS  = -ldflags '-X XXX_IMPORT_PATH_XXX/cmd.Version=${VERSION}'

XXX_ONLY_APP_XXX

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

	strOnlyApp := `build:  ## Build binary
	env GO111MODULE=on GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BUILD) $(GO_LDFLAGS) -o $(APP) XXX_CODE_XXX`

	if libProject {
		code = strings.Replace(code, "XXX_ONLY_APP_XXX", "", 1)
	} else if cli {
		code = strings.Replace(code, "XXX_ONLY_APP_XXX", strOnlyApp, 1)
		code = strings.Replace(code, "XXX_CODE_XXX", filepath.Join("main.go"), 1)
	} else {
		code = strings.Replace(code, "XXX_ONLY_APP_XXX", strOnlyApp, 1)
		code = strings.Replace(code, "XXX_CODE_XXX", filepath.Join("cmd", name, "main.go"), 1)
	}
	code = strings.Replace(code, "XXX_APP_XXX", name, 1)
	code = strings.Replace(code, "XXX_IMPORT_PATH_XXX", importPath, 1)
	return path, code
}

func changelogFile(name string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = "Changelog.md"
	} else {
		path = filepath.Join(name, "Changelog.md")
	}

	data := `# Changelog
All notable changes to this project will be documented in this file.  
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).   
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
`
	return path, data
}

func githubBuildYml(name string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join(".github", "workflows", "build.yml")
	} else {
		path = filepath.Join(name, ".github", "workflows", "build.yml")
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
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: "XXX_VER_XXX"

    - name: Build
      run: make build
`
	data = strings.Replace(data, "XXX_VER_XXX", gotool.Version(), 1)
	return path, data
}

func githubUnitTestYml(name string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join(".github", "workflows", "unit_test.yml")
	} else {
		path = filepath.Join(name, ".github", "workflows", "unit_test.yml")
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
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: "XXX_VER_XXX"

    - name: UnitTest
      run: make test
`
	data = strings.Replace(data, "XXX_VER_XXX", gotool.Version(), 1)
	return path, data
}

func githubReviewDog(name string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join(".github", "workflows", "reviewdog.yml")
	} else {
		path = filepath.Join(name, ".github", "workflows", "reviewdog.yml")
	}
	data := `name: reviewdog
on: [pull_request]

jobs:
  golangci-lint:
    name: golangci-lint
    runs-on: ubuntu-latest
    steps:
      - name: Check out code into the Go module directory
        uses: actions/checkout@v3
        with:
          persist-credentials: false
      - name: golangci-lint
        uses: reviewdog/action-golangci-lint@v2
        with:
          reporter: github-pr-review
          level: warning

  misspell:
    name: misspell
    runs-on: ubuntu-latest
    steps:
      - name: Check out code into the Go module directory
        uses: actions/checkout@v3
        with:
          persist-credentials: false
      - name: misspell
        uses: reviewdog/action-misspell@v1
        with:
          reporter: github-pr-review
          level: warning
          locale: "US"

  actionlint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: reviewdog/action-actionlint@v1
        with:
          reporter: github-pr-review
          level: warning
`
	return path, data
}

func githubContributors(name string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join(".github", "workflows", "contributors.yml")
	} else {
		path = filepath.Join(name, ".github", "workflows", "contributors.yml")
	}
	data := `name: Contributors

on:
  pull_request:
    branches: [main]

jobs:
  contributors:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: "XXX_VER_XXX"

      - name: Generate Contributors
        run: |
          go install github.com/nao1215/contributor@latest
          git remote set-url origin https://github-actions:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}
          git config --global user.name "${GITHUB_ACTOR}"
          git config --global user.email "${GITHUB_ACTOR}@users.noreply.github.com"
          contributor --file
          git add .
          git commit -m "Update Contributors List"
          git push origin HEAD:${GITHUB_REF}
`
	data = strings.Replace(data, "XXX_VER_XXX", gotool.Version(), 1)
	return path, data
}

func githubRelease(name string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join(".github", "workflows", "release.yml")
	} else {
		path = filepath.Join(name, ".github", "workflows", "release.yml")
	}
	data := `name: Release

on:
  push:
    tags:
      - "v*"

jobs:
  release:
    name: Release
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
        with:
          fetch-depth: 0
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: "XXX_VER_XXX"
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v2
        with:
          version: latest
          args: release --rm-dist
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
`
	data = strings.Replace(data, "XXX_VER_XXX", gotool.Version(), 1)
	return path, data
}

func goreleaser(name, importPath string, noRoot, cli bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join(".goreleaser.yml")
	} else {
		path = filepath.Join(name, ".goreleaser.yml")
	}
	data := `project_name: XXX_APP_NAME_XXX
env:
  - GO111MODULE=on
before:
  hooks:
    - go mod tidy
    - go generate ./...
builds:
  - main: XXX_BUILD_TARGET_XXX
    ldflags:
      - -s -w -X XXX_IMPORT_PATH_XXX/cmd.Version=v{{ .Version }}
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
archives:
  - name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    replacements:
      darwin: Darwin
      linux: Linux
      windows: Windows
      386: i386
      amd64: x86_64
    format_overrides:
      - goos: windows
        format: zip
checksum:
  name_template: "checksums.txt"
snapshot:
  name_template: "{{ incpatch .Version }}-next"
changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"	
`
	data = strings.Replace(data, "XXX_APP_NAME_XXX", name, 1)
	data = strings.Replace(data, "XXX_BUILD_TARGET_XXX", ".", 1)
	data = strings.Replace(data, "XXX_IMPORT_PATH_XXX", name, 1)
	return path, data
}

func rootFile(name, importPath string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join("cmd", "root.go")
	} else {
		path = filepath.Join(name, "cmd", "root.go")
	}
	data := `package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"XXX_PATH_XXX/internal/print"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "XXX_CMD_XXX",
}

// Execute start command.
func Execute() {
	if isWindows() {
		print.Err("not support windows")
		os.Exit(1)
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SilenceErrors = true
	deployShellCompletionFileIfNeeded(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

// deployShellCompletionFileIfNeeded creates the shell completion file.
// If the file with the same contents already exists, it is not created.
func deployShellCompletionFileIfNeeded(cmd *cobra.Command) {
	if !isWindows() {
		makeBashCompletionFileIfNeeded(cmd)
		makeFishCompletionFileIfNeeded(cmd)
		makeZshCompletionFileIfNeeded(cmd)
	}
}

func makeBashCompletionFileIfNeeded(cmd *cobra.Command) {
	if existSameBashCompletionFile(cmd) {
		return
	}

	path := bashCompletionFilePath()
	bashCompletion := new(bytes.Buffer)
	if err := cmd.GenBashCompletion(bashCompletion); err != nil {
		print.Err(fmt.Errorf("can not generate bash completion content: %w", err))
		return
	}

	if !isFile(path) {
		fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0664)
		if err != nil {
			print.Err(fmt.Errorf("can not open .bash_completion: %w", err))
			return
		}

		if _, err := fp.WriteString(bashCompletion.String()); err != nil {
			print.Err(fmt.Errorf("can not write .bash_completion %w", err))
			return
		}

		if err := fp.Close(); err != nil {
			print.Err(fmt.Errorf("can not close .bash_completion %w", err))
			return
		}
		return
	}

	fp, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0664)
	if err != nil {
		print.Err(fmt.Errorf("can not append .bash_completion: %w", err))
		return
	}

	if _, err := fp.WriteString(bashCompletion.String()); err != nil {
		print.Err(fmt.Errorf("can not write .bash_completion: %w", err))
		return
	}

	if err := fp.Close(); err != nil {
		print.Err(fmt.Errorf("can not close .bash_completion: %w", err))
		return
	}
	print.Info("append bash-completion: " + path)
}

func makeFishCompletionFileIfNeeded(cmd *cobra.Command) {
	if isSameFishCompletionFile(cmd) {
		return
	}

	path := fishCompletionFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0775); err != nil {
		print.Err(fmt.Errorf("can not create fish-completion file: %w", err))
		return
	}

	if err := cmd.GenFishCompletionFile(path, false); err != nil {
		print.Err(fmt.Errorf("can not create fish-completion file: %w", err))
		return
	}
	print.Info("create fish-completion file: " + path)
}

func makeZshCompletionFileIfNeeded(cmd *cobra.Command) {
	if isSameZshCompletionFile(cmd) {
		return
	}

	path := zshCompletionFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0775); err != nil {
		print.Err(fmt.Errorf("can not create zsh-completion file: %w", err))
		return
	}

	if err := cmd.GenZshCompletionFile(path); err != nil {
		print.Err(fmt.Errorf("can not create zsh-completion file: %w", err))
		return
	}
	print.Info("create zsh-completion file: " + path)

	appendFpathAtZshrcIfNeeded()
}

func appendFpathAtZshrcIfNeeded() {
	const zshFpath = XXX_ZSH_FPATH_XXX

	zshrcPath := zshrcPath()
	if !isFile(zshrcPath) {
		fp, err := os.OpenFile(zshrcPath, os.O_RDWR|os.O_CREATE, 0664)
		if err != nil {
			print.Err(fmt.Errorf("can not open .zshrc: %w", err).Error())
			return
		}

		if _, err := fp.WriteString(zshFpath); err != nil {
			print.Err(fmt.Errorf("can not write zsh $fpath in .zshrc: %w", err).Error())
			return
		}

		if err := fp.Close(); err != nil {
			print.Err(fmt.Errorf("can not close .zshrc: %w", err).Error())
			return
		}
		return
	}

	zshrc, err := os.ReadFile(zshrcPath)
	if err != nil {
		print.Err(fmt.Errorf("can not read .zshrc: %w", err).Error())
		return
	}

	if strings.Contains(string(zshrc), zshFpath) {
		return
	}

	fp, err := os.OpenFile(zshrcPath, os.O_RDWR|os.O_APPEND, 0664)
	if err != nil {
		print.Err(fmt.Errorf("can not open .zshrc: %w", err).Error())
		return
	}

	if _, err := fp.WriteString(zshFpath); err != nil {
		print.Err(fmt.Errorf("can not write zsh $fpath in .zshrc: %w", err).Error())
		return
	}

	if err := fp.Close(); err != nil {
		print.Err(fmt.Errorf("can not close .zshrc: %w", err).Error())
		return
	}
}

func existSameBashCompletionFile(cmd *cobra.Command) bool {
	if !isFile(bashCompletionFilePath()) {
		return false
	}
	return hasSameBashCompletionContent(cmd)
}

func hasSameBashCompletionContent(cmd *cobra.Command) bool {
	bashCompletionFileInLocal, err := os.ReadFile(bashCompletionFilePath())
	if err != nil {
		print.Err(fmt.Errorf("can not read .bash_completion: %w", err).Error())
		return false
	}

	currentBashCompletion := new(bytes.Buffer)
	if err := cmd.GenBashCompletion(currentBashCompletion); err != nil {
		return false
	}
	if !strings.Contains(string(bashCompletionFileInLocal), currentBashCompletion.String()) {
		return false
	}
	return true
}

func isSameFishCompletionFile(cmd *cobra.Command) bool {
	path := fishCompletionFilePath()
	if !isFile(path) {
		return false
	}

	currentFishCompletion := new(bytes.Buffer)
	if err := cmd.GenFishCompletion(currentFishCompletion, false); err != nil {
		return false
	}

	fishCompletionInLocal, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	if bytes.Compare(currentFishCompletion.Bytes(), fishCompletionInLocal) != 0 {
		return false
	}
	return true
}

func isSameZshCompletionFile(cmd *cobra.Command) bool {
	path := zshCompletionFilePath()
	if !isFile(path) {
		return false
	}

	currentZshCompletion := new(bytes.Buffer)
	if err := cmd.GenZshCompletion(currentZshCompletion); err != nil {
		return false
	}

	zshCompletionInLocal, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	if bytes.Compare(currentZshCompletion.Bytes(), zshCompletionInLocal) != 0 {
		return false
	}
	return true
}

// bashCompletionFilePath return bash-completion file path.
func bashCompletionFilePath() string {
	return filepath.Join(os.Getenv("HOME"), ".bash_completion")
}

// fishCompletionFilePath return fish-completion file path.
func fishCompletionFilePath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "fish", "completions", Name+".fish")
}

// zshCompletionFilePath return zsh-completion file path.
func zshCompletionFilePath() string {
	return filepath.Join(os.Getenv("HOME"), ".zsh", "completion", "_"+Name)
}

// zshrcPath return .zshrc path.
func zshrcPath() string {
	return filepath.Join(os.Getenv("HOME"), ".zshrc")
}

// isFile reports whether the path exists and is a file.
func isFile(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil) && (!stat.IsDir())
}
`

	zshFpath := "`"
	zshFpath += `
# setting for XXX_NAME_XXX command (auto generate)
fpath=(~/.zsh/completion $fpath)
autoload -Uz compinit && compinit -i
`
	zshFpath += "`"

	data = strings.Replace(data, "XXX_PATH_XXX", importPath, -1)
	data = strings.Replace(data, "XXX_NAME_XXX", name, -1)
	data = strings.Replace(data, "XXX_ZSH_FPATH_XXX", zshFpath, 1)

	return path, data
}

func versionFile(name string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join("cmd", "version.go")
	} else {
		path = filepath.Join(name, "cmd", "version.go")
	}
	data := `package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show " + Name + " command version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// Version value is set by s
var Version string

// Name is command name
const Name = "morrigan"

// getVersion return gup command version.
// Version global variable is set by s.
func getVersion() string {
	version := "unknown"
	if Version != "" {
		version = Version
	} else if buildInfo, ok := debug.ReadBuildInfo(); ok {
		version = buildInfo.Main.Version
	}
	return fmt.Sprintf("%s version %s", Name, version)
}
`
	return path, data
}

func printFile(name, importPath string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join("internal", "print", "print.go")
	} else {
		path = filepath.Join(name, "internal", "print", "print.go")
	}
	data := `// Package print defines functions to accept colored standard output and user input
package print

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var (
	// Stdout is new instance of Writer which handles escape sequence for stdout.
	Stdout = colorable.NewColorableStdout()
	// Stderr is new instance of Writer which handles escape sequence for stderr.
	Stderr = colorable.NewColorableStderr()
)

// Info print information message at STDOUT in green.
// This function is used to print some information (that is not error) to the user.
func Info(msg string) {
	fmt.Fprintf(Stdout, "%s: %s\n", color.GreenString("INFO "), msg)
}

// Warn print warning message at STDERR in yellow.
// This function is used to print warning message to the user.
func Warn(err interface{}) {
	fmt.Fprintf(Stderr, "%s: %v\n", color.YellowString("WARN "), err)
}

// Err print error message at STDERR in yellow.
// This function is used to print error message to the user.
func Err(err interface{}) {
	fmt.Fprintf(Stderr, "%s: %v\n", color.HiYellowString("ERROR"), err)
}

// OsExit is wrapper for  os.Exit(). It's for unit test.
var OsExit = os.Exit

// Fatal print dying message at STDERR in red.
// After print message, process will exit
func Fatal(err interface{}) {
	fmt.Fprintf(Stderr, "%s: %v\n", color.RedString("FATAL"), err)
	OsExit(1)
}

// FmtScanln is wrapper for fmt.Scanln(). It's for unit test.
var FmtScanln = fmt.Scanln

// Question displays the question in the terminal and receives an answer from the user.
func Question(ask string) bool {
	var response string

	fmt.Fprintf(Stdout, "%s: %s", color.GreenString("CHECK"), ask+" [Y/n] ")
	_, err := FmtScanln(&response)
	if err != nil {
		// If user input only enter.
		if strings.Contains(err.Error(), "expected newline") {
			return Question(ask)
		}
		fmt.Fprint(os.Stderr, err.Error())
		return false
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return Question(ask)
	}
}
`
	data = strings.Replace(data, "XXX_PATH_XXX", importPath, -1)
	return path, data
}

func printTestFile(name string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join("internal", "print", "print_test.go")
	} else {
		path = filepath.Join(name, "internal", "print", "print_test.go")
	}
	data := `// Package print defines functions to accept colored standard output and user input
package print

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInfo(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Print message",
			args: args{
				msg: "test message",
			},
			want: []string{"INFO : test message", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgStdout := Stdout
			orgStderr := Stderr
			pr, pw, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			Stdout = pw
			Stderr = pw

			Info(tt.args.msg)
			pw.Close()
			Stdout = orgStdout
			Stderr = orgStderr

			if err != nil {
				return
			}

			buf := bytes.Buffer{}
			_, err = io.Copy(&buf, pr)
			if err != nil {
				t.Error(err)
			}
			got := strings.Split(buf.String(), "\n")

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Print message",
			args: args{
				msg: "test message",
			},
			want: []string{"WARN : test message", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgStdout := Stdout
			orgStderr := Stderr
			pr, pw, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			Stdout = pw
			Stderr = pw

			Warn(tt.args.msg)
			pw.Close()
			Stdout = orgStdout
			Stderr = orgStderr

			if err != nil {
				return
			}

			buf := bytes.Buffer{}
			_, err = io.Copy(&buf, pr)
			if err != nil {
				t.Error(err)
			}
			got := strings.Split(buf.String(), "\n")

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestErr(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Print message",
			args: args{
				msg: "test message",
			},
			want: []string{"ERROR: test message", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgStdout := Stdout
			orgStderr := Stderr
			pr, pw, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			Stdout = pw
			Stderr = pw

			Err(tt.args.msg)
			pw.Close()
			Stdout = orgStdout
			Stderr = orgStderr

			if err != nil {
				return
			}

			buf := bytes.Buffer{}
			_, err = io.Copy(&buf, pr)
			if err != nil {
				t.Error(err)
			}
			got := strings.Split(buf.String(), "\n")

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFatal(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name     string
		args     args
		want     []string
		exitcode int
	}{
		{
			name: "Print message",
			args: args{
				msg: "test message",
			},
			want:     []string{"FATAL: test message", ""},
			exitcode: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgStdout := Stdout
			orgStderr := Stderr
			pr, pw, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			Stdout = pw
			Stderr = pw

			orgOsExit := OsExit
			exitCode := 0
			OsExit = func(code int) {
				exitCode = code
			}
			defer func() { OsExit = orgOsExit }()

			Fatal(tt.args.msg)
			pw.Close()
			Stdout = orgStdout
			Stderr = orgStderr

			if err != nil {
				return
			}

			buf := bytes.Buffer{}
			_, err = io.Copy(&buf, pr)
			if err != nil {
				t.Error(err)
			}
			got := strings.Split(buf.String(), "\n")

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}

			if exitCode != tt.exitcode {
				t.Errorf("value is mismatch. want=%d got=%d", exitCode, tt.exitcode)
			}
		})
	}
}

func TestQuestion(t *testing.T) {
	type args struct {
		ask string
	}
	tests := []struct {
		name  string
		args  args
		input string
		want  bool
	}{
		{
			name:  "user input 'y'",
			args:  args{"no check"},
			input: "y",
			want:  true,
		},
		{
			name:  "user input 'yes'",
			args:  args{"no check"},
			input: "yes",
			want:  true,
		},
		{
			name:  "user input 'n'",
			args:  args{"no check"},
			input: "n",
			want:  false,
		},
		{
			name:  "user input 'no'",
			args:  args{"no check"},
			input: "no",
			want:  false,
		},
		{
			name:  "user input 'yes' after 'a'",
			args:  args{"no check"},
			input: "a\nyes",
			want:  true,
		},
		{
			name:  "user only input enter",
			args:  args{"no check"},
			input: "\nyes",
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funcDefer, err := mockStdin(t, tt.input)
			if err != nil {
				t.Fatal(err)
			}
			defer funcDefer()

			if got := Question(tt.args.ask); got != tt.want {
				t.Errorf("Question() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestion_FmtScanlnErr(t *testing.T) {
	t.Run("fmt.Scanln() return error", func(t *testing.T) {
		orgFmtScanln := FmtScanln
		FmtScanln = func(a ...any) (n int, err error) {
			return -1, errors.New("some error")
		}
		defer func() { FmtScanln = orgFmtScanln }()

		if got := Question("no check"); got != false {
			t.Errorf("Question() = %v, want %v", got, false)
		}
	})
}

// mockStdin is a helper function that lets the test pretend dummyInput as os.Stdin.
// It will return a function for defer to clean up after the test.
func mockStdin(t *testing.T, dummyInput string) (funcDefer func(), err error) {
	t.Helper()

	oldOsStdin := os.Stdin
	tmpFile, err := os.CreateTemp(t.TempDir(), "morrigan_")

	if err != nil {
		return nil, err
	}

	content := []byte(dummyInput)

	if _, err := tmpFile.Write(content); err != nil {
		return nil, err
	}

	if _, err := tmpFile.Seek(0, 0); err != nil {
		return nil, err
	}

	// Set stdin to the temp file
	os.Stdin = tmpFile

	return func() {
		// clean up
		os.Stdin = oldOsStdin
		os.Remove(tmpFile.Name())
	}, nil
}
`
	return path, data
}

func codeOfConduct(name string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join("CODE_OF_CONDUCT.md")
	} else {
		path = filepath.Join(name, "CODE_OF_CONDUCT.md")
	}
	data := `# Contributor Covenant Code of Conduct

## Our Pledge

We as members, contributors, and leaders pledge to make participation in our
community a harassment-free experience for everyone, regardless of age, body
size, visible or invisible disability, ethnicity, sex characteristics, gender
identity and expression, level of experience, education, socio-economic status,
nationality, personal appearance, race, religion, or sexual identity
and orientation.

We pledge to act and interact in ways that contribute to an open, welcoming,
diverse, inclusive, and healthy community.

## Our Standards

Examples of behavior that contributes to a positive environment for our
community include:

* Demonstrating empathy and kindness toward other people
* Being respectful of differing opinions, viewpoints, and experiences
* Giving and gracefully accepting constructive feedback
* Accepting responsibility and apologizing to those affected by our mistakes,
  and learning from the experience
* Focusing on what is best not just for us as individuals, but for the
  overall community

Examples of unacceptable behavior include:

* The use of sexualized language or imagery, and sexual attention or
  advances of any kind
* Trolling, insulting or derogatory comments, and personal or political attacks
* Public or private harassment
* Publishing others' private information, such as a physical or email
  address, without their explicit permission
* Other conduct which could reasonably be considered inappropriate in a
  professional setting

## Enforcement Responsibilities

Community leaders are responsible for clarifying and enforcing our standards of
acceptable behavior and will take appropriate and fair corrective action in
response to any behavior that they deem inappropriate, threatening, offensive,
or harmful.

Community leaders have the right and responsibility to remove, edit, or reject
comments, commits, code, wiki edits, issues, and other contributions that are
not aligned to this Code of Conduct, and will communicate reasons for moderation
decisions when appropriate.

## Scope

This Code of Conduct applies within all community spaces, and also applies when
an individual is officially representing the community in public spaces.
Examples of representing our community include using an official e-mail address,
posting via an official social media account, or acting as an appointed
representative at an online or offline event.

## Enforcement

Instances of abusive, harassing, or otherwise unacceptable behavior may be
reported to the community leaders responsible for enforcement at GitHub Issue.
All complaints will be reviewed and investigated promptly and fairly.

All community leaders are obligated to respect the privacy and security of the
reporter of any incident.

## Enforcement Guidelines

Community leaders will follow these Community Impact Guidelines in determining
the consequences for any action they deem in violation of this Code of Conduct:

### 1. Correction

**Community Impact**: Use of inappropriate language or other behavior deemed
unprofessional or unwelcome in the community.

**Consequence**: A private, written warning from community leaders, providing
clarity around the nature of the violation and an explanation of why the
behavior was inappropriate. A public apology may be requested.

### 2. Warning

**Community Impact**: A violation through a single incident or series
of actions.

**Consequence**: A warning with consequences for continued behavior. No
interaction with the people involved, including unsolicited interaction with
those enforcing the Code of Conduct, for a specified period of time. This
includes avoiding interactions in community spaces as well as external channels
like social media. Violating these terms may lead to a temporary or
permanent ban.

### 3. Temporary Ban

**Community Impact**: A serious violation of community standards, including
sustained inappropriate behavior.

**Consequence**: A temporary ban from any sort of interaction or public
communication with the community for a specified period of time. No public or
private interaction with the people involved, including unsolicited interaction
with those enforcing the Code of Conduct, is allowed during this period.
Violating these terms may lead to a permanent ban.

### 4. Permanent Ban

**Community Impact**: Demonstrating a pattern of violation of community
standards, including sustained inappropriate behavior,  harassment of an
individual, or aggression toward or disparagement of classes of individuals.

**Consequence**: A permanent ban from any sort of public interaction within
the community.

## Attribution

This Code of Conduct is adapted from the [Contributor Covenant][homepage],
version 2.0, available at
https://www.contributor-covenant.org/version/2/0/code_of_conduct.html.

Community Impact Guidelines were inspired by [Mozilla's code of conduct
enforcement ladder](https://github.com/mozilla/diversity).

[homepage]: https://www.contributor-covenant.org

For answers to common questions about this code of conduct, see the FAQ at
https://www.contributor-covenant.org/faq. Translations are available at
https://www.contributor-covenant.org/translations.
`
	return path, data
}

func dependabot(name string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join(".github", "dependabot.yml")
	} else {
		path = filepath.Join(name, ".github", "dependabot.yml")
	}
	data := `version: 2
updates:
  - package-ecosystem: gomod
    directory: "/"
    schedule:
      interval: daily
      time: "20:00"
    open-pull-requests-limit: 10
`
	return path, data
}

func bugReportTemplate(name string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join(".github", "ISSUE_TEMPLATE", "bug_report.md")
	} else {
		path = filepath.Join(name, ".github", "ISSUE_TEMPLATE", "bug_report.md")
	}
	data := `---
name: Bug report
about: Create a report to help us improve
title: "[BUG] XXX"
labels: bug
assignees: ''

---

**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**Screenshots**
If applicable, add screenshots to help explain your problem.

**Desktop (please complete the following information):**
 - OS: [e.g. Ubuntu]
 - Go Version [e.g. 1.17]
 - Application Version [e.g. 1.0.1]

**Additional context**
Add any other context about the problem here.
`
	return path, data
}

func issueTemplate(name string, noRoot bool) (string, string) {
	var path string
	if noRoot {
		path = filepath.Join(".github", "ISSUE_TEMPLATE", "issue.md")
	} else {
		path = filepath.Join(name, ".github", "ISSUE_TEMPLATE", "issue.md")
	}
	data := `---
name: Task
about: Describe this issue
title: ''
labels: ''
assignees: ''

---

## What

Describe what this issue should address.

## How

Describe how to address the issue.

## Checklist

- [ ] Finish implementation of the issue
- [ ] Test all functions
- [ ] Have enough logs to trace activities
- [ ] Notify developers of necessary actions
`
	return path, data
}
