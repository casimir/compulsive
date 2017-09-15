package providers

import (
	"sync"

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

func insertProvider(list ProviderList, idx int, withInstance bool, entry providerMapEntry) {
	provider := ProviderEntry{
		Name:      entry.name,
		Available: entry.availableFunc(),
	}
	if provider.Available && withInstance {
		provider.Instance = entry.newFunc()
	}
	list[idx] = provider
}

func list(filterFunc func(providerMapEntry) bool, withInstance bool) ProviderList {
	var filteredProviderMap []providerMapEntry
	for _, it := range providerMap {
		if filterFunc(it) {
			filteredProviderMap = append(filteredProviderMap, it)
		}
	}
	providers := make(ProviderList, len(filteredProviderMap))
	var wg sync.WaitGroup
	for i, it := range filteredProviderMap {
		if withInstance {
			wg.Add(1)
			go func(idx int, entry providerMapEntry) {
				defer wg.Done()
				insertProvider(providers, idx, withInstance, entry)
			}(i, it)
		} else {
			insertProvider(providers, i, withInstance, it)
		}
	}
	wg.Wait()
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
