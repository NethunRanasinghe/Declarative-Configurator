package main

import (
	"declarative-configurator/internal/helper"
	"fmt"
)

func main(){
	var osDetails = osinfo.GetOsDetails()
	fmt.Println("Arch    : ",osDetails.Arch)
	fmt.Println("Os    : ", osDetails.Os)
	fmt.Println("Distro    : ", osDetails.Distro)
	fmt.Println("Distro Base    : ", osDetails.Base)
	fmt.Println("Hostname    : ", osDetails.Hostname)
}