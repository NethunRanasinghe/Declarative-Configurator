package main

import (
	"fmt"
	"log"
	"declarative-configurator/internal/helper"
)

func main() {
	// Get OS Details
	var osDetails = helper.GetOsDetails()

	// Get Package Data from the YAML
	packageDetails, err := helper.GetPackageDetails("Config/package.yaml", osDetails.Distro)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Package Details for", osDetails.Distro)
	for _, details := range packageDetails {
		fmt.Printf("%+v\n", details)
	}
}