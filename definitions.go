package compulsive

type PackageState rune

const (
	StateUnknown  PackageState = '?'
	StateOutdated PackageState = '+'
	StateUpToDate PackageState = '='
)

type (
	Package interface {
		Provider() Provider
		Name() string
		Label() string
		State() PackageState
		Version() string
		NextVersion() string
	}

	Provider interface {
		Name() string
		IsAvailable() bool
		Sync() error
		List() ([]Package, error)
		UpdateCommand(...Package) string
	}
)
