package helper

import (
	"github.com/goccy/go-yaml"
	"os"
	"log"
	"reflect"
)

type StateChanges struct{
	NativeToInstall []string
	FlatpakToInstall []string
	LocalToInstall []string
	NativeToRemove []string
	FlatpakToRemove []string
	LocalToRemove []string
}

const stateFile string = ".state.yaml"

func CheckState(appPackages AppPackages)(StateChanges, bool){
	var stateChanges StateChanges
	var statePackages AppPackages

	// Check if state file exists
	_, err := os.Stat(stateFile)

	// If it doesn't exist
	if err != nil{
		// Create StateFile
		_, err := os.Create(stateFile)
		if err != nil{
			log.Panic(err)
		}

		// Assign config packages to to-install
		stateChanges.NativeToInstall = appPackages.Native
		stateChanges.FlatpakToInstall = appPackages.Flatpaks
		stateChanges.LocalToInstall = appPackages.Local

		return stateChanges, true
	}

	// If it exist
	stateDetails, err := os.ReadFile(stateFile)
	if err != nil{
		panic(err)
	}

	// Create a map
	statePackagesMap := make(map[string]AppPackages)

	// Unmarshal the YAML into the map
	if err := yaml.Unmarshal(stateDetails, &statePackagesMap); err != nil {
		panic(err)
	}

	// Get Packages Data
	statePackages, ok1 := statePackagesMap["Packages"]
	if !ok1 {
		log.Fatal("Error Reading State File !")
	}

	// Compare State and Config
	isEqual := reflect.DeepEqual(statePackages, appPackages)
	if isEqual{
		return stateChanges, false
	}

	GetPackageDifferences(statePackages, appPackages, &stateChanges)
	return stateChanges, true
}

// Differences - Package ------ Start

func GetPackageDifferences(config AppPackages, state AppPackages, changes *StateChanges) {
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

// Differences - Package ------ End