package homebrew

import (
	"encoding/json"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/casimir/compulsive"
)

func Available() bool {
	out, err := exec.Command("brew", "--version").Output()
	if err != nil {
		return false
	}
	matched, err := regexp.Match(`\bHomebrew \d+.\d+.\d+-\d+-\w+\b`, out)
	return err == nil && matched
}

type pkgInfo struct {
	Name_     string `json:"name"`
	FullName  string `json:"full_name"`
	Outdated  bool   `json:"outdated"`
	Installed []struct {
		Version string `json:"version"`
	} `json:"installed"`
}

func (p pkgInfo) Name() string {
	return p.FullName
}

func (p pkgInfo) Label() string {
	return p.Name_
}

func (p pkgInfo) State() compulsive.PackageState {
	if p.Outdated {
		return compulsive.StateOutdated
	}
	return compulsive.StateUpToDate
}

func (p pkgInfo) Version() string {
	var versions []string
	for _, it := range p.Installed {
		versions = append(versions, it.Version)
	}
	return strings.Join(versions, "/")
}

func (p pkgInfo) NextVersion() string {
	return "?"
}

func loadPackages() []pkgInfo {
	if err := exec.Command("brew", "update").Run(); err != nil {
		log.Printf("error while syncing repository: %s\n", err)
	}
	out, err := exec.Command("brew", "info", "--json=v1", "--installed").Output()
	if err != nil {
		log.Printf("error while fetching packages: %s\n", err)
	}
	var pkgs []pkgInfo
	if err := json.Unmarshal(out, &pkgs); err != nil {
		log.Printf("failed to decode package info: %s", err)
	}
	return pkgs
}

type Homebrew struct {
	packages []pkgInfo
}

func (p Homebrew) List() []compulsive.Package {
	var pkgs []compulsive.Package
	for _, it := range p.packages {
		pkgs = append(pkgs, it)
	}
	return pkgs
}

func (p Homebrew) UpdateCommand(pkgs ...compulsive.Package) string {
	var names []string
	for _, it := range p.packages {
		names = append(names, it.Name())
	}
	return "brew upgrade " + strings.Join(names, " ") + "\n"
}

func New() compulsive.Provider {
	return Homebrew{
		packages: loadPackages(),
	}
}
