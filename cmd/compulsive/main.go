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
		var line []string
		state := it.State()
		if *aAll {
			line = append(line, string(state))
		} else if state != compulsive.StateOutdated {
			continue
		}
		line = append(line, it.Label())
		if state == compulsive.StateOutdated {
			line = append(line, "("+it.Version()+" → "+it.NextVersion()+")")
		} else {
			line = append(line, "("+it.Version()+")")
		}
		fmt.Println(strings.Join(line, " "))
	}
}

func runListPackages() {
	for _, p := range providers.ListAvailable(true) {
		fmt.Println(p.Name)
		for _, it := range p.Instance.List() {
			var line []string
			state := it.State()
			sign := ' '
			if *aAll {
				sign = rune(state)
			} else if state != compulsive.StateOutdated {
				continue
			}
			line = append(line, string(sign))
			line = append(line, it.Label())
			if state == compulsive.StateOutdated {
				line = append(line, "("+it.Version()+" → "+it.NextVersion()+")")
			} else {
				line = append(line, "("+it.Version()+")")
			}
			fmt.Println(strings.Join(line, " "))
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
