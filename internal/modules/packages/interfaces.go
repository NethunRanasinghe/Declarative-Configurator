package packages

type PackageManager interface{
	Install(pkg string) error
	Remove(pkg string) error
	Update() error
}