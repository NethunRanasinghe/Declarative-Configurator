package osinfo

import (
	"os"
	"runtime"
	"strings"
)

type OsDetailsObject struct{
	Arch string
	Os string
	Distro string
	Base string
	Hostname string
}

func GetOsDetails() OsDetailsObject{
	var details OsDetailsObject
	var distro string = "Unknown"
	var base string = "Unknown"

	// Hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	// Distro Details
	releaseDetails,err := os.ReadFile("/etc/os-release")

	if err == nil {
		releaseDetailsSplit := strings.Split(string(releaseDetails), "\n")
		for _, releaseDetail := range releaseDetailsSplit{
			releaseDetailParts := strings.SplitN(string(releaseDetail), "=", 2)

			if len(releaseDetailParts) != 2{
				continue
			}

			key := releaseDetailParts[0]
			value := strings.Trim(releaseDetailParts[1], `"`)

			if key == "ID"{
				distro = value
			}

			if key == "ID_LIKE"{
				base = value
			}
		}
	}

	if distro != "Unknown" && base == "Unknown"{
		base = distro
	}

	// Architecture
	architecture := runtime.GOARCH

	// OS
	operatingSystem := runtime.GOOS

	// Assign to struct
	details.Arch = architecture
	details.Os = operatingSystem
	details.Distro = distro
	details.Base = base
	details.Hostname = hostname

	return details
}