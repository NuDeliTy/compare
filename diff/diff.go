package diff

import (
	"fmt"
	"strings"

	"koodcompare/fileio"
)

func GenerateTextDiff(expected, actual []string) string {
	var result strings.Builder

	result.WriteString("DIFF REPORT\n")
	result.WriteString("===========\n")
	result.WriteString(fmt.Sprintf("Program: %s\n\n", fileio.CurrentProgram))

	result.WriteString("Expected                         | Actual\n")
	result.WriteString("--------------------------------|--------------------------------\n")

	maxLen := len(expected)
	if len(actual) > maxLen {
		maxLen = len(actual)
	}

	differencesFound := false

	for i := 0; i < maxLen; i++ {
		var expLine, actLine string
		if i < len(expected) {
			expLine = expected[i]
		}
		if i < len(actual) {
			actLine = actual[i]
		}

		expFormatted := fmt.Sprintf("%-30s", expLine)
		actFormatted := fmt.Sprintf("%-30s", actLine)

		// Check if we're at the position where the extra line appears
		if i >= len(expected) && i < len(actual) {
			// Extra line in actual
			differencesFound = true
			result.WriteString(fmt.Sprintf("%s| %s\n", strings.Repeat(" ", 30), actFormatted))
			result.WriteString(fmt.Sprintf("%s| ^^^ EXTRA LINE: '%s'\n", strings.Repeat(" ", 30), actLine))
			result.WriteString("--------------------------------|--------------------------------\n")
			continue
		}

		if i < len(expected) && i >= len(actual) {
			// Missing line in actual
			differencesFound = true
			result.WriteString(fmt.Sprintf("%s| %s\n", expFormatted, strings.Repeat(" ", 30)))
			result.WriteString(fmt.Sprintf("%s| vvv MISSING LINE: '%s'\n", strings.Repeat(" ", 30), expLine))
			result.WriteString("--------------------------------|--------------------------------\n")
			continue
		}

		if expLine == actLine {
			result.WriteString(fmt.Sprintf("%s| %s\n", expFormatted, actFormatted))
		} else {
			differencesFound = true
			result.WriteString(fmt.Sprintf("%s| %s\n", expFormatted, actFormatted))

			if expLine == "" {
				result.WriteString(fmt.Sprintf("%s| ^^^ EXTRA LINE: '%s'\n", strings.Repeat(" ", 30), actLine))
			} else if actLine == "" {
				result.WriteString(fmt.Sprintf("%s| vvv MISSING LINE: '%s'\n", strings.Repeat(" ", 30), expLine))
			} else {
				result.WriteString(fmt.Sprintf("%s| >>> CONTENT DIFFERENCE <<<\n", strings.Repeat(" ", 30)))

				// Show character-by-character difference
				minLen := len(expLine)
				if len(actLine) < minLen {
					minLen = len(actLine)
				}

				differPos := -1
				for j := 0; j < minLen; j++ {
					if expLine[j] != actLine[j] {
						differPos = j
						break
					}
				}

				if differPos != -1 {
					indicator := strings.Repeat(" ", 30) + "| "
					indicator += strings.Repeat(" ", differPos) + "^"
					result.WriteString(fmt.Sprintf("%s\n", indicator))
				}
			}
			result.WriteString("--------------------------------|--------------------------------\n")
		}
	}

	if !differencesFound {
		return "NO DIFFERENCES FOUND - OUTPUTS MATCH EXACTLY\n"
	}

	return result.String()
}

func Compare(expected, actual []string) (string, string, bool) {
	textDiff := GenerateTextDiff(expected, actual)
	htmlDiff := GenerateHTMLDiff(expected, actual)

	hasDifferences := false
	lines := strings.Split(textDiff, "\n")
	for _, line := range lines {
		if strings.Contains(line, "EXTRA LINE") || strings.Contains(line, "MISSING LINE") || strings.Contains(line, "CONTENT DIFFERENCE") {
			hasDifferences = true
			break
		}
	}

	return textDiff, htmlDiff, hasDifferences
}
