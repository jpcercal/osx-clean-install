package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
)

// -------------- Themed Standard Output --------------

// Theme defines the styling theme for the output
type Theme struct {
	PrimaryColor      pterm.Style
	SecondaryColor    pterm.Style
	ErrorColor        pterm.Style
	InfoColor         pterm.Style
	SuccessColor      pterm.Style
	NumberOfColumns   uint8
	TitleLineSymbol   string
	SectionLineSymbol string
}

// ThemedStandardOutput manages themed output to console
type ThemedStandardOutput struct {
	theme Theme
}

// NewThemedStandardOutput creates a new ThemedStandardOutput instance
func NewThemedStandardOutput(theme Theme) *ThemedStandardOutput {
	return &ThemedStandardOutput{theme: theme}
}

// RenderStyledText renders text with a given style
func (ts *ThemedStandardOutput) RenderStyledText(style pterm.Style, text string) {
	style.Println(text)
}

// RenderLine renders a line with a given style and symbol
func (ts *ThemedStandardOutput) RenderLine(style pterm.Style, symbol string) {
	line := strings.Repeat(symbol, int(ts.theme.NumberOfColumns))
	style.Println(line)
}

// PrintTitle prints a titled section
func (ts *ThemedStandardOutput) PrintTitle(title string) {
	padding := (int(ts.theme.NumberOfColumns) / 2) - (len(title) / 2)
	symbol := " "

	message := fmt.Sprintf("%s%s%s", strings.Repeat(symbol, padding), strings.ToUpper(title), strings.Repeat(symbol, padding))
	message = strings.ReplaceAll(message, ";", " ")

	ts.RenderLine(ts.theme.PrimaryColor, ts.theme.TitleLineSymbol)
	ts.RenderStyledText(ts.theme.PrimaryColor, message)
	fmt.Println()
	ts.RenderLine(ts.theme.PrimaryColor, ts.theme.TitleLineSymbol)
}

// PrintSection prints a section with a given title
func (ts *ThemedStandardOutput) PrintSection(section string) {
	padding := (int(ts.theme.NumberOfColumns) / 2) - (len(section) / 2)
	symbol := " "

	message := fmt.Sprintf("%s%s%s", strings.Repeat(symbol, padding), section, strings.Repeat(symbol, padding))
	message = strings.ReplaceAll(message, ";", " ")

	ts.RenderLine(ts.theme.SecondaryColor, ts.theme.SectionLineSymbol)
	ts.RenderStyledText(ts.theme.SecondaryColor, message)
	fmt.Println()
	ts.RenderLine(ts.theme.SecondaryColor, ts.theme.SectionLineSymbol)
}

// PrintParagraph prints a paragraph with a specified style
func (ts *ThemedStandardOutput) PrintParagraph(style pterm.Style, message string) {
	formattedMessage := FoldText(message, ts.theme.NumberOfColumns)
	ts.RenderStyledText(style, formattedMessage)
	fmt.Println()
}

// FoldText folds text to fit within a given width
func FoldText(text string, width uint8) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var buffer bytes.Buffer
	line := words[0]

	for _, word := range words[1:] {
		if len(line)+len(word)+1 > int(width) {
			buffer.WriteString(line)
			buffer.WriteString("\n")
			line = word
		} else {
			line += " " + word
		}
	}
	buffer.WriteString(line)

	return buffer.String()
}

// Info prints an informational message
func (ts *ThemedStandardOutput) Info(message string) {
	ts.PrintParagraph(ts.theme.InfoColor, fmt.Sprintf("[info] %s", message))
}

// Success prints a success message
func (ts *ThemedStandardOutput) Success(message string) {
	ts.PrintParagraph(ts.theme.SuccessColor, fmt.Sprintf("[success] %s", message))
}

// Error prints an error message
func (ts *ThemedStandardOutput) Error(message string) {
	ts.PrintParagraph(ts.theme.ErrorColor, fmt.Sprintf("[error] %s", message))
}

// Command executes a command and prints the output
func (ts *ThemedStandardOutput) Command(command, description string, execute bool) {
	message := FoldText(command, ts.theme.NumberOfColumns)

	ts.RenderLine(ts.theme.PrimaryColor, ts.theme.SectionLineSymbol)
	ts.RenderStyledText(pterm.Style{}, fmt.Sprintf("$ %s\n", message))
	ts.RenderStyledText(pterm.Style{pterm.FgGreen}, fmt.Sprintf("%s\n", description))

	if !execute {
		return
	}

	start := time.Now()
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	end := time.Now()
	elapsed := end.Sub(start).Seconds()

	if err == nil {
		ts.Success(fmt.Sprintf("It took \"%.2fs\" to complete this job.", elapsed))
	} else {
		ts.Error(fmt.Sprintf("It took \"%.2fs\" to complete this job.", elapsed))
		logrus.WithError(err).Error("Command execution failed")
	}
}
