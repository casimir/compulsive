package providers

import (
	"github.com/casimir/compulsive"
	"github.com/casimir/compulsive/providers/golang"
	"github.com/casimir/compulsive/providers/homebrew"
	"github.com/casimir/compulsive/providers/pip"
)

type providerMapEntry struct {
	name          string
	newFunc       func() compulsive.Provider
	availableFunc func() bool
}

var providerMap = []providerMapEntry{
	{"Go", golang.New, golang.Available},
	{"Homebrew", homebrew.New, homebrew.Available},
	{"Pip", pip.NewV(""), pip.AvailableV("")},
	{"Pip2", pip.NewV("2"), pip.AvailableV("2")},
	{"Pip3", pip.NewV("3"), pip.AvailableV("3")},
}

type (
	ProviderEntry struct {
		Name      string
		Available bool
		Instance  compulsive.Provider
	}
	ProviderList []ProviderEntry
)

func New(name string) (compulsive.Provider, bool) {
	for _, it := range providerMap {
		if it.name == name && it.availableFunc() {
			return it.newFunc(), true
		}
	}
	return nil, false
}

func list(filterFunc func(providerMapEntry) bool, withInstance bool) ProviderList {
	var providers ProviderList
	for _, it := range providerMap {
		if filterFunc(it) {
			entry := ProviderEntry{
				Name:      it.name,
				Available: it.availableFunc(),
			}
			if entry.Available && withInstance {
				entry.Instance = it.newFunc()
			}
			providers = append(providers, entry)
		}
	}
	return providers
}

// ListAll gives the list of all providers along with their availability.
func ListAll(withInstance bool) ProviderList {
	filter := func(_ providerMapEntry) bool { return true }
	return list(filter, withInstance)
}

// ListAvailable gives the list of available providers.
func ListAvailable(withInstance bool) ProviderList {
	filter := func(entry providerMapEntry) bool {
		return entry.availableFunc()
	}
	return list(filter, withInstance)
}
