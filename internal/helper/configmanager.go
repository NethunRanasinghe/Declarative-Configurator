package helper

import (
	"github.com/goccy/go-yaml"
	"os"
	"fmt"
)

func GetPackageDetails(configPath string, baseOs string) ([]AppPackages, error){
	var allPackageDetails []AppPackages

	packageDetails, err := os.ReadFile(configPath)
	if err != nil{
		panic(err)
	}

	// Create a map
	osPackages := make(map[string]AppPackages)

	// Unmarshal the YAML into the map
	if err := yaml.Unmarshal(packageDetails, &osPackages); err != nil {
		panic(err)
	}

	// Get Packages Data
	osPackageDetails, ok1 := osPackages[baseOs]
	anyPackageDetails, ok2 := osPackages["any"]

	if ok1{
		allPackageDetails = append(allPackageDetails, osPackageDetails)
	}

	if ok2{
		allPackageDetails = append(allPackageDetails, anyPackageDetails)
	}

	if len(allPackageDetails) == 0 {
		return nil, fmt.Errorf("no package information found for OS: %s", baseOs)
	}

	return allPackageDetails, nil
}