package pip

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/casimir/compulsive"
)

var (
	defaultChecked    = false
	defaultPythonRoot = ""
)

var pipRe = regexp.MustCompile(`pip (?P<version>\d+.\d+.\d+) from (?P<root>.+) \((?P<pyversion>.+)\)`)

func checkVersion(version string) (bool, string) {
	out, err := exec.Command("pip"+version, "--version").Output()
	if err != nil {
		return false, ""
	}
	matches := pipRe.FindSubmatch(out)
	return true, string(matches[2])
}

type pkgInfo struct {
	provider      compulsive.Provider
	Name          string `json:"name"`
	Version       string `json:"version"`
	LatestVersion string `json:"latest_version"`
}

type Pip struct {
	version  string
	bin      string
	packages []pkgInfo
}

func (p *Pip) Name() string {
	return p.bin
}

func (p *Pip) IsAvailable() bool {
	if !defaultChecked {
		defaultChecked, defaultPythonRoot = checkVersion("")
	}
	if p.version == "" {
		return defaultChecked && defaultPythonRoot != ""
	}
	available, pythonRoot := checkVersion(p.version)
	if defaultChecked {
		return available && pythonRoot != defaultPythonRoot
	}
	return available

}

func (p *Pip) Sync() error {
	return exec.Command(p.bin, "install", "--upgrade", "pip").Run()
}

func (p *Pip) List() ([]compulsive.Package, error) {
	outOutdated, err := exec.Command(p.bin, "list", "--format", "json", "--outdated").Output()
	if err != nil {
		return nil, fmt.Errorf("error while fetching packages: %s\n", err)
	}
	var pkgsOutdated []pkgInfo
	if err := json.Unmarshal(outOutdated, &pkgsOutdated); err != nil {
		return nil, fmt.Errorf("failed to decode package info: %s", err)
	}
	outdatedMap := make(map[string]pkgInfo, len(pkgsOutdated))
	for _, it := range pkgsOutdated {
		outdatedMap[it.Name] = it
	}
	outAll, err := exec.Command(p.bin, "list", "--format", "json").Output()
	if err != nil {
		return nil, fmt.Errorf("error while fetching packages: %s\n", err)
	}
	var pkgsAll []pkgInfo
	if err := json.Unmarshal(outAll, &pkgsAll); err != nil {
		return nil, fmt.Errorf("failed to decode package info: %s", err)
	}
	var pkgs []compulsive.Package
	for _, it := range pkgsAll {
		pkg := compulsive.Package{
			Provider: p,
			Name:     it.Name,
			Label:    it.Name,
			Version:  it.Version,
			State:    compulsive.StateUpToDate,
		}
		if outdatedPkg, ok := outdatedMap[it.Name]; ok {
			pkg.NextVersion = outdatedPkg.LatestVersion
			pkg.State = compulsive.StateOutdated
		}
		pkgs = append(pkgs, pkg)
	}
	return pkgs, nil
}

func (p *Pip) UpdateCommand(pkgs ...compulsive.Package) string {
	var names []string
	for _, it := range pkgs {
		names = append(names, it.Name)
	}
	return p.bin + " install --upgrade " + strings.Join(names, " ")
}

func New(version string) compulsive.Provider {
	return &Pip{version: version, bin: "pip" + version}
}
