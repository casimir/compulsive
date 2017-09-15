package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/casimir/compulsive"
	"github.com/casimir/compulsive/providers"
)

var (
	aAll           = flag.Bool("a", false, "include up-to-date packages/unavailable providers")
	aListProviders = flag.Bool("l", false, "list providers")
	aProvider      = flag.String("p", "", "run only the given provider")
)

func formatPackageLine(pkg compulsive.Package) string {
	var line []string
	line = append(line, pkg.Label())
	if pkg.State() == compulsive.StateOutdated {
		line = append(line, "("+pkg.Version()+" → "+pkg.NextVersion()+")")
	} else {
		line = append(line, "("+pkg.Version()+")")
	}
	return strings.Join(line, " ")
}

func runListProviders() {
	for _, p := range providers.ListAll(false) {
		var line []string
		if *aAll {
			if p.Available {
				line = append(line, "*")
			} else {
				line = append(line, " ")
			}
		} else if !p.Available {
			continue
		}
		line = append(line, p.Name)
		fmt.Println(strings.Join(line, " "))
	}
}

func runProvider(name string) {
	provider, found := providers.New(name)
	if !found {
		os.Exit(1)
	}
	for _, it := range provider.List() {
		state := it.State()
		if *aAll {
			fmt.Printf("%c ", state)
		} else if state != compulsive.StateOutdated {
			continue
		}
		fmt.Println(formatPackageLine(it))
	}
}

func runListPackages() {
	for _, p := range providers.ListAvailable(true) {
		fmt.Println(p.Name)
		for _, it := range p.Instance.List() {
			state := it.State()
			sign := ' '
			if *aAll {
				sign = rune(state)
			} else if state != compulsive.StateOutdated {
				continue
			}
			fmt.Printf("%c %s\n", sign, formatPackageLine(it))
		}
	}
}

func main() {
	flag.Parse()
	if *aListProviders {
		runListProviders()
	} else if *aProvider != "" {
		runProvider(*aProvider)
	} else {
		runListPackages()
	}
}
