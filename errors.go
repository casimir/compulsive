package compulsive

import (
	"errors"
	"os"
	"strings"
)

var (
	ErrPackageName         = errors.New("not a package name")
	ErrPackageNotFound     = errors.New("package not found")
	ErrProviderNotFound    = errors.New("provider not found")
	ErrProviderUnavailable = errors.New("provider unavailable")
	ErrSudoNeeded          = errors.New("sudo needed for this operation")
)

func CheckSudo() error {
	if os.Getuid() > 0 {
		return ErrSudoNeeded
	}
	return nil
}

func IsPackageName(name string) bool {
	return len(strings.SplitN(name, "/", 2)) == 2
}
