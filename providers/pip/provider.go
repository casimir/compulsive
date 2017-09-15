package pip

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/casimir/compulsive"
)

func pipBin(version string) string {
	return "pip" + version
}

func AvailableV(version string) func() bool {
	if version == "" {
		return func() bool {
			return os.Getenv("VIRTUAL_ENV") != ""
		}
	}
	return func() bool {
		out, err := exec.Command(pipBin(version), "--version").Output()
		if err != nil {
			return false
		}
		matched, err := regexp.Match(`\bpip \d+.\d+.\d+\s`, out)
		return err == nil && matched
	}
}

type pkgInfo struct {
	Name_         string `json:"name"`
	Version_      string `json:"version"`
	LatestVersion string `json:"latest_version"`
}

func (p pkgInfo) Name() string {
	return p.Name_
}

func (p pkgInfo) Label() string {
	return p.Name_
}

func (p pkgInfo) State() compulsive.PackageState {
	if p.LatestVersion != "" {
		return compulsive.StateOutdated
	}
	return compulsive.StateUpToDate
}

func (p pkgInfo) Version() string {
	return p.Version_
}

func (p pkgInfo) NextVersion() string {
	return p.LatestVersion
}

func loadPackages(bin string) []pkgInfo {
	if err := exec.Command(bin, "install", "--upgrade", "pip").Run(); err != nil {
		log.Printf("error while syncing repository: %s\n", err)
	}
	outOutdated, err := exec.Command(bin, "list", "--format", "json", "--outdated").Output()
	if err != nil {
		log.Printf("error while fetching packages: %s\n", err)
	}
	var pkgsOutdated []pkgInfo
	if err := json.Unmarshal(outOutdated, &pkgsOutdated); err != nil {
		log.Printf("failed to decode package info: %s", err)
	}
	outdatedMap := make(map[string]pkgInfo, len(pkgsOutdated))
	for _, it := range pkgsOutdated {
		outdatedMap[it.Name()] = it
	}
	outAll, err := exec.Command(bin, "list", "--format", "json").Output()
	if err != nil {
		log.Printf("error while fetching packages: %s\n", err)
	}
	var pkgsAll []pkgInfo
	if err := json.Unmarshal(outAll, &pkgsAll); err != nil {
		log.Printf("failed to decode package info: %s", err)
	}
	var pkgs []pkgInfo
	for _, it := range pkgsAll {
		pkg := pkgInfo{
			Name_:    it.Name_,
			Version_: it.Version_,
		}
		if outdatedPkg, ok := outdatedMap[pkg.Name_]; ok {
			pkg.LatestVersion = outdatedPkg.LatestVersion
		}
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}

type Pip struct {
	bin      string
	packages []pkgInfo
}

func (p Pip) List() []compulsive.Package {
	var pkgs []compulsive.Package
	for _, it := range p.packages {
		pkgs = append(pkgs, it)
	}
	return pkgs
}

func (p Pip) UpdateCommand(pkgs ...compulsive.Package) string {
	var names []string
	for _, it := range p.packages {
		names = append(names, it.Name())
	}
	return p.bin + " install --upgrade " + strings.Join(names, " ")
}

func NewV(version string) func() compulsive.Provider {
	bin := pipBin(version)
	return func() compulsive.Provider {
		return Pip{
			bin:      bin,
			packages: loadPackages(bin),
		}
	}
}
