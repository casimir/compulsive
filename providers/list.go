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

type ProviderEntry struct {
	Name      string
	Available bool
	Instance  compulsive.Provider
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
		return compulsive.ProviderNotFoundError
	}
	if !provider.IsAvailable() {
		return compulsive.ProviderUnavailableError
	}
	return nil
}

func list(filterFunc func(compulsive.Provider) bool) []ProviderEntry {
	var providerList []ProviderEntry
	for _, pvd := range providers {
		if filterFunc == nil || filterFunc(pvd) {
			provider := ProviderEntry{
				Name:      pvd.Name(),
				Available: pvd.IsAvailable(),
			}
			if provider.Available {
				provider.Instance = pvd
			}
			providerList = append(providerList, provider)
		}
	}
	return providerList
}

// ListAll gives the list of all providers along with their availability.
func ListAll() []ProviderEntry {
	return list(nil)
}

// ListAvailable gives the list of available providers.
func ListAvailable() []ProviderEntry {
	filter := func(pvd compulsive.Provider) bool {
		return pvd.IsAvailable()
	}
	return list(filter)
}
