package resource

type Manager interface {
	Create(options Options) error
	Delete(options Options) error
	SetDryRun()
	Update(options Options) error
}
