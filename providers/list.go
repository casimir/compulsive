package providers

import (
	"sort"

	"github.com/casimir/compulsive"
)

func initInstances() map[string]compulsive.Provider {
	instances := []compulsive.Provider{
		NewGo(),
		NewBrew(),
		NewCargo(),
		NewPip(""),
		NewPip("2"),
		NewPip("3"),
	}
	instanceMap := make(map[string]compulsive.Provider, len(instances))
	for _, it := range instances {
		instanceMap[it.Name()] = it
	}
	return instanceMap
}

var Instances = initInstances()

func Check(name string) error {
	pvd, ok := Instances[name]
	if !ok {
		return compulsive.ErrProviderNotFound
	}
	if !pvd.IsAvailable() {
		return compulsive.ErrProviderUnavailable
	}
	return nil
}

func list(filterFunc func(compulsive.Provider) bool) []compulsive.Provider {
	var pvds []compulsive.Provider
	for _, pvd := range Instances {
		if filterFunc == nil || filterFunc(pvd) {
			pvds = append(pvds, pvd)
		}
	}
	sort.Slice(pvds, func(i, j int) bool { return pvds[i].Name() < pvds[j].Name() })
	return pvds
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
