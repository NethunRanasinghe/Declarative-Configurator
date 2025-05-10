package packages

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type DnfManager struct{}

const breakLine = "--------------------------------------------------"

func (d DnfManager) runDnfCommand(action string, args ...string) error {
	// Construct the full command
	cmdArgs := append([]string{"dnf", action}, args...)
	cmdArgs = append(cmdArgs, "-y")
	cmd := exec.Command("sudo", cmdArgs...)

	// Connect output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println(breakLine)
	if len(args) > 0 {
		fmt.Printf("%s %s\n", action, args[0])
	} else {
		fmt.Printf("%s Packages\n", action)
	}
	fmt.Println(breakLine)

	// Execute command
	if err := cmd.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return fmt.Errorf("dnf %s failed with exit code %d", action, exitError.ExitCode())
		}
		return fmt.Errorf("failed to execute dnf %s: %w", action, err)
	}

	fmt.Printf("dnf %s completed successfully.\n", action)
	return nil
}

func (d DnfManager) Install(pkg string) error {
	return d.runDnfCommand("install", pkg)
}

func (d DnfManager) Remove(pkg string) error {
	return d.runDnfCommand("remove", pkg)
}

func (d DnfManager) Update() error {
	return d.runDnfCommand("update")
}
