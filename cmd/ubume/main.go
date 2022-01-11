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
	"runtime"
	"strings"

	"github.com/jessevdk/go-flags"
)

type options struct {
	Version bool `short:"v" long:"version" description:"Show ubume command version"`
}

var osExit = os.Exit

const cmdName string = "ubume"
const version string = "0.5.0"

const (
	exitSuccess int = iota // 0
	exitFailure
)

type project struct {
	importPath string
	name       string
	version    string
}

func main() {
	var opts options
	var args []string
	var err error

	if args, err = parseArgs(&opts); err != nil {
		osExit(exitFailure)
	}

	prj := initProject(args)
	prj.canMake()
	prj.make()
}

func (p project) make() {
	p.makeProjectDirs()
	p.makeFiles()
	p.goModInit()
}

func (p project) makeProjectDirs() {
	dirs := []string{
		filepath.Join(p.name, "cmd", p.name),
	}

	for _, path := range dirs {
		target := os.ExpandEnv(path)
		err := os.MkdirAll(target, 0755)
		if err != nil {
			die(err.Error())
		}
	}
}

func (p project) makeFiles() {
	p.makeMainSourecCodeFile()
	p.makeMainSourecCodeTestFile()
	p.makeMakefile()
	p.makeChangelogFile()
}

func (p project) makeMainSourecCodeFile() {
	path := filepath.Join(p.name, "cmd", p.name, "main.go")
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

func (p project) makeMainSourecCodeTestFile() {
	path := filepath.Join(p.name, "cmd", p.name, "main_test.go")
	code := `package main

import "testing"

func TestHelloWorld(t *testing.T) {
	if HelloWorld() != "Hello, World" {
		t.Errorf("HelloWorlf = %s, want \"Hello, World\"", HelloWorld())
	}
}
`
	writeFile(code, path)
}

func (p project) makeMakefile() {
	path := filepath.Join(p.name, "Makefile")
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
GO_PACKAGES = $(shell $(GO_LIST) $(GO_PKGROOT) | grep -v vendor)

build: deps ## Build binary 
	env GO111MODULE=on GOOS=$(GOOS) $(GO_BUILD) $(GO_LDFLAGS) -o $(APP) XXX_CODE_XXX

clean: ## Clean project
	-rm -rf ./vendor $(APP) cover.out cover.html

deps: ## Dependency resolution for build
	$(GO_DEP) vendor

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

func (p project) makeChangelogFile() {
	path := filepath.Join(p.name, "Changelog.md")
	data := `# Changelog
All notable changes to this project will be documented in this file.  
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).   
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
`
	writeFile(data, path)
}

func (p project) goModInit() {
	err := os.Chdir(p.name)
	if err != nil {
		die(err.Error())
	}

	err = exec.Command("go", "mod", "init", p.importPath).Run()
	if err != nil {
		die(err.Error())
	}
}

func (p project) canMake() {
	p.canUseGoCmd()
	p.canMakePrjDir()
}

func (p project) canUseGoCmd() {
	_, err := exec.LookPath("go")
	if err != nil {
		die("this system does not install go cmd. Please download golang")
	}
}

func (p project) canMakePrjDir() {
	if exists(p.name) {
		die("same name project already exists at current directory")
	}
	if p.name == "" {
		die("project name is empty (import path end with \"/\"?)")
	}
}

func die(msg string) {
	fmt.Fprintln(os.Stderr, cmdName+": "+msg)
	osExit(exitFailure)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return (err == nil)
}

func writeFile(text string, path string) {
	file, err := os.Create(path)
	if err != nil {
		die(err.Error())
	}
	defer file.Close()
	file.Write(([]byte)(text))
}

func initProject(arg []string) project {
	var prj project
	prj.importPath = arg[0]

	if strings.Contains(prj.importPath, "/") {
		arr := strings.Split(prj.importPath, "/")
		prj.name = arr[len(arr)-1]
	} else {
		prj.name = prj.importPath
	}
	prj.version = "0.0.1"
	return prj
}

func parseArgs(opts *options) ([]string, error) {
	p := initParser(opts)

	args, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if opts.Version {
		showVersion(cmdName, version)
		osExit(exitSuccess)
	}

	if len(args) != 1 {
		showHelp(p)
		osExit(exitFailure)
	}
	return args, nil
}

func showHelp(p *flags.Parser) {
	p.WriteHelp(os.Stdout)
}

func initParser(opts *options) *flags.Parser {
	parser := flags.NewParser(opts, flags.Default)
	parser.Name = cmdName
	parser.Usage = "[OPTIONS] IMPORT_PATH    ※ IMPORT_PATH is same as $ go mod init IMPORT_PATH"

	return parser
}
func showVersion(cmdName string, version string) {
	description := cmdName + " version " + version + " (under Apache License verison 2.0)"
	fmt.Fprintln(os.Stdout, description)
}
