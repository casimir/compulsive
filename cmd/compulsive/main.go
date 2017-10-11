package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/casimir/compulsive"
	"github.com/casimir/compulsive/index"
	"github.com/casimir/compulsive/providers"
)

type (
	command struct {
		help    string
		runFunc func(options, ...string) error
	}

	options struct {
		all      bool
		provider string
		sync     bool
	}
)

var (
	commandMap = map[string]command{
		"info":      {"\tprint detailed information about one or more packages", runInfoPackage},
		"packages":  {"list packages (default)", runListPackages},
		"providers": {"list providers", runListProviders},
	}
	cliOpts options
)

func init() {
	flag.BoolVar(&cliOpts.all, "a", false, "include up-to-date packages/unavailable providers")
	flag.StringVar(&cliOpts.provider, "p", "", "apply the command for this provider only")
	flag.BoolVar(&cliOpts.sync, "s", false, "sync providers before listing packages")
}

func runListProviders(opts options, _ ...string) error {
	for _, pvd := range providers.ListAll() {
		var line []string
		if opts.all {
			if pvd.IsAvailable() {
				line = append(line, "*")
			} else {
				line = append(line, " ")
			}
		} else if !pvd.IsAvailable() {
			continue
		}
		line = append(line, pvd.Name())
		fmt.Println(strings.Join(line, " "))
	}
	return nil
}

func runInfoPackage(opts options, args ...string) error {
	if len(args) != 1 || !compulsive.IsPackageName(args[0]) {
		return compulsive.ErrPackageName
	}
	parts := strings.SplitN(args[0], "/", 2)
	idx, err := index.NewFor([]string{parts[0]}, opts.sync)
	if err != nil {
		return fmt.Errorf("could not build index: %s", err)
	}
	pvd, ok := idx.FindProviderByName(parts[0])
	if !ok {
		return compulsive.ErrPackageNotFound
	}
	pkg, ok := idx[pvd][parts[1]]
	if !ok {
		return compulsive.ErrPackageNotFound
	}
	fmt.Print(compulsive.FmtPkgDesc(pkg))
	return nil
}

func runProvider(opts options, _ ...string) error {
	if err := providers.Check(opts.provider); err != nil {
		return err
	}
	idx, err := index.NewFor([]string{opts.provider}, opts.sync)
	if err != nil {
		return fmt.Errorf("could not build index: %s", err)
	}
	for _, it := range idx.ListProviderPackages(opts.provider) {
		state := it.State()
		if opts.all {
			fmt.Printf("%c ", state)
		} else if state != compulsive.StateOutdated {
			continue
		}
		fmt.Println(compulsive.FmtPkgLine(it))
	}
	return nil
}

func runListPackages(opts options, args ...string) error {
	if opts.provider != "" {
		return runProvider(opts)
	}
	idx, err := index.New(opts.sync)
	if err != nil {
		return fmt.Errorf("could not build index: %s", err)
	}
	for _, pvd := range providers.ListAvailable() {
		for _, it := range idx.ListProviderPackages(pvd.Name()) {
			state := it.State()
			if opts.all {
				fmt.Printf("%c %s\n", state, compulsive.FmtPkgLine(it))
			} else if state == compulsive.StateOutdated {
				fmt.Println(compulsive.FmtPkgLine(it))
			}
		}
	}
	return nil
}

func printUsage() {
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
	fmt.Fprintln(os.Stderr, "  help\t\tprint the current message")
	for _, it := range commands {
		fmt.Fprintf(os.Stderr, "  %s\t%s\n", it, commandMap[it].help)
	}
}

func main() {
	flag.Usage = printUsage
	flag.Parse()

	commandName := "packages"
	args := flag.Args()
	if flag.NArg() > 0 {
		commandName = args[0]
		args = args[1:]
	}

	var err error
	if command, ok := commandMap[commandName]; ok {
		err = command.runFunc(cliOpts, args...)
	} else if commandName == "help" {
		printUsage()
	} else {
		err = fmt.Errorf("unknown command: %s", commandName)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
