package helper

import (
	"fmt"
	"os"
)

const cmdFormat string = "%-20s : %2s\n"

var helpCommand = map[string]string{"refresh": "refresh all configs",
	"refresh all": "refresh all configs",
	"refresh <module_name> (Ex:- refresh packages)": "refresh specific config"}

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

func HandleCMDArgs(args []string, stringFormat string) int {
	if len(args) < 2 {
		fmt.Printf(stringFormat, "\nWarning", "No argument provided, use 'help' to view available commands!")
	}

	if args[1] == "help" {
		for key, val := range helpCommand {
			fmt.Printf(cmdFormat, key, val)
		}
		fmt.Print("\n\n")
		os.Exit(0)
	}

	if args[1] == "refresh" {
		if len(args) < 3 {
			return 0
		}

		if args[2] == "all" {
			return 1
		}

		if args[2] == "packages" {
			return 2
		}

		fmt.Printf(stringFormat, "\nError", "Invalid Sub Command!")
	}

	fmt.Printf(stringFormat, "\nError", "Invalid Sub Command!")
	return -1
}
