package main

import (
	"declarative-configurator/internal/helper"
	"declarative-configurator/internal/modules/packages"
	"fmt"
	"log"
)

func main() {
	// Start Program
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
		fmt.Println("Packages : No Changed !")
		return
	}

	fmt.Println(changes)
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

func packageOps(pm packages.PackageManager, pkg string) {
	err := pm.Install(pkg)
	if err != nil {
		return
	}
}

//endregion
