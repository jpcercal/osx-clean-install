package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/jpcercal/dotfiles/internal"
	"github.com/k0kubun/go-ansi"
	"github.com/pterm/pterm"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

// setupLogger configures the log output and format
func setupLogger(ctx context.Context, logFilename string) (*os.File, error) {
	logFile, err := os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.WithContext(ctx).WithField("filename", logFilename).WithError(err).Fatal("failed to open log file")
		return nil, err
	}
	log.SetOutput(logFile)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	return logFile, nil
}

// handleSignals creates a context that cancels when an interrupt signal is received
func handleSignals(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
			return
		}
	}()
	return ctx, cancel
}

// runCommands executes the commands and handles results
func runCommands(ctx context.Context, commands []internal.Command, resultChan chan<- internal.CommandResult, semaphore chan struct{}, wg *sync.WaitGroup) {
	for _, command := range commands {
		wg.Add(1)
		semaphore <- struct{}{} // acquire a token
		go func(ctx context.Context, command internal.Command) {
			defer wg.Done()
			defer func() { <-semaphore }()
			command.Run(ctx, resultChan)
		}(ctx, command)
	}
}

// handleResults processes command results and updates progress
func handleResults(ctx context.Context, resultChan <-chan internal.CommandResult, pc *PostponedCommands) {
	for {
		select {
		case <-ctx.Done():
			// do not return it
			// return
		case result, ok := <-resultChan:
			if !ok {
				return
			}

			func() {
				defer progressBar.UpdateTitle(fmt.Sprintf("%-40.40s", string(result.Command)))
				defer progressBar.Increment()
			}()

			if result.Success {
				pterm.Success.Println("completed successfully.")
				log.WithContext(ctx).WithField("result", result).Info("installation completed successfully.")
			} else {
				handleFailedCommand(ctx, result, pc)
			}
		}
	}
}

// handleFailedCommand processes failed command results
func handleFailedCommand(ctx context.Context, result internal.CommandResult, pc *PostponedCommands) {
	if strings.Contains(strings.ToLower(result.Message), "operation already in progress") {
		pc.mu.Lock()
		pc.pendingCommands = append(pc.pendingCommands, result.Command)
		progressBar.Stop()
		progressBar, _ = progressBar.WithTotal(progressBar.Total + 1).WithCurrent(progressBar.Current).Start()
		pc.mu.Unlock()
		pterm.Info.Println("dependency installation was already in progress, postponing it to be installed later.")
		log.WithContext(ctx).WithField("result", result).Warn("dependency installation was already in progress, postponing it to be installed later.")
		return
	}

	if result.OperationCancelled {
		pterm.Info.Println(result.Message)
		log.WithContext(ctx).WithField("result", result).Info(result.Message)
		return
	}

	pterm.Error.Println(result.Message)
	log.WithContext(ctx).WithField("result", result).Error("installation failed.")
}

type PostponedCommands struct {
	mu              sync.Mutex
	pendingCommands []internal.Command
}

var progressBar *pterm.ProgressbarPrinter

func main() {
	// Define colors based on the theme
	primaryColor := pterm.NewStyle(pterm.FgYellow)
	secondaryColor := pterm.NewStyle(pterm.FgCyan)
	errorColor := pterm.NewStyle(pterm.FgRed)
	infoColor := pterm.NewStyle(pterm.FgYellow)
	successColor := pterm.NewStyle(pterm.FgGreen)

	theme := internal.Theme{
		PrimaryColor:      *primaryColor,
		SecondaryColor:    *secondaryColor,
		ErrorColor:        *errorColor,
		InfoColor:         *infoColor,
		SuccessColor:      *successColor,
		NumberOfColumns:   80,
		TitleLineSymbol:   "=",
		SectionLineSymbol: "-",
	}

	// Start a new fullscreen centered area.
	// This area will be used to display the bar chart.
	area, _ := pterm.DefaultArea..WithFullscreen().WithCenter().Start()

	// Ensure the area stops updating when we're done.
	defer area.Stop()

	// bar, _ := pterm.DefaultProgressbar.WithBarStyle().WithTotal(5).WithTitle("Downloading stuff").Start()

	bar := progressbar.NewOptions(5,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()), //you should install "github.com/k0kubun/go-ansi"
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan][1/3][reset] Writing moshable file..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

        bar.
	tsOutput := internal.NewThemedStandardOutput(theme)
	bar.Add(1)

	tsOutput.PrintTitle("INSTALL DEPENDENCIES")
	tsOutput.PrintSection("Installing Homebrew")
	bar.Add(1)

	tsOutput.Info("Info message example")
	tsOutput.Success("Success message example")
	tsOutput.Error("Error message example")
	bar.Add(1)

	tsOutput.Command("echo Hello, world!", "This is a description of the command.", true)
	bar.Add(1)

	tsOutput.RenderLine(theme.PrimaryColor, theme.SectionLineSymbol)
	bar.Add(1)

	os.Exit(1)

	ctx, cancel := handleSignals(context.Background())
	defer cancel()

	logFilename := "commands.log"
	logFile, err := setupLogger(ctx, logFilename)
	if err != nil {
		os.Exit(1)
	}
	defer logFile.Close()

	commands, err := internal.ResolveCommandsFromFile(ctx, "apps.yaml")
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to resolve commands")
		os.Exit(1)
	}

	progressBar, err = internal.NewProgressBar(len(commands))
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to create progress bar")
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	pc := &PostponedCommands{}

	resultChan := make(chan internal.CommandResult)
	go handleResults(ctx, resultChan, pc)

	semaphore := make(chan struct{}, internal.ResolveNumberOfParallelCommandExecution(ctx))
	runCommands(ctx, commands, resultChan, semaphore, &wg)
	wg.Wait()
	close(semaphore)

	// Handle postponed commands
	pc.mu.Lock()
	pendingCommands := make([]internal.Command, len(pc.pendingCommands))
	copy(pendingCommands, pc.pendingCommands)
	pc.pendingCommands = nil // clear pending commands
	pc.mu.Unlock()

	if len(pendingCommands) > 0 {
		semaphore2 := make(chan struct{}, 1)
		runCommands(ctx, commands, resultChan, semaphore2, &wg)
		wg.Wait()
		close(semaphore2)
	}

	close(resultChan)
	progressBar.Stop()
}
