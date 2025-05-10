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

	for _, details := range packageDetails {
		for _, pkg := range details.Native {
			if !(helper.Contains(allPackages.Native, pkg)) {
				allPackages.Native = append(allPackages.Native, pkg)
			}
		}

		for _, pkg := range details.Flatpaks {
			if !(helper.Contains(allPackages.Flatpaks, pkg)) {
				allPackages.Flatpaks = append(allPackages.Flatpaks, pkg)
			}
		}

		for _, pkg := range details.Local {
			if !(helper.Contains(allPackages.Local, pkg)) {
				allPackages.Local = append(allPackages.Local, pkg)
			}
		}
	}

	// Check State
	changes, hasChanged := helper.CheckState(allPackages)
	if !hasChanged {
		return
	}

	fmt.Println(changes)
}

func packageOps(pm packages.PackageManager, pkg string) {
	err := pm.Install(pkg)
	if err != nil {
		return
	}
}
