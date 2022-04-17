// Package project manage and operate the project information to be generated.
package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/nao1215/mkgoprj/internal/gotool"
	"github.com/nao1215/mkgoprj/internal/ioutils"
	"github.com/nao1215/mkgoprj/internal/target"
)

type Option struct {
}

// Project have project information to be generated.
type Project struct {
	importPath string            // same as "$ git mod init <importPath>"
	name       string            // project (command) name
	library    bool              // it mean library project
	cli        bool              // it means cli project with cobra
	noRoot     bool              // whether create project root directory or not
	files      map[string]string // File to be created: key=file path, value=text in file
	dirs       []string          // directory to be created
}

// NewProject return initialized project struct.
func NewProject(importPath string, lib, cli, noRoot bool) *Project {
	var prj Project
	prj.importPath = importPath
	prj.name = filepath.Base(prj.importPath)
	prj.library = lib
	prj.cli = cli
	prj.noRoot = noRoot
	prj.files = target.Files(prj.name, importPath, lib, cli, noRoot)
	prj.dirs = target.Dirs(prj.name, lib, cli, noRoot)
	return &prj
}

// Make generate project directory and files.
func (p *Project) Make() {
	now := time.Now()

	p.printStartBanner()
	p.canMake()
	p.makeProjectDirs()
	p.makeProjectFiles()
	p.printDirTree()
	p.goModInit()
	if p.cli {
		p.goModTidy()
	}

	ms := time.Since(now).Milliseconds()
	p.printEndBanner(ms)
}

// printStartBanner displays a banner to start creating a project.
func (p *Project) printStartBanner() {
	kind := "application"
	if p.library {
		kind = "library"
	}
	fmt.Printf("%s starts creating the '%s' %s project (import path='%s')\n\n",
		color.HiYellowString("mkgoprj"), color.GreenString(p.name), kind,
		color.GreenString(p.importPath))
}

// printEndBanner displays a banner to end creating a project.
func (p *Project) printEndBanner(ms int64) {
	fmt.Println("")
	fmt.Printf("%s in %d[ms]\n", color.GreenString("BUILD SUCCESSFUL"), ms)
}

func (p *Project) printDirTree() {
	projectDir := p.name
	if p.noRoot {
		projectDir = "."
	}

	fmt.Printf("        %s (your project root)\n", color.YellowString(projectDir))
	ioutils.Tree(projectDir)
}

// canMake check whether can create project template or not.
// If it can't create the project, exit command.
func (p *Project) canMake() {
	fmt.Printf("[%s] check if %s can create the project\n",
		color.GreenString("START"), ioutils.CmdName)

	gotool.CanUseGoCmd()
	if strings.Trim(p.name, " ") == "" {
		ioutils.Die("project name is empty (import path end with \"/ \"?)")
	}
	p.canMakePrjFile()
}

// makeProjectDirs create all directories in project template.
// If it can not make directories, exit command.
func (p *Project) makeProjectDirs() {
	fmt.Printf("[%s] create directories\n", color.GreenString("START"))
	ioutils.MkDirs(p.dirs)
}

// makeProjectDirs create all files in application project template.
func (p *Project) makeProjectFiles() {
	fmt.Printf("[%s] create files\n", color.GreenString("START"))
	for path, code := range p.files {
		ioutils.WriteFile(code, path)
	}
}

// canMakePrjFile check whether the file mkgoprj is trying to generate already exists.
func (p *Project) canMakePrjFile() {
	var files []string
	for k := range p.files {
		files = append(files, k)
	}
	files = append(files, p.dirs...)

	for _, v := range files {
		if ioutils.Exists(v) {
			ioutils.Die("same name file (" + v + ") already exists")
		}
	}
}

// goModInit execute "$ go mod init <importPath>"
// If it can not execute "$ go mod", exit command.
func (p *Project) goModInit() {
	fmt.Printf("[%s] Execute 'go mod init %s'\n", color.GreenString("START"), p.importPath)
	preDir, err := os.Getwd()
	if err != nil {
		ioutils.Die(err.Error())
	}

	if !p.noRoot {
		err := os.Chdir(p.name)
		if err != nil {
			ioutils.Die(err.Error())
		}
	}
	gotool.ModInit(p.importPath)

	err = os.Chdir(preDir)
	if err != nil {
		ioutils.Die(err.Error())
	}
}

// goModTidy execute "$ go mod tidy"
// If it can not execute "$ go mod", exit command.
func (p *Project) goModTidy() {
	fmt.Printf("[%s] Execute 'go mod tidy'\n", color.GreenString("START"))
	preDir, err := os.Getwd()
	if err != nil {
		ioutils.Die(err.Error())
	}

	if !p.noRoot {
		err := os.Chdir(p.name)
		if err != nil {
			ioutils.Die(err.Error())
		}
	}
	gotool.ModTidy()

	err = os.Chdir(preDir)
	if err != nil {
		ioutils.Die(err.Error())
	}
}
