package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/casimir/compulsive"
	"github.com/casimir/compulsive/providers"
)

var (
	aAll      = flag.Bool("a", false, "include up-to-date packages/unavailable providers")
	aProvider = flag.String("p", "", "apply the command for this provider only")
)

type (
	options struct {
		all      bool
		provider string
	}

	command struct {
		help    string
		runFunc func(options, ...string)
	}
)

var commandMap = map[string]command{
	"packages":  {"list packages", runListPackages},
	"providers": {"list providers", runListProviders},
}

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

func runListProviders(opts options, _ ...string) {
	for _, p := range providers.ListAll(false) {
		var line []string
		if opts.all {
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

func runProvider(opts options, _ ...string) {
	provider, found := providers.New(opts.provider)
	if !found {
		os.Exit(1)
	}
	for _, it := range provider.List() {
		state := it.State()
		if opts.all {
			fmt.Printf("%c ", state)
		} else if state != compulsive.StateOutdated {
			continue
		}
		fmt.Println(formatPackageLine(it))
	}
}

func runListPackages(opts options, args ...string) {
	if opts.provider != "" {
		runProvider(opts)
		return
	}

	for _, p := range providers.ListAvailable(true) {
		fmt.Println(p.Name)
		for _, it := range p.Instance.List() {
			state := it.State()
			sign := ' '
			if opts.all {
				sign = rune(state)
			} else if state != compulsive.StateOutdated {
				continue
			}
			fmt.Printf("%c %s\n", sign, formatPackageLine(it))
		}
	}
}

func print_usage() {
	var commands sort.StringSlice
	for it := range commandMap {
		commands = append(commands, it)
	}
	commands.Sort()

	fmt.Fprintf(os.Stderr, "Usage: %s [flags...] <command> [<parameters>...]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Commands:")
	for _, it := range commands {
		fmt.Fprintf(os.Stderr, "  %s\t%s\n", it, commandMap[it].help)
	}
}

func main() {
	flag.Usage = print_usage
	flag.Parse()

	commandName := "packages"
	args := flag.Args()
	if flag.NArg() > 0 {
		commandName = args[0]
		args = args[1:]
	}
	command, ok := commandMap[commandName]
	if ok {
		opts := options{all: *aAll, provider: *aProvider}
		command.runFunc(opts, args...)
	} else {
		fmt.Fprintf(os.Stderr, "error: unknown command: %s\n", commandName)
		os.Exit(1)
	}
}
