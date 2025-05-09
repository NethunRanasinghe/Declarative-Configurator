package helper

import (
	"slices"
)

type AppPackages struct {
	Native   []string `yaml:"Native"`
	Flatpaks []string `yaml:"Flatpaks"`
	Local    []string `yaml:"Local"`
}

func Contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}