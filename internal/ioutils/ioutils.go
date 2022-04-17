// Package ioutils provides APIs related to input and output.
package ioutils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var (
	// Stdout is new instance of Writer which handles escape sequence for stdout.
	Stdout = colorable.NewColorableStdout()
	// Stderr is new instance of Writer which handles escape sequence for stderr.
	Stderr = colorable.NewColorableStderr()
)

// CmdName is this command name.
const CmdName string = "mkgoprj"

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

// Warn print warning message
func Warn(msg string) {
	fmt.Fprintf(Stderr, "[%s] %s: %s\n", color.YellowString("WARN "), CmdName, msg)
	os.Exit(1)
}

// Die exit program with message.
func Die(msg string) {
	fmt.Fprintf(Stderr, "[%s] %s: %s\n", color.RedString("ERROR"), CmdName, msg)
	os.Exit(1)
}

// Tree display directories in the tree structure.
func Tree(path string) {
	if err := tree("        ", path); err != nil {
		Die("can not print directory-tree: " + err.Error())
	}
}

func tree(indent, path string) error {
	dirs, files, err := readDirs(path)
	if err != nil {
		return err
	}

	for i, v := range files {
		s := indent + " ├─"
		if len(dirs) == 0 && i == len(files)-1 {
			s = indent + " └─"
		}

		fmt.Printf("%s %s\n", s, v)
	}

	for i, v := range dirs {
		s := indent + " ├─"
		a := " │ "
		if i == len(dirs)-1 {
			s = indent + " └─"
			a = "   "
		}
		fmt.Printf("%s %s\n", s, v)

		if err := tree(indent+a, filepath.Join(path, v)); err != nil {
			return err
		}
	}
	return nil
}

func readDirs(name string) ([]string, []string, error) {
	fp, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}

	list, err := fp.Readdir(-1)
	fp.Close()
	if err != nil {
		return nil, nil, err
	}

	dirs, files := []string{}, []string{}
	for _, v := range list {
		if v.IsDir() {
			dirs = append(dirs, v.Name())
		} else {
			files = append(files, v.Name())
		}
	}
	return dirs, files, nil
}

// IsFile reports whether the path exists and is a file.
func IsFile(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil) && (!stat.IsDir())
}
