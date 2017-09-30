package homebrew

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/casimir/compulsive"
)

type pkgInfo struct {
	provider compulsive.Provider
	Name_    string `json:"name"`
	FullName string `json:"full_name"`
	Outdated bool   `json:"outdated"`
	Versions struct {
		Stable string `json:"stable"`
	} `json:"versions"`
	Installed []struct {
		Version string `json:"version"`
	} `json:"installed"`
}

func (p pkgInfo) Provider() compulsive.Provider {
	return p.provider
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
	return p.Versions.Stable
}

type Homebrew struct{}

func (p *Homebrew) Name() string {
	return "homebrew"
}

func (p *Homebrew) IsAvailable() bool {
	out, err := exec.Command("brew", "--version").Output()
	if err != nil {
		return false
	}
	matched, err := regexp.Match(`\bHomebrew \d+.\d+.\d+-\d+-\w+\b`, out)
	return err == nil && matched

}

func (p *Homebrew) Sync() error {
	return exec.Command("brew", "update").Run()
}

func (p *Homebrew) List() ([]compulsive.Package, error) {
	out, err := exec.Command("brew", "info", "--json=v1", "--installed").Output()
	if err != nil {
		return nil, fmt.Errorf("error while fetching packages: %s\n", err)
	}
	var pkgsInfo []pkgInfo
	if err := json.Unmarshal(out, &pkgsInfo); err != nil {
		return nil, fmt.Errorf("failed to decode package info: %s", err)
	}
	var pkgs []compulsive.Package
	for _, pkg := range pkgsInfo {
		pkg.provider = p
		pkgs = append(pkgs, pkg)
	}
	return pkgs, nil
}

func (p *Homebrew) UpdateCommand(pkgs ...compulsive.Package) string {
	var names []string
	for _, it := range pkgs {
		names = append(names, it.Name())
	}
	return "brew upgrade " + strings.Join(names, " ")
}

func New() compulsive.Provider {
	return &Homebrew{}
}
