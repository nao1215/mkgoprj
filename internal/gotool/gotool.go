// Package gotool handles information about go commands
package gotool

import (
	"os/exec"
	"regexp"
	"runtime"

	"github.com/nao1215/mkgoprj/internal/ioutils"
)

// Version return runtime golang version(only number, not include "go" or "cpu name")
func Version() string {
	runtimeVer := runtime.Version()

	rex := regexp.MustCompile("[0-9]\\.[0-9]*")
	return rex.FindString(runtimeVer)
}

// ModInit execute "$ go mod init <importPath>"
// If it can not execute "$ go mod", exit command.
func ModInit(importPath string) {
	if err := exec.Command("go", "mod", "init", importPath).Run(); err != nil {
		ioutils.Die(err.Error())
	}
}

// ModTidy execute "$ go mod tidy"
// If it can not execute "$ go mod", exit command.
func ModTidy() {
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		ioutils.Die(err.Error())
	}
}

// CanUseGoCmd check whether go command install in the system.
// If not install, exit command.
func CanUseGoCmd() {
	_, err := exec.LookPath("go")
	if err != nil {
		ioutils.Die("this system does not install go cmd. Please download golang")
	}
}
