package application

type Installer interface {
	Install(options Options) error
	Uninstall(options Options) error
	SetDryRun()
}
