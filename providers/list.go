package providers

import (
	"github.com/casimir/compulsive"
	"github.com/casimir/compulsive/providers/golang"
	"github.com/casimir/compulsive/providers/homebrew"
)

type providerMapEntry struct {
	name          string
	newFunc       func() compulsive.Provider
	availableFunc func() bool
}

var providerMap = []providerMapEntry{
	{"Go", golang.New, golang.Available},
	{"Homebrew", homebrew.New, homebrew.Available},
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

func ListAll(withInstance bool) ProviderList {
	filter := func(_ providerMapEntry) bool { return true }
	return list(filter, withInstance)
}

func ListAvailable(withInstance bool) ProviderList {
	filter := func(entry providerMapEntry) bool {
		return entry.availableFunc()
	}
	return list(filter, withInstance)
}
