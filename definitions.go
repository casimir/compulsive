package compulsive

type PackageState rune

const (
	StateUnknown  PackageState = '?'
	StateOutdated PackageState = '+'
	StateUpToDate PackageState = '='
)

type (
	Package interface {
		Name() string
		Label() string
		State() PackageState
		Version() string
		NextVersion() string
	}

	Provider interface {
		List() []Package
		UpdateCommand(...Package) string
	}
)
