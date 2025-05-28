package packages

import (
	"declarative-configurator/internal/helper"
	"fmt"
	"path/filepath"
	"strings"
)

type LocalManager struct{}

var LocalInstaller string

const cmdFormat string = "%-20s : %2s\n"

func (l LocalManager) runLocalCommand(action string, args ...string) error {
	refreshLocalInstaller()
	return RunPackageCommand(LocalInstaller, action, true, args...)
}

func (l LocalManager) Install(pkg string) error {
	refreshLocalInstaller()
	err := l.runLocalCommand("install", pkg)
	if err != nil {
		return err
	}

	installStateConfig := CreateStateConfigHelper(pkg, 1, -1)
	err = helper.StateManager(installStateConfig)
	if err != nil {
		return err
	}
	return nil
}

func (l LocalManager) Remove(pkg string) error {
	formattedPkg := formatPackageNameForUninstall(pkg)
	err := l.runLocalCommand("remove", formattedPkg)
	if err != nil {
		return err
	}

	removeStateConfig := CreateStateConfigHelper(pkg, 0, -1)
	err = helper.StateManager(removeStateConfig)
	if err != nil {
		return err
	}
	return nil
}

func (l LocalManager) Update() error {
	// In the future
	// Get the current version
	// Get the version in yaml
	// Perform in-place update

	//  Current - Temp
	fmt.Printf(cmdFormat, "Info", "To Update Local Packages, Change the Current Installed Package Version in the YAML and Refresh !")
	return nil
}

func refreshLocalInstaller() {
	LocalInstaller = helper.DistroAndLocalInstaller[helper.GetOsDetails().Distro]
}

func formatPackageNameForUninstall(pkg string) string {
	// Extract the base file name
	base := filepath.Base(pkg)

	// Remove the extension
	extIndex := strings.LastIndex(base, ".")
	if extIndex == -1 {
		return base
	}

	return base[:extIndex]
}
