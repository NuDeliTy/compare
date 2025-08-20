package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"

	"koodcompare/ui"
)

var currentProgram string

func main() {
	for {
		clearScreen()
		showMenu()

		choice := getMenuChoice()
		if choice == -1 {
			continue
		}

		handleMenuChoice(choice)
	}
}

func showMenu() {
	titleColor := color.New(color.FgHiCyan, color.Bold)
	optionColor := color.New(color.FgHiWhite)

	titleColor.Println("╔══════════════════════════════╗")
	titleColor.Println("║        TEST MANAGER          ║")
	titleColor.Println("╚══════════════════════════════╝")

	if currentProgram != "" {
		optionColor.Printf("Current Program: ")
		color.New(color.FgHiGreen).Printf("%s\n\n", currentProgram)
	} else {
		color.New(color.FgHiRed).Println("No program selected!\n")
	}

	optionColor.Println("1.  Create/select program")
	optionColor.Println("2.  Show existing programs")
	optionColor.Println("3.  Add expected output")
	optionColor.Println("4.  Add actual output")
	optionColor.Println("5.  Compare expected vs actual")
	optionColor.Println("6.  Show last comparison (text)")
	optionColor.Println("7.  Show last comparison (HTML)")
	optionColor.Println("8.  Delete a program")
	optionColor.Println("9.  Delete comparison for current program")
	optionColor.Println("10. Exit")

	color.New(color.FgHiYellow).Print("\nChoose: ")
}

func getMenuChoice() int {
	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil {
		// Clear input buffer on error
		var discard string
		fmt.Scanln(&discard)
		color.New(color.FgHiRed).Println("Invalid input! Please enter a number.")
		fmt.Print("Press Enter to continue...")
		fmt.Scanln()
		return -1
	}
	return choice
}

func handleMenuChoice(choice int) {
	switch choice {
	case 1:
		ui.SelectProgram()
	case 2:
		ui.ShowPrograms()
	case 3:
		ui.AddOutput("expected")
	case 4:
		ui.AddOutput("actual")
	case 5:
		ui.CompareOutputs()
	case 6:
		ui.ShowLastDiff("text")
	case 7:
		ui.ShowLastDiff("html")
	case 8:
		ui.DeleteProgram()
	case 9:
		ui.DeleteDiff()
	case 10:
		color.New(color.FgHiGreen).Println("Goodbye! 👋")
		os.Exit(0)
	default:
		color.New(color.FgHiRed).Println("Invalid choice! Please select 1-10.")
		fmt.Print("Press Enter to continue...")
		fmt.Scanln()
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J") // ANSI escape codes to clear screen
}

func pause() {
	color.New(color.FgHiYellow).Print("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
