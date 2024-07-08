package internal

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const MAX_ALLOWED_PARALLEL_INSTALLATION_PROCESSES = 4

type Command string

type CommandResult struct {
	Command            Command `yaml:"command" json:"command"`
	Executed           bool    `yaml:"executed" json:"executed"`
	ExecutionTimeTs    int64   `yaml:"execution_time" json:"execution_time"`
	Success            bool    `yaml:"success" json:"success"`
	OperationCancelled bool    `yaml:"operation_cancelled" json:"operation_cancelled"`
	Message            string  `yaml:"message" json:"message"`
}

type Brew struct {
	Taps           []string `yaml:"taps" json:"taps"`
	Formulas       []string `yaml:"formulas" json:"formulas"`
	Casks          []string `yaml:"casks" json:"casks"`
	CustomCommands []string `yaml:"customCommands" json:"customCommands"`
}

type Gem struct {
	Rubygems []string `yaml:"rubygems" json:"rubygems"`
}

type Npm struct {
	Global struct {
		Packages interface{} `yaml:"packages" json:"packages"`
	} `yaml:"global" json:"global"`
}

type Mas struct {
	Apps interface{} `yaml:"apps" json:"apps"`
}

type SymbolicLink struct {
	From struct {
		RelativePath string `yaml:"relative_path" json:"relative_path"`
	} `yaml:"from" json:"from"`
	To struct {
		AbsolutePath string `yaml:"absolute_path" json:"absolute_path"`
	} `yaml:"to" json:"to"`
}

type Dockutil struct {
	Before struct {
		Reset     bool `yaml:"reset" json:"reset"`
		RemoveAll bool `yaml:"removeAll" json:"removeAll"`
	} `yaml:"_before" json:"_before"`
	Add []struct {
		App   string `yaml:"app" json:"app"`
		After string `yaml:"after" json:"after"`
	} `yaml:"add" json:"add"`
}

type Apps struct {
	Install struct {
		Brew Brew `yaml:"brew" json:"brew"`
		Gem  Gem  `yaml:"gem" json:"gem"`
		Npm  Npm  `yaml:"npm" json:"npm"`
		Mas  Mas  `yaml:"mas" json:"mas"`
	} `yaml:"install" json:"install"`
	Config struct {
		Mkdir         []string       `yaml:"mkdir" json:"mkdir"`
		SymbolicLinks []SymbolicLink `yaml:"symbolic_links" json:"symbolic_links"`
		Dockutil      Dockutil       `yaml:"dockutil" json:"dockutil"`
	} `yaml:"config" json:"config"`
}

func ResolveNumberOfParallelCommandExecution(ctx context.Context) int {
	// Get the number of logical CPUs
	numCPU := runtime.NumCPU()
	log.WithContext(ctx).Infof("number of logical CPUs: %d", numCPU)

	// Set the maximum number of CPUs that can execute simultaneously
	runtime.GOMAXPROCS(numCPU)

	// Set the number of goroutines that can run in parallel
	n := numCPU
	if n > MAX_ALLOWED_PARALLEL_INSTALLATION_PROCESSES {
		n = MAX_ALLOWED_PARALLEL_INSTALLATION_PROCESSES
	}

	log.WithContext(ctx).Infof("number of goroutines that can run in parallel: %d", n)

	return n
}

func ResolveCommandsFromFile(ctx context.Context, filename string) ([]Command, error) {
	// Read the YAML file
	file, err := os.ReadFile(filename)
	if err != nil {
		log.WithContext(ctx).WithField("filename", filename).WithError(err).Fatal("failed to read YAML file")
		return nil, fmt.Errorf("failed to read YAML file %v", err)
	}

	// Unmarshal the YAML file
	var apps Apps
	if err := yaml.Unmarshal(file, &apps); err != nil {
		log.WithContext(ctx).WithField("filename", filename).WithError(err).Fatal("failed to parse YAML file")
		return nil, fmt.Errorf("failed to parse YAML file %v", err)
	}

	// Configure the way commands will be executed
	var commands []Command

	// Homebrew related commands
	// TODO: think about executing brew update and brew tap before any other command
	// commands = append(commands, Command("brew update"))
	// for _, tap := range apps.Install.Brew.Taps {
	// 	commands = append(commands, Command(fmt.Sprintf("brew tap %s", tap)))
	// }
	for _, formula := range apps.Install.Brew.Formulas {
		commands = append(commands, Command(fmt.Sprintf("brew install %s --force --verbose", formula)))
	}
	for _, customCommand := range apps.Install.Brew.CustomCommands {
		commands = append(commands, Command(customCommand))
	}
	for _, cask := range apps.Install.Brew.Casks {
		commands = append(commands, Command(fmt.Sprintf("brew install --cask %s --force --verbose", cask)))
	}

	return commands, nil
}

func (c Command) Run(ctx context.Context, resultChan chan<- CommandResult) bool {
	select {
	case <-ctx.Done():
		resultChan <- CommandResult{
			Command:            c,
			Executed:           false,
			ExecutionTimeTs:    0,
			Success:            false,
			OperationCancelled: true,
			Message:            "cancelled by the user (ctrl+c).",
		}
		return false
	default:
		// Split the command into name and args
		args := strings.Fields(string(c))
		name := args[0]
		args = args[1:]

		// Setup the command
		cmd := cmd.NewCmd(name, args...)
		statusChan := cmd.Start() // non-blocking

		// Setup a ticker to print the last line of stdout every 2s
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		// Print last line of stdout every 2s
		go func() {
			for range ticker.C {
				select {
				case <-ctx.Done():
					log.WithContext(ctx).WithField("cmd", cmd.Status().Cmd).Info("stopped the command execution by the user (ctrl+c).")
					cmd.Stop()
					return
				default:
					status := cmd.Status()
					if len(status.Stdout) > 0 {
						log.WithContext(ctx).WithField("cmd", cmd.Status().Cmd).Info(status.Stdout[len(status.Stdout)-1])
					}
				}
			}
		}()

		// Block waiting for command to exit
		finalStatus := <-statusChan

		// Prepare result
		result := CommandResult{
			Command:            c,
			Executed:           finalStatus.Complete,
			ExecutionTimeTs:    finalStatus.StopTs - finalStatus.StartTs, // in milliseconds,
			Success:            finalStatus.Exit == 0,
			OperationCancelled: false,
			Message:            "",
		}

		// Get error message
		if len(finalStatus.Stderr) > 0 {
			result.Message = strings.Join(finalStatus.Stderr, "\n")
		}
		if !finalStatus.Complete {
			result.Message = "command stopped or signaled, see the logs for more info."
		}

		// Send result to channel
		resultChan <- result

		return result.Success
	}
}
