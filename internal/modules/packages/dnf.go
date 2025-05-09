package packages

import (
	"fmt"
)

type DnfManager struct {}

func (d DnfManager) Install(pkg string) error{

	fmt.Println("Installing ", pkg)
	return nil
}