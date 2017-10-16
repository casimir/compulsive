package compulsive

type PackageState rune

const (
	StateUnknown  PackageState = '?'
	StateOutdated PackageState = '+'
	StateUpToDate PackageState = '='
)

type (
	Package struct {
		Provider    Provider
		Name        string
		Label       string
		Summary     string
		Binaries    []string
		State       PackageState
		Version     string
		NextVersion string
	}

	Provider interface {
		Name() string
		IsAvailable() bool
		Sync() error
		List() ([]Package, error)
		UpdateCommand(...Package) string
	}
)
