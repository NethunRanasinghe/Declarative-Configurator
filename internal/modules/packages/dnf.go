package packages

import (
	"declarative-configurator/internal/helper"
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
	err := d.runDnfCommand("install", pkg)
	if err != nil {
		return err
	}

	installStateConfig := createStateConfigHelper(pkg, 1)
	err = helper.StateManager(installStateConfig)
	if err != nil {
		return err
	}
	return nil
}

func (d DnfManager) Remove(pkg string) error {
	err := d.runDnfCommand("remove", pkg)
	if err != nil {
		return err
	}

	removeStateConfig := createStateConfigHelper(pkg, 0)
	err = helper.StateManager(removeStateConfig)
	if err != nil {
		return err
	}
	return nil
}

func (d DnfManager) Update() error {
	return d.runDnfCommand("update")
}

func createStateConfigHelper(pkg string, addOrRemove int) helper.StateConfig {
	var stateConfig helper.StateConfig
	stateConfig.ModuleType = 0
	stateConfig.PackageType = 0
	stateConfig.AddOrRemove = addOrRemove
	stateConfig.PackageName = pkg

	return stateConfig
}
