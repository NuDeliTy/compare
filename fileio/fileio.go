package fileio

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var CurrentProgram string

func ProgramDir() string {
	return filepath.Join("programs", CurrentProgram)
}

func ProgramPath(filename string) string {
	return filepath.Join(ProgramDir(), filename)
}

func EnsureProgramDir() error {
	return os.MkdirAll(ProgramDir(), 0755)
}

func SelectProgram() string {
	fmt.Print("Enter program name: ")
	var programName string
	fmt.Scanln(&programName)
	CurrentProgram = programName
	EnsureProgramDir()
	return programName
}

func ShowPrograms() {
	programsDir := "programs"
	files, err := os.ReadDir(programsDir)
	if err != nil {
		fmt.Println("No programs found.")
		return
	}

	fmt.Println("\nAvailable Programs:")
	fmt.Println("──────────────────")
	for _, f := range files {
		if f.IsDir() {
			fmt.Printf("• %s\n", f.Name())
		}
	}
	fmt.Println()
}

func AddOutput(fileType string) error {
	if CurrentProgram == "" {
		return fmt.Errorf("no program selected")
	}

	if err := EnsureProgramDir(); err != nil {
		return fmt.Errorf("error creating program directory: %v", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Paste %s output (end with END on a line by itself):\n", fileType)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "END" {
			break
		}
		lines = append(lines, line)
	}
	output := strings.Join(lines, "\n")
	return os.WriteFile(ProgramPath(fileType+".txt"), []byte(output), 0644)
}

func ReadOutput(fileType string) ([]string, error) {
	data, err := os.ReadFile(ProgramPath(fileType + ".txt"))
	if err != nil {
		return nil, err
	}

	content := strings.TrimRight(string(data), "\n")
	if content == "" {
		return []string{}, nil
	}
	return strings.Split(content, "\n"), nil
}

func DeleteProgram() error {
	if CurrentProgram == "" {
		return fmt.Errorf("no program selected")
	}

	err := os.RemoveAll(ProgramDir())
	if err != nil {
		return err
	}

	CurrentProgram = ""
	return nil
}

func DeleteDiff() error {
	if CurrentProgram == "" {
		return fmt.Errorf("no program selected")
	}

	os.Remove(ProgramPath("diff.txt"))
	os.Remove(ProgramPath("diff.html"))
	return nil
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
