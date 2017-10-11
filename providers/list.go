package providers

import (
	"github.com/casimir/compulsive"
	"github.com/casimir/compulsive/providers/golang"
	"github.com/casimir/compulsive/providers/homebrew"
	"github.com/casimir/compulsive/providers/pip"
)

var providers = []compulsive.Provider{
	golang.New(),
	homebrew.New(),
	pip.New(""),
	pip.New("2"),
	pip.New("3"),
}

func Check(name string) error {
	var provider compulsive.Provider
	for _, pvd := range providers {
		if pvd.Name() == name {
			provider = pvd
			break
		}
	}
	if provider == nil {
		return compulsive.ErrProviderNotFound
	}
	if !provider.IsAvailable() {
		return compulsive.ErrProviderUnavailable
	}
	return nil
}

func list(filterFunc func(compulsive.Provider) bool) []compulsive.Provider {
	var providerList []compulsive.Provider
	for _, pvd := range providers {
		if filterFunc == nil || filterFunc(pvd) {
			providerList = append(providerList, pvd)
		}
	}
	return providerList
}

// ListAll gives the list of all providers along with their availability.
func ListAll() []compulsive.Provider {
	return list(nil)
}

// ListAvailable gives the list of available providers.
func ListAvailable() []compulsive.Provider {
	filter := func(pvd compulsive.Provider) bool {
		return pvd.IsAvailable()
	}
	return list(filter)
}
