package diff

import (
	"html/template"
	"strings"

	"koodcompare/fileio"
)

func GenerateHTMLDiff(expected, actual []string) string {
	const htmlTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>Diff Report - {{.Program}}</title>
    <style>
        body { font-family: monospace; margin: 20px; }
        .diff-table { border-collapse: collapse; width: 100%; }
        .diff-table td, .diff-table th { border: 1px solid #ddd; padding: 8px; }
        .expected { background-color: #ffcccc; }
        .actual { background-color: #ccffcc; }
        .match { background-color: #f0f0f0; }
        .header { background-color: #e0e0e0; font-weight: bold; }
        .extra { background-color: #ccffcc; color: #006600; }
        .missing { background-color: #ffcccc; color: #990000; }
        .difference { background-color: #ffffcc; color: #996600; }
    </style>
</head>
<body>
    <h1>Diff Report - {{.Program}}</h1>
    <table class="diff-table">
        <tr class="header">
            <th width="50%">Expected</th>
            <th width="50%">Actual</th>
        </tr>
        {{range .Lines}}
        <tr>
            <td class="{{.ExpectedClass}}">{{.Expected}}</td>
            <td class="{{.ActualClass}}">{{.Actual}}</td>
        </tr>
        {{if .Note}}
        <tr>
            <td colspan="2" class="{{.NoteClass}}">{{.Note}}</td>
        </tr>
        {{end}}
        {{end}}
    </table>
</body>
</html>`

	type DiffLine struct {
		Expected      string
		Actual        string
		ExpectedClass string
		ActualClass   string
		Note          string
		NoteClass     string
	}

	type TemplateData struct {
		Program string
		Lines   []DiffLine
	}

	var lines []DiffLine
	maxLen := len(expected)
	if len(actual) > maxLen {
		maxLen = len(actual)
	}

	for i := 0; i < maxLen; i++ {
		var expLine, actLine string
		if i < len(expected) {
			expLine = expected[i]
		}
		if i < len(actual) {
			actLine = actual[i]
		}

		diffLine := DiffLine{
			Expected: expLine,
			Actual:   actLine,
		}

		if expLine == actLine {
			diffLine.ExpectedClass = "match"
			diffLine.ActualClass = "match"
		} else {
			if expLine == "" {
				diffLine.ExpectedClass = "match"
				diffLine.ActualClass = "extra"
				diffLine.Note = "EXTRA LINE IN ACTUAL OUTPUT"
				diffLine.NoteClass = "extra"
			} else if actLine == "" {
				diffLine.ExpectedClass = "missing"
				diffLine.ActualClass = "match"
				diffLine.Note = "MISSING LINE IN ACTUAL OUTPUT"
				diffLine.NoteClass = "missing"
			} else {
				diffLine.ExpectedClass = "expected"
				diffLine.ActualClass = "actual"
				diffLine.Note = "CONTENT DIFFERENCE"
				diffLine.NoteClass = "difference"
			}
		}

		lines = append(lines, diffLine)
	}

	tmpl, _ := template.New("diff").Parse(htmlTemplate)
	var result strings.Builder
	tmpl.Execute(&result, TemplateData{
		Program: fileio.CurrentProgram,
		Lines:   lines,
	})

	return result.String()
}
