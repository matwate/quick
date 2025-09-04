package renderer

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/pmezard/go-difflib/difflib"
)

func Render(s string) {
	fmt.Println(s)
}

func RenderGreen(s string) {
	color.Green(s)
}

func RenderRed(s string) {
	color.Red(s)
}

func RenderYellow(s string) {
	color.Yellow(s)
}

func RenderBlue(s string) {
	color.Blue(s)
}

func RenderDiff(expected, actual string) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(expected),
		B:        difflib.SplitLines(actual),
		FromFile: "Expected",
		ToFile:   "Actual",
		Context:  3,
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	fmt.Println(text)
}

func RenderTable(header []string, data [][]string) {
	// Find the maximum width for each column
	widths := make([]int, len(header))
	for i, h := range header {
		widths[i] = len(h)
	}
	for _, row := range data {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print header
	head := ""
	for i, h := range header {
		head += fmt.Sprintf("%-*s", widths[i]+2, h)
	}
	Render(head)
	Render(strings.Repeat("-", len(head)))

	// Print data
	for _, row := range data {
		line := ""
		for i, cell := range row {
			line += fmt.Sprintf("%-*s", widths[i]+2, cell)
		}
		Render(line)
	}
}

