package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"

	"koodcompare/diff"
	"koodcompare/fileio"
)

func PrintColoredDiff(diffText string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	lines := strings.Split(diffText, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Handle header lines
		if strings.HasPrefix(trimmed, "DIFF REPORT") {
			fmt.Println(cyan(line))
			continue
		}
		if strings.HasPrefix(trimmed, "===========") {
			fmt.Println(cyan(line))
			continue
		}
		if strings.HasPrefix(trimmed, "Program:") {
			fmt.Println(blue(line))
			continue
		}

		// Handle column headers and separators
		if strings.Contains(line, "Expected") && strings.Contains(line, "Actual") {
			fmt.Println(blue(line))
			continue
		}
		if strings.Contains(line, "--------------------------------") {
			fmt.Println(blue(line))
			continue
		}

		// Handle difference indicators
		if strings.Contains(line, "EXTRA LINE") {
			fmt.Println(green(line))
			continue
		}
		if strings.Contains(line, "MISSING LINE") {
			fmt.Println(red(line))
			continue
		}
		if strings.Contains(line, "CONTENT DIFFERENCE") {
			fmt.Println(yellow(line))
			continue
		}
		if strings.Contains(line, "^") {
			fmt.Println(yellow(line))
			continue
		}

		// Handle the main comparison lines
		if strings.Contains(line, "|") {
			parts := strings.SplitN(line, "|", 2)
			if len(parts) != 2 {
				fmt.Println(line)
				continue
			}

			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			// Check if this is a content difference line
			if left != right && left != "" && right != "" {
				// PRESERVE THE ORIGINAL FORMATTING while adding colors
				leftFormatted := parts[0]  // Keep the original formatting with spaces
				rightFormatted := parts[1] // Keep the original formatting with spaces

				// Colorize the content but preserve the spacing
				coloredLeft := strings.Replace(leftFormatted, left, red(left), 1)
				coloredRight := strings.Replace(rightFormatted, right, green(right), 1)

				fmt.Printf("%s|%s\n", coloredLeft, coloredRight)
			} else {
				// Regular line or other cases
				fmt.Println(line)
			}
			continue
		}

		// Default case
		fmt.Println(line)
	}
}

func ShowDiffPreview(diffText string) {
	fmt.Println("\n" + strings.Repeat("═", 50))
	color.New(color.FgHiCyan).Println("COLORED DIFF PREVIEW")
	fmt.Println(strings.Repeat("═", 50))
	PrintColoredDiff(diffText)
	fmt.Println(strings.Repeat("═", 50))
}

func Pause() {
	color.New(color.FgHiYellow).Print("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func SelectProgram() {
	color.New(color.FgHiYellow).Print("Enter program name: ")
	fmt.Scanln(&fileio.CurrentProgram)
	color.New(color.FgHiGreen).Printf("Selected program: %s\n", fileio.CurrentProgram)
	fileio.EnsureProgramDir()
}

func ShowPrograms() {
	color.New(color.FgHiCyan).Println("\nAvailable Programs:")
	color.New(color.FgHiCyan).Println("──────────────────")
	fileio.ShowPrograms()
}

func AddOutput(fileType string) {
	if fileio.CurrentProgram == "" {
		color.New(color.FgHiRed).Println("No program selected!")
		return
	}

	err := fileio.AddOutput(fileType)
	if err != nil {
		color.New(color.FgHiRed).Printf("Error: %v\n", err)
		return
	}

	color.New(color.FgHiGreen).Printf("%s output saved for %s\n", fileType, fileio.CurrentProgram)
}

func CompareOutputs() {
	if fileio.CurrentProgram == "" {
		color.New(color.FgHiRed).Println("No program selected!")
		return
	}

	expected, err1 := fileio.ReadOutput("expected")
	actual, err2 := fileio.ReadOutput("actual")

	if err1 != nil || err2 != nil {
		color.New(color.FgHiRed).Println("Error reading output files. Make sure both expected and actual outputs exist.")
		return
	}

	textDiff, htmlDiff, hasDifferences := diff.Compare(expected, actual)

	// Save diffs to files
	os.WriteFile(fileio.ProgramPath("diff.txt"), []byte(textDiff), 0644)
	os.WriteFile(fileio.ProgramPath("diff.html"), []byte(htmlDiff), 0644)

	if !hasDifferences {
		color.New(color.FgHiGreen).Println("No differences found!")
		return
	}

	color.New(color.FgHiGreen).Println("Comparison saved:")
	color.New(color.FgHiGreen).Printf("  Text: %s\n", fileio.ProgramPath("diff.txt"))
	color.New(color.FgHiGreen).Printf("  HTML: %s\n", fileio.ProgramPath("diff.html"))

	// Show colored preview
	ShowDiffPreview(textDiff)
	Pause()
}

func ShowLastDiff(format string) {
	if fileio.CurrentProgram == "" {
		color.New(color.FgHiRed).Println("No program selected!")
		return
	}

	var filename string
	switch format {
	case "text":
		filename = "diff.txt"
	case "html":
		filename = "diff.html"
	default:
		color.New(color.FgHiRed).Println("Invalid format")
		return
	}

	if !fileio.FileExists(fileio.ProgramPath(filename)) {
		color.New(color.FgHiRed).Printf("No %s diff found for %s\n", format, fileio.CurrentProgram)
		return
	}

	if format == "text" {
		data, _ := os.ReadFile(fileio.ProgramPath(filename))
		ShowDiffPreview(string(data))
		Pause()
	} else {
		color.New(color.FgHiGreen).Printf("HTML diff saved to: %s\n", fileio.ProgramPath("diff.html"))
		color.New(color.FgHiYellow).Println("Open this file in a web browser to view colored diff")
		Pause()
	}
}

func DeleteProgram() {
	if fileio.CurrentProgram == "" {
		color.New(color.FgHiRed).Println("No program selected!")
		return
	}

	err := fileio.DeleteProgram()
	if err != nil {
		color.New(color.FgHiRed).Printf("Error deleting program: %v\n", err)
	} else {
		color.New(color.FgHiGreen).Printf("Program deleted: %s\n", fileio.CurrentProgram)
	}
}

func DeleteDiff() {
	if fileio.CurrentProgram == "" {
		color.New(color.FgHiRed).Println("No program selected!")
		return
	}

	err := fileio.DeleteDiff()
	if err != nil {
		color.New(color.FgHiRed).Printf("Error: %v\n", err)
	} else {
		color.New(color.FgHiGreen).Printf("Deleted diff files for %s\n", fileio.CurrentProgram)
	}
}
