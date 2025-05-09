package packages

import (
	"fmt"
	"os/exec"
	"os"
)

type DnfManager struct {}
const breakLine string = "--------------------------------------------------"

func (d DnfManager) Install(pkg string) error{

	// Command
	cmd := exec.Command("sudo", "dnf", "install", pkg, "-y")

	// Connect command output to current stdout/stderr
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

	fmt.Println(breakLine)
	fmt.Println("Installing ", pkg)
	fmt.Println(breakLine)

	// Run and check for errors
    if err := cmd.Run(); err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok{
			return fmt.Errorf("dnf install failed with exit code %d", exitError.ExitCode())
		}
        return fmt.Errorf("failed to execute dnf install: %w", err)
    }

    fmt.Println("Package(s) installed successfully.")
	return nil
}

func (d DnfManager) Remove(pkg string) error{

	// Command
	cmd := exec.Command("sudo", "dnf", "remove", pkg, "-y")

	// Connect command output to current stdout/stderr
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

	fmt.Println(breakLine)
	fmt.Println("Removing ", pkg)
	fmt.Println(breakLine)

	// Run and check for errors
    if err := cmd.Run(); err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok{
			return fmt.Errorf("dnf uninstall failed with exit code %d", exitError.ExitCode())
		}
        return fmt.Errorf("failed to execute dnf remove: %w", err)
    }

    fmt.Println("Package(s) removed successfully.")
	return nil
}

func (d DnfManager) Update() error{

	// Command
	cmd := exec.Command("sudo", "dnf", "update", "-y")

	// Connect command output to current stdout/stderr
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

	fmt.Println(breakLine)
	fmt.Println("Updating Packages")
	fmt.Println(breakLine)

	// Run and check for errors
    if err := cmd.Run(); err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok{
			return fmt.Errorf("dnf update failed with exit code %d", exitError.ExitCode())
		}
        return fmt.Errorf("failed to execute dnf update: %w", err)
    }

    fmt.Println("System updated successfully.")
	return nil
}