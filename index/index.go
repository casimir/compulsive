package index

import (
	"fmt"
	"sort"

	"github.com/casimir/compulsive"
	"github.com/casimir/compulsive/providers"
)

type Index map[providers.ProviderEntry]map[string]compulsive.Package

func (idx Index) FindProviderByName(name string) (providers.ProviderEntry, bool) {
	for pvd := range idx {
		if pvd.Name == name {
			return pvd, true
		}
	}
	return providers.ProviderEntry{}, false
}

func (idx Index) ListProviderPackages(providerName string) []compulsive.Package {
	var list []compulsive.Package
	provider, ok := idx.FindProviderByName(providerName)
	if !ok {
		return list
	}
	for _, pkg := range idx[provider] {
		list = append(list, pkg)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list
}

func NewFor(names []string, sync bool) (Index, error) {
	index := make(map[providers.ProviderEntry]map[string]compulsive.Package)
	sort.Strings(names)
	for _, pvd := range providers.ListAvailable() {
		searchIdx := sort.SearchStrings(names, pvd.Name)
		if searchIdx < len(names) && names[searchIdx] == pvd.Name {
			pvdIndex := make(map[string]compulsive.Package)
			if sync {
				if err := pvd.Instance.Sync(); err != nil {
					return index, fmt.Errorf("could not sync provider: %s", err)
				}
			}
			list, err := pvd.Instance.List()
			if err != nil {
				return index, err
			}
			for _, pkg := range list {
				pvdIndex[pkg.Name()] = pkg
			}
			index[pvd] = pvdIndex
		}
	}
	return index, nil
}

func New(sync bool) (Index, error) {
	var names []string
	for _, pvd := range providers.ListAvailable() {
		names = append(names, pvd.Name)
	}
	return NewFor(names, sync)
}
