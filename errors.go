package compulsive

import (
	"errors"
	"os"
	"strings"
)

var (
	PackageNameError         = errors.New("not a package name")
	PackageNotFoundError     = errors.New("package not found")
	ProviderNotFoundError    = errors.New("provider not found")
	ProviderUnavailableError = errors.New("provider unavailable")
	SudoNeededError          = errors.New("sudo needed for this operation")
)

func CheckSudo() error {
	if os.Getuid() > 0 {
		return SudoNeededError
	}
	return nil
}

func IsPackageName(name string) bool {
	return len(strings.SplitN(name, "/", 2)) == 2
}
