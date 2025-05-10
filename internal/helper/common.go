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
	err = os.WriteFile(StateFile, data, 0644)
	if err != nil {
		log.Panic(err)
	}
}
