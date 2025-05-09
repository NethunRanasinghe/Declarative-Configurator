package main

import (
	"log"
	"declarative-configurator/internal/helper"
	"declarative-configurator/internal/modules/packages"
)

func main() {
	// Start Program
	startProgram()
}

func startProgram(){
	// Get OS Details
	var osDetails = helper.GetOsDetails()

	// Get Package Data from the YAML
	packageDetails, err := helper.GetPackageDetails("Config/package.yaml", osDetails.Distro)
	if err != nil {
		log.Fatal(err)
	}

	// Install Package
	for _, details := range packageDetails {
		for _, pkgDetails := range details.Native{
			
			var pm packages.PackageManager = packages.DnfManager{}
			packageOps(pm, pkgDetails)
		}
	}
}

func packageOps(pm packages.PackageManager, pkg string){
	pm.Install(pkg)
}