package helper

import (
	"github.com/goccy/go-yaml"
	"log"
	"os"
	"reflect"
)

type StateChanges struct {
	NativeToInstall  []string
	FlatpakToInstall []string
	LocalToInstall   []string
	NativeToRemove   []string
	FlatpakToRemove  []string
	LocalToRemove    []string
}

const StateFile string = ".state.yaml"

var StateDetails []byte

func CheckPackageState(appPackages AppPackages) (StateChanges, bool) {
	var stateChanges StateChanges
	var statePackages AppPackages

	// Check if state file exists
	_, err := os.Stat(StateFile)

	// If it doesn't exist
	if err != nil {
		CreateStateYaml()

		// Assign config packages to to-install
		stateChanges.NativeToInstall = appPackages.Native
		stateChanges.FlatpakToInstall = appPackages.Flatpaks
		stateChanges.LocalToInstall = appPackages.Local
	}

	// Refresh State
	RefreshState()

	// Create a map
	statePackagesMap := make(map[string]AppPackages)

	// Unmarshal the YAML into the map
	if err := yaml.Unmarshal(StateDetails, &statePackagesMap); err != nil {
		panic(err)
	}

	// Get Packages Data
	statePackages, ok1 := statePackagesMap["packages"]
	if !ok1 {
		log.Fatal("Error Reading State File !")
	}

	// Normalize both structs
	normalizePackageStructs(&statePackages)
	normalizePackageStructs(&appPackages)

	// Compare State and Config
	isEqual := reflect.DeepEqual(statePackages, appPackages)
	if isEqual {
		return stateChanges, false
	}

	getPackageDifferences(appPackages, statePackages, &stateChanges)
	return stateChanges, true
}

func StateManager(stateConfig StateConfig) error {
	var stateMap map[string]interface{}

	// Refresh State
	RefreshState()

	// Unmarshal YAML
	if err := yaml.Unmarshal(StateDetails, &stateMap); err != nil {
		return err
	}

	// Update Package Section
	if stateConfig.ModuleType == 0 {
		// Extract Section
		yamlPackagesSection := stateMap["packages"]
		yamlPackageSectionBytes, _ := yaml.Marshal(yamlPackagesSection)

		var yamlPackages AppPackages
		err := yaml.Unmarshal(yamlPackageSectionBytes, &yamlPackages)
		if err != nil {
			return err
		}

		// Modify
		stateManagerPackageModify(&yamlPackages, &stateConfig)

		// Recreate
		newYamlPackageSectionBytes, _ := yaml.Marshal(yamlPackages)
		var newUpdatedPackageMap map[string]interface{}
		err1 := yaml.Unmarshal(newYamlPackageSectionBytes, &newUpdatedPackageMap)
		if err1 != nil {
			return err1
		}
		stateMap["packages"] = newUpdatedPackageMap
	}

	// Write To File
	newUpdatedYaml, _ := yaml.Marshal(stateMap)

	// Write to file
	err2 := os.WriteFile(StateFile, newUpdatedYaml, 0644)
	if err2 != nil {
		return err2
	}

	return nil
}

//region Package Module

func stateManagerPackageModify(stateFilePackagesSection *AppPackages, stateConfig *StateConfig) {
	// Native - 0, Flatpak - 2, Local - -1
	// Add - 1, Remove - 0
	// **PMU** (3)

	// Native
	if stateConfig.PackageType == 0 {
		if stateConfig.AddOrRemove == 1 {
			stateFilePackagesSection.Native = append(stateFilePackagesSection.Native, stateConfig.PackageName)
		} else {
			stateFilePackagesSection.Native = RemoveFromSlice(stateFilePackagesSection.Native, stateConfig.PackageName)
		}
	}

	// Flatpak
	if stateConfig.PackageType == 2 {
		if stateConfig.AddOrRemove == 1 {
			stateFilePackagesSection.Flatpaks = append(stateFilePackagesSection.Flatpaks, stateConfig.PackageName)
		} else {
			stateFilePackagesSection.Flatpaks = RemoveFromSlice(stateFilePackagesSection.Flatpaks, stateConfig.PackageName)
		}
	}

	// Local
	if stateConfig.PackageType == -1 {
		if stateConfig.AddOrRemove == 1 {
			stateFilePackagesSection.Local = append(stateFilePackagesSection.Local, stateConfig.PackageName)
		} else {
			stateFilePackagesSection.Local = RemoveFromSlice(stateFilePackagesSection.Local, stateConfig.PackageName)
		}
	}
}

func getPackageDifferences(config AppPackages, state AppPackages, changes *StateChanges) {
	changes.NativeToInstall = diffAdd(config.Native, state.Native)
	changes.FlatpakToInstall = diffAdd(config.Flatpaks, state.Flatpaks)
	changes.LocalToInstall = diffAdd(config.Local, state.Local)

	changes.NativeToRemove = diffRemove(config.Native, state.Native)
	changes.FlatpakToRemove = diffRemove(config.Flatpaks, state.Flatpaks)
	changes.LocalToRemove = diffRemove(config.Local, state.Local)
}

// Added (To-Install)
func diffAdd(list1, list2 []string) []string {
	var result []string
	for _, val := range list1 {
		if !Contains(list2, val) {
			result = append(result, val)
		}
	}
	return result
}

// Removed (To-Remove)
func diffRemove(list1, list2 []string) []string {
	var result []string
	for _, val := range list2 {
		if !Contains(list1, val) {
			result = append(result, val)
		}
	}
	return result
}

func normalizePackageStructs(s *AppPackages) {
	if s.Native == nil {
		s.Native = []string{}
	}

	if s.Flatpaks == nil {
		s.Flatpaks = []string{}
	}

	if s.Local == nil {
		s.Local = []string{}
	}
}

//endregion
