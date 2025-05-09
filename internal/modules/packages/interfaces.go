package packages

type PackageManager interface{
	Install(pkg string) error
}