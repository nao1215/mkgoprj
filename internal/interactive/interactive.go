package interactive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/nao1215/ubume/internal/ioutils"
)

// ImportPath return project import path from user input.
func ImportPath() (string, error) {
	response, err := question("[" + color.GreenString("Import PATH") + "]: ")
	if err != nil {
		return "", fmt.Errorf("can not get import path: %w", err)
	}
	return response, nil
}

// ProjectKind return project kind from user input.
func ProjectKind() (string, error) {
	intro := "[" + color.GreenString("Project Kind") + "]\n"
	kinds := []string{
		"  1: Simple Application\n",
		"  2: Library\n",
		"  3: CLI tool with cobra\n",
	}
	q := "  Select value with 1-" + fmt.Sprint(len(kinds)) + ": "

	response, err := question(intro + strings.Join(kinds, "") + q)
	if err != nil {
		return "", fmt.Errorf("can not get project kind: %w", err)
	}

	switch response {
	case "1", "2", "3":
		return response, nil
	default:
		return ProjectKind()
	}
}

// ProjectRootDir returns where to create the project(project path)
func ProjectRootDir(projectName string) (string, error) {
	intro := "[" + color.GreenString("Project Root PATH") + "]\n"
	root, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("can not get current directory path: %w", err)
	}
	prjRoot := "  1: " + filepath.Join(root, projectName) + "\n"
	noRoot := "  2: " + root + "\n"
	q := "  Select value with 1-2: "

	response, err := question(intro + prjRoot + noRoot + q)
	if err != nil {
		return "", fmt.Errorf("can not get import path: %w", err)
	}
	switch response {
	case "1", "2":
		return response, nil
	default:
		return ProjectRootDir(projectName)
	}
}

// IsNoRoot return whether create project root directory
func IsNoRoot(no string) bool {
	return no == "2"
}

// IsSimpleApp return whether the number specifies the creation of the simple project
func IsSimpleApp(no string) bool {
	return no == "1"
}

// IsLibrary return whether the number specifies the creation of the library project
func IsLibrary(no string) bool {
	return no == "2"
}

// IsCLI return whether the number specifies the creation of the cli project
func IsCLI(no string) bool {
	return no == "3"
}

// FinalConfirm check if you can create a project with the above contents
func FinalConfirm(importPath, projectName, projectKind, noRoot string) bool {
	name := "[" + color.GreenString("Project Name") + "] " + projectName + "\n"
	kind := "[" + color.GreenString("Project Kind") + "] " + projectKindStr(projectKind) + "\n"
	root := "[" + color.GreenString("Project Root") + "] " + rootDirPATHStr(noRoot, projectName) + "\n"
	ip := "[" + color.GreenString("Import  PATH") + "] " + importPath + "\n"
	q := " Create the project with the above contents[" + color.YellowString("Y/N") + "]: "

	fmt.Println("")
	fmt.Println(color.YellowString("===Confirm the input information==="))
	response, err := question(name + kind + root + ip + q)
	if err != nil {
		ioutils.Die("can not get 'Y' or 'N' from user input: " + err.Error())
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return FinalConfirm(importPath, projectName, projectKind, noRoot)
	}
}

func projectKindStr(no string) string {
	switch no {
	case "1":
		return "Simple Application"
	case "2":
		return "Library"
	case "3":
		return "CLI tool with cobra"
	default:
		return ""
	}
}

func rootDirPATHStr(no, projectName string) string {
	root, err := os.Getwd()
	if err != nil {
		ioutils.Die("can not get current directory path: " + err.Error())
	}
	prjRoot := filepath.Join(root, projectName)
	noRoot := root
	switch no {
	case "1":
		return prjRoot
	case "2":
		return noRoot
	default:
		return ""
	}
}

func question(ask string) (string, error) {
	var response string
	fmt.Fprintf(os.Stdout, ask)
	_, err := fmt.Scanln(&response)
	if err != nil {
		if strings.Contains(err.Error(), "expected newline") {
			return question(ask) // If user input only enter.
		}
		return "", err
	}
	return response, nil
}
