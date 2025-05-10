package main

import (
	"declarative-configurator/internal/helper"
	"declarative-configurator/internal/modules/packages"
	"fmt"
	"log"
)

func main() {
	// Start Program

	// Start Package Module
	startPackageModule()
}

//region Package Module

func startPackageModule() {
	// Get OS Details
	var osDetails = helper.GetOsDetails()

	// Get Package Data from the YAML
	packageDetails, err := helper.GetPackageDetails("Config/package.yaml", osDetails.Distro)
	if err != nil {
		log.Fatal(err)
	}

	// Merge
	var allPackages helper.AppPackages
	updateAllPackages(packageDetails, &allPackages)

	// Check State
	changes, hasChanged := helper.CheckState(allPackages)
	if !hasChanged {
		fmt.Println("Packages : No Changes !")
		return
	}

	// Separate Changes
	seperatedChanges := separateChanges(changes, &osDetails.Distro)

	// Perform Package Operations
	for pm, ch := range seperatedChanges {
		packageOps(pm, ch)
	}
}

func updateAllPackages(packageDetails []helper.AppPackages, allPackages *helper.AppPackages) {
	var nativeAll []string
	var flatpakAll []string
	var localAll []string

	for _, details := range packageDetails {
		nativeAll = append(nativeAll, mergePackages(details.Native, allPackages.Native)...)
		flatpakAll = append(flatpakAll, mergePackages(details.Flatpaks, allPackages.Flatpaks)...)
		localAll = append(localAll, mergePackages(details.Local, allPackages.Local)...)
	}

	allPackages.Native = nativeAll
	allPackages.Flatpaks = flatpakAll
	allPackages.Local = localAll
}

func mergePackages(pkg1, pkg2 []string) []string {
	var result []string
	for _, value := range pkg1 {
		if !(helper.Contains(pkg2, value)) {
			result = append(result, value)
		}
	}

	return result
}

func separateChanges(changes helper.StateChanges, distro *string) map[packages.PackageManager]helper.PackageOperation {
	var allChanges = map[packages.PackageManager]helper.PackageOperation{}

	// Separate Native Changes
	var nativeChanges helper.PackageOperation
	nativeChanges.Install = changes.NativeToInstall
	nativeChanges.Remove = changes.NativeToRemove
	allChanges[separatedChangesNative(distro)] = nativeChanges

	// Return all changes
	return allChanges
}

func separatedChangesNative(distro *string) packages.PackageManager {
	var pmTypeMap = map[string]int{"dnf": 0, "apt": 1}
	var pm packages.PackageManager

	pmType := pmTypeMap[helper.DistroAndPackageManager[*distro]]

	// Perform Package Operations
	if pmType == 0 {
		pm = packages.DnfManager{}
	} else {
		log.Fatal("Package Manager is not supported")
	}

	return pm
}

func packageOps(pm packages.PackageManager, changes helper.PackageOperation) {
	// Install Package
	for _, pkg := range changes.Install {
		err := pm.Install(pkg)
		if err != nil {
			return
		}
	}

	// Remove Package
	for _, pkg := range changes.Remove {
		err := pm.Remove(pkg)
		if err != nil {
			return
		}
	}
}

//endregion
