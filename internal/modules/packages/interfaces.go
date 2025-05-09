package packages

type PackageManager interface{
	Install(pkgs []string) error
	Update() error
	Remove(pkgs []string) error
}