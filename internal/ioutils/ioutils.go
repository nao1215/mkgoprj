// Package ioutils provides APIs related to input and output.
package ioutils

import (
	"fmt"
	"os"
)

// CmdName is this command name.
const CmdName string = "ubume"

// Exists check whether file or directory exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return (err == nil)
}

// WriteFile write string to file.
// If it can not create file, exit command.
func WriteFile(text string, path string) {
	file, err := os.Create(path)
	if err != nil {
		Die(err.Error())
	}
	defer file.Close()
	file.Write(([]byte)(text))
}

// MkDirs create multiple specified directories.
// If the parent directory does not exist, create the parent directory as well.
// If an error occurs, exit command.
func MkDirs(paths []string) {
	for _, path := range paths {
		target := os.ExpandEnv(path)
		err := os.MkdirAll(target, 0755)
		if err != nil {
			Die(err.Error())
		}
	}
}

// Die exit program with message.
func Die(msg string) {
	fmt.Fprintln(os.Stderr, CmdName+": "+msg)
	os.Exit(1)
}
