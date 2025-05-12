package packages

import (
	"bufio"
	"declarative-configurator/internal/helper"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
)

const BreakLine = "--------------------------------------------------"

// region Run Command

// Prepares the command arguments
func prepareCommandArgs(manager, action string, args ...string) []string {
	cmdArgs := append([]string{manager, action}, args...)
	if manager == "dnf" || manager == "flatpak" {
		cmdArgs = append(cmdArgs, "-y")
	}

	return cmdArgs
}

// Creates the appropriate command
func createCommand(manager string, cmdArgs []string, isRoot bool) *exec.Cmd {
	var cmd *exec.Cmd
	if isRoot {
		cmd = exec.Command("sudo", cmdArgs...)
	} else {
		// Temporary fix before a more good implementation
		tmpFix := cmdArgs[1:]
		cmd = exec.Command(manager, tmpFix...)
		fmt.Println("CMD : ", cmd.String())
	}
	return cmd
}

// Reads and prints output from a pipe, stripping ANSI escape sequences
func streamOutput(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		// Strip ANSI escape sequences
		line := stripAnsi(scanner.Text())
		fmt.Println(line)
	}
}

// Prints information about the package command
func printCommandInfo(action string, args ...string) {
	fmt.Println(BreakLine)
	if len(args) > 0 {
		fmt.Printf("%s %s\n", action, args[0])
	} else {
		fmt.Printf("%s Packages\n", action)
	}
	fmt.Println(BreakLine)
}

// Executes a package management command

func RunPackageCommand(manager string, action string, isRoot bool, args ...string) error {
	cmdArgs := prepareCommandArgs(manager, action, args...)
	cmd := createCommand(manager, cmdArgs, isRoot)

	// Create pipes to capture output
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Stream output in goroutines
	go streamOutput(stdoutPipe)
	go streamOutput(stderrPipe)

	// Print command information
	printCommandInfo(action, args...)

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return fmt.Errorf("%s %s failed with exit code %d", manager, action, exitError.ExitCode())
		}
		return fmt.Errorf("failed to execute %s %s: %w", manager, action, err)
	}

	fmt.Printf("%s %s completed successfully.\n", manager, action)
	return nil
}

// Removes ANSI escape sequences from a string
func stripAnsi(str string) string {
	ansiEscapeRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiEscapeRegex.ReplaceAllString(str, "")
}

// endregion

// Create state config

func CreateStateConfigHelper(pkg string, addOrRemove int, packageType int) helper.StateConfig {
	var stateConfig helper.StateConfig
	stateConfig.ModuleType = 0
	stateConfig.PackageType = packageType
	stateConfig.AddOrRemove = addOrRemove
	stateConfig.PackageName = pkg

	return stateConfig
}
