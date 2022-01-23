// Package project manage and operate the project information to be generated.
package project

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/nao1215/ubume/internal/gotool"
	"github.com/nao1215/ubume/internal/ioutils"
	"github.com/nao1215/ubume/internal/target"
)

// Project have project information to be generated.
type Project struct {
	importPath string            // same as "$ git mod init <importPath>"
	name       string            // project (command) name
	version    string            // project version
	library    bool              // it's mean library project
	noRoot     bool              // whether create project root directory or not
	files      map[string]string // File to be created: key=file path, value=text in file
	dirs       []string          // directory to be created
}

// NewProject return initialized project struct.
func NewProject(importPath string, lib, noRoot bool) *Project {
	var prj Project
	prj.importPath = importPath
	prj.name = filepath.Base(prj.importPath)
	prj.version = "0.0.1"
	prj.library = lib
	prj.noRoot = noRoot
	prj.files = target.Files(prj.name, lib, noRoot)
	prj.dirs = target.Dirs(prj.name, lib, noRoot)
	return &prj
}

// Make generate project directory and files.
func (p *Project) Make() {
	p.canMake()
	p.makeProjectDirs()
	p.makeProjectFiles()
	p.goModInit()
}

// canMake check whether can create project template or not.
// If it can't create the project, exit command.
func (p *Project) canMake() {
	gotool.CanUseGoCmd()
	if strings.Trim(p.name, " ") == "" {
		ioutils.Die("project name is empty (import path end with \"/ \"?)")
	}
	p.canMakePrjFile()
}

// makeProjectDirs create all directories in project template.
// If it can not make directories, exit command.
func (p *Project) makeProjectDirs() {
	ioutils.MkDirs(p.dirs)
}

// makeProjectDirs create all files in application project template.
func (p *Project) makeProjectFiles() {
	for path, code := range p.files {
		ioutils.WriteFile(code, path)
	}
}

// canMakePrjFile check whether the file ubume is trying to generate already exists.
func (p *Project) canMakePrjFile() {
	var files []string
	for k := range p.files {
		files = append(files, k)
	}
	files = append(files, p.dirs...)

	for _, v := range files {
		if ioutils.Exists(v) {
			ioutils.Die("same name file (" + v + ") already exists at current directory")
		}
	}
}

// goModInit execute "$ go mod init <importPath>"
// If it can not execute "$ go mod", exit command.
func (p *Project) goModInit() {
	if !p.noRoot {
		err := os.Chdir(p.name)
		if err != nil {
			ioutils.Die(err.Error())
		}
	}
	gotool.ModInit(p.importPath)
}
