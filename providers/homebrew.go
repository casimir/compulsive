package providers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/casimir/compulsive"
)

type brewPkgInfo struct {
	provider compulsive.Provider
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Outdated bool   `json:"outdated"`
	Versions struct {
		Stable string `json:"stable"`
	} `json:"versions"`
	Installed []struct {
		Version string `json:"version"`
	} `json:"installed"`
}

type Brew struct{}

func (p *Brew) Name() string {
	return "brew"
}

func (p *Brew) IsAvailable() bool {
	out, err := exec.Command("brew", "--version").Output()
	if err != nil {
		return false
	}
	matched, err := regexp.Match(`\bHomebrew \d+.\d+.\d+-\d+-\w+\b`, out)
	return err == nil && matched

}

func (p *Brew) Sync() error {
	return exec.Command("brew", "update").Run()
}

func (p *Brew) List() ([]compulsive.Package, error) {
	out, err := exec.Command("brew", "info", "--json=v1", "--installed").Output()
	if err != nil {
		return nil, fmt.Errorf("error while fetching packages: %s", err)
	}
	var pkgsInfo []brewPkgInfo
	if err := json.Unmarshal(out, &pkgsInfo); err != nil {
		return nil, fmt.Errorf("failed to decode package info: %s", err)
	}
	var pkgs []compulsive.Package
	for _, it := range pkgsInfo {
		var versions []string
		for _, installed := range it.Installed {
			versions = append(versions, installed.Version)
		}
		pkg := compulsive.Package{
			Provider:    p,
			Name:        it.FullName,
			Label:       it.Name,
			State:       compulsive.StateUpToDate,
			Version:     strings.Join(versions, "/"),
			NextVersion: it.Versions.Stable,
		}
		if it.Outdated {
			pkg.State = compulsive.StateOutdated
		}
		pkgs = append(pkgs, pkg)
	}
	return pkgs, nil
}

func (p *Brew) UpdateCommand(pkgs ...compulsive.Package) string {
	var names []string
	for _, it := range pkgs {
		names = append(names, it.Name)
	}
	return "brew upgrade " + strings.Join(names, " ")
}

func NewBrew() compulsive.Provider {
	return &Brew{}
}
