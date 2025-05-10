package helper

import (
	"github.com/goccy/go-yaml"
	"log"
	"os"
	"slices"
)

type AppPackages struct {
	Native   []string `yaml:"Native"`
	Flatpaks []string `yaml:"Flatpaks"`
	Local    []string `yaml:"Local"`
}

type PackageOperation struct {
	Install []string
	Remove  []string
}

type StateConfig struct {
	ModuleType  int
	PackageType int
	AddOrRemove int
	PackageName string
}

var DistroAndPackageManager = map[string]string{
	"fedora":    "dnf",
	"rhel":      "dnf",
	"centos":    "dnf",
	"rocky":     "dnf",
	"almalinux": "dnf",
	"ubuntu":    "apt",
	"debian":    "apt",
}

type StateTemplate struct {
	Packages AppPackages `yaml:"packages"`
}

func Contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

func CreateStateYaml() {
	stateTemplate := StateTemplate{
		Packages: AppPackages{
			Native:   []string{},
			Flatpaks: []string{},
			Local:    []string{},
		},
	}

	// Marshal the state to YAML
	data, err := yaml.Marshal(stateTemplate)
	if err != nil {
		log.Panic(err)
	}

	// Write to file
	err2 := os.WriteFile(StateFile, data, 0644)
	if err2 != nil {
		log.Panic(err2)
	}
}

func ReadFileData(filePath string) []byte {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}
	return fileData
}

func RefreshState() {
	StateDetails = ReadFileData(StateFile)
}

func RemoveFromSlice(s []string, item string) []string {
	for index, value := range s {
		if value == item {
			return append(s[:index], s[index+1:]...)
		}
	}

	return s
}
