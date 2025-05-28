package packages

import (
	"declarative-configurator/internal/helper"
)

type AptManager struct{}

func (a AptManager) runAptCommand(action string, args ...string) error {
	return RunPackageCommand("apt", action, true, args...)
}

func (a AptManager) Install(pkg string) error {
	err := a.runAptCommand("install", pkg)
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

func (a AptManager) Remove(pkg string) error {
	err := a.runAptCommand("remove", pkg)
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

func (a AptManager) Update() error {
	return a.runAptCommand("update")
}
