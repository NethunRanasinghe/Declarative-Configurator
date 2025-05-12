package packages

import (
	"declarative-configurator/internal/helper"
)

type FlatpakManager struct{}

func (f FlatpakManager) runFlatpakCommand(action string, args ...string) error {
	return RunPackageCommand("flatpak", action, false, args...)
}

func (f FlatpakManager) Install(pkg string) error {
	err := f.runFlatpakCommand("install", pkg)
	if err != nil {
		return err
	}

	installStateConfig := CreateStateConfigHelper(pkg, 1, 2)
	err = helper.StateManager(installStateConfig)
	if err != nil {
		return err
	}
	return nil
}

func (f FlatpakManager) Remove(pkg string) error {
	err := f.runFlatpakCommand("uninstall", pkg)
	if err != nil {
		return err
	}

	removeStateConfig := CreateStateConfigHelper(pkg, 0, 2)
	err = helper.StateManager(removeStateConfig)
	if err != nil {
		return err
	}
	return nil
}

func (f FlatpakManager) Update() error {
	return f.runFlatpakCommand("update")
}
