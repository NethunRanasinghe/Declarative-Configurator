package packages

import (
	"declarative-configurator/internal/helper"
)

type DnfManager struct{}

func (d DnfManager) runDnfCommand(action string, args ...string) error {
	return RunPackageCommand("dnf", action, true, args...)
}

func (d DnfManager) Install(pkg string) error {
	err := d.runDnfCommand("install", pkg)
	if err != nil {
		return err
	}

	installStateConfig := CreateStateConfigHelper(pkg, 1, 0)
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

	removeStateConfig := CreateStateConfigHelper(pkg, 0, 0)
	err = helper.StateManager(removeStateConfig)
	if err != nil {
		return err
	}
	return nil
}

func (d DnfManager) Update() error {
	return d.runDnfCommand("update")
}
