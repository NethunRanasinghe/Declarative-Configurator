package packages

import (
	"declarative-configurator/internal/helper"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const BreakLine = "--------------------------------------------------"

// Run Command

func RunPackageCommand(manager string, action string, isRoot bool, args ...string) error {
	cmdArgs := append([]string{manager, action}, args...)

	// Package manager based options
	if manager == "dnf" || manager == "flatpak" {
		cmdArgs = append(cmdArgs, "-y")
	}

	// If root add sudo
	var cmd *exec.Cmd
	if isRoot {
		cmd = exec.Command("sudo", cmdArgs...)
	} else {
		tmpFix := cmdArgs[1:] // Temporary fix before a more good implementation
		cmd = exec.Command(manager, tmpFix...)
		fmt.Println("CMD : ", cmd.String())
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println(BreakLine)
	if len(args) > 0 {
		fmt.Printf("%s %s\n", action, args[0])
	} else {
		fmt.Printf("%s Packages\n", action)
	}
	fmt.Println(BreakLine)

	if err := cmd.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return fmt.Errorf("%s %s failed with exit code %d", manager, action, exitError.ExitCode())
		}
		return fmt.Errorf("failed to execute %s %s: %w", manager, action, err)
	}

	fmt.Printf("%s %s completed successfully.\n", manager, action)
	return nil
}

// Create state config

func CreateStateConfigHelper(pkg string, addOrRemove int, packageType int) helper.StateConfig {
	var stateConfig helper.StateConfig
	stateConfig.ModuleType = 0
	stateConfig.PackageType = packageType
	stateConfig.AddOrRemove = addOrRemove
	stateConfig.PackageName = pkg

	return stateConfig
}
