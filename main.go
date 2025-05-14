package main

import (
	"declarative-configurator/internal/helper"
	"declarative-configurator/internal/modules/packages"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	const stringFormat string = "%-8s : %2s\n"

	// Handle CMD Arguments
	helper.HandleCMDArgs(os.Args, stringFormat)

	// CMD Flags
	configPtr := flag.String("config", "", "Config File Location")
	configSPtr := flag.String("c", "", "Config File Location")

	// Get Config
	flag.Parse()
	configPath := helper.GetConfigPath(configPtr, configSPtr)

	// Get OS Details
	var osDetails = helper.GetOsDetails()

	// Start Program
	showWelcome(stringFormat, osDetails)
	result := helper.HandleCMDArgs(os.Args, stringFormat)

	if result == -1 {
		os.Exit(1)
	}

	if result == 0 || result == 1 {
		fmt.Printf(stringFormat, "\nInfo", "Refreshing all!")
		startPackageModule(configPath, osDetails)
	}

	if result == 2 {
		fmt.Printf(stringFormat, "\nInfo", "Refreshing packages!")
		startPackageModule(configPath, osDetails)
	}

}

func showWelcome(stringFormat string, osDetails helper.OsDetailsObject) {
	fmt.Println("----------------------------------Welcome----------------------------------")
	fmt.Printf(stringFormat, "OS", osDetails.Os)
	fmt.Printf(stringFormat, "Distro", osDetails.Distro)
	fmt.Printf(stringFormat, "Base", osDetails.Base)
	fmt.Printf(stringFormat, "Arch", osDetails.Arch)
	fmt.Printf(stringFormat, "Hostname", osDetails.Hostname)
}

//region Package Module

func startPackageModule(configDirLoc string, osDetails helper.OsDetailsObject) {

	fmt.Println("\nPackage Module : Start")

	// Set Config Path
	packageConfigPath := fmt.Sprintf("%s/%s", configDirLoc, "package.yaml")

	// Get Package Data from the YAML
	packageDetails, err := helper.GetPackageDetails(packageConfigPath, osDetails.Distro)
	if err != nil {
		log.Fatal(err)
	}

	// Merge
	var allPackages helper.AppPackages
	updateAllPackages(packageDetails, &allPackages)

	// Check State
	changes, hasChanged := helper.CheckPackageState(allPackages)
	if !hasChanged {
		fmt.Println("\nPackages : No Changes !")
		return
	}

	// Separate Changes
	seperatedChanges := separateChanges(changes, &osDetails.Distro)

	// Perform Package Operations
	for pm, ch := range seperatedChanges {
		packageOps(pm, ch)
	}

	fmt.Println("\nPackage Module : End")
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

	// **PMU** (1)

	// Separate Native Changes
	var nativeChanges helper.PackageOperation
	nativeChanges.Install = changes.NativeToInstall
	nativeChanges.Remove = changes.NativeToRemove
	allChanges[separatedChangesNativeOrLocal(distro, false)] = nativeChanges

	// Separate Flatpak Changes
	var flatpakChanges helper.PackageOperation
	flatpakChanges.Install = changes.FlatpakToInstall
	flatpakChanges.Remove = changes.FlatpakToRemove
	allChanges[separateChangesSandboxed("flatpak")] = flatpakChanges

	// Separate Local Changes
	var localChanges helper.PackageOperation
	localChanges.Install = changes.LocalToInstall
	localChanges.Remove = changes.LocalToRemove
	allChanges[separatedChangesNativeOrLocal(distro, true)] = localChanges

	// Return all changes
	return allChanges
}

func separatedChangesNativeOrLocal(distro *string, localOrNot bool) packages.PackageManager {
	var pmTypeMap = map[string]int{"dnf": 0, "apt": 1} // **PMU** (2)
	var pm packages.PackageManager

	pmType := pmTypeMap[helper.DistroAndPackageManager[*distro]]

	// Perform Package Operations
	if pmType == 0 {
		if !localOrNot {
			pm = packages.DnfManager{}
		} else {
			pm = packages.LocalManager{}
		}
	} else {
		log.Fatal("Package Manager is not supported")
	}

	return pm
}

func separateChangesSandboxed(sandboxType string) packages.PackageManager {
	var sandboxTypeMap = map[string]int{"flatpak": 0} // **PMU** (2)
	var spm packages.PackageManager

	spmType := sandboxTypeMap[helper.DistroAndPackageManager[sandboxType]]

	// Perform Package Operations
	if spmType == 0 {
		spm = packages.FlatpakManager{}
	} else {
		log.Fatal("Sandbox Package Manager is not supported")
	}

	return spm
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
