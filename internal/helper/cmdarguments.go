package helper

import "fmt"

func GetConfigPath(lptr *string, sptr *string) string {
	if *lptr == "" && *sptr == "" {
		fmt.Println("Warning : No config path found, using default './Config'")
		return "./Config"
	}

	configPath := *lptr
	if configPath == "" {
		configPath = *sptr
	}

	return configPath
}
