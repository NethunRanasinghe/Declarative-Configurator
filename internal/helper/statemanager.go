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

func GetPackageDifferences(configPackages AppPackages, statePackages AppPackages, changes *StateChanges){

	// Added (To-Install)
	for _, value := range configPackages.Native{
		if !(Contains(statePackages.Native, value)){
			changes.NativeToInstall = append(changes.NativeToInstall, value)
		}
	}

	for _, value := range configPackages.Flatpaks{
		if !(Contains(statePackages.Flatpaks, value)){
			changes.FlatpakToInstall = append(changes.FlatpakToInstall, value)
		}
	}

	for _, value := range configPackages.Local{
		if !(Contains(statePackages.Local, value)){
			changes.LocalToInstall = append(changes.LocalToInstall, value)
		}
	}

	// Removed (To-Remove)
	for _, value := range statePackages.Native{
		if !(Contains(configPackages.Native, value)){
			changes.NativeToRemove = append(changes.NativeToRemove, value)
		}
	}

	for _, value := range statePackages.Flatpaks{
		if !(Contains(configPackages.Flatpaks, value)){
			changes.FlatpakToRemove = append(changes.FlatpakToRemove, value)
		}
	}

	for _, value := range statePackages.Local{
		if !(Contains(configPackages.Local, value)){
			changes.LocalToRemove = append(changes.LocalToRemove, value)
		}
	}
}