package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/casimir/compulsive"
)

var (
	cargoRe      = regexp.MustCompile(`^cargo (?P<version>\d+\.\d+\.\d+)`)
	cargoEntryRe = regexp.MustCompile(`"(?P<name>\S+) (?P<version>\S+) \((?P<uri>\S+)\)" = \[(?P<binaries>[^]]+)\]`)
)

type (
	cargoManifestEntry struct {
		name     string
		version  string
		uri      string
		binaries []string
	}

	cargoPkgPayload struct {
		Crate struct {
			Description string
			MaxVersion  string `json:"max_version"`
		}
	}
)

func cleanManifestBinaries(raw []byte) []string {
	var names []string
	for _, it := range bytes.Split(raw, []byte(", ")) {
		names = append(names, string(bytes.Trim(it, "\"")))
	}
	return names
}

func unmarshalManifest(raw []byte) []cargoManifestEntry {
	var manifest []cargoManifestEntry
	lookingForSection := true
	for _, line := range bytes.Split(raw, []byte("\n")) {
		if lookingForSection {
			if bytes.Equal(line, []byte("[v1]")) {
				lookingForSection = false
			}
			continue
		}
		matches := cargoEntryRe.FindSubmatch(line)
		if len(matches) == 0 {
			continue
		}
		entry := cargoManifestEntry{
			name:     string(matches[1]),
			version:  string(matches[2]),
			uri:      string(matches[3]),
			binaries: cleanManifestBinaries(matches[4]),
		}
		manifest = append(manifest, entry)
	}
	return manifest
}

func fetchPkgInfo(uri string, pkg *compulsive.Package) error {
	if uri != "registry+https://github.com/rust-lang/crates.io-index" {
		return nil
	}
	resp, err := http.Get("https://crates.io/api/v1/crates/" + pkg.Name)
	if err != nil {
		return err
	}
	var payload cargoPkgPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}
	pkg.Summary = payload.Crate.Description
	pkg.NextVersion = payload.Crate.MaxVersion
	if pkg.Version == pkg.NextVersion {
		pkg.State = compulsive.StateUpToDate
	} else {
		pkg.State = compulsive.StateOutdated
	}
	return nil
}

type Cargo struct {
	manifest []cargoManifestEntry
}

func (p *Cargo) Name() string {
	return "cargo"
}

func (p *Cargo) IsAvailable() bool {
	out, err := exec.Command("cargo", "version").Output()
	if err != nil {
		return false
	}
	return cargoRe.Match(out)
}

func (p *Cargo) Sync() error {
	return nil
}

func (p *Cargo) loadManifest() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	manifestPath := filepath.Join(usr.HomeDir, ".cargo", ".crates.toml")
	raw, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return err
	}
	p.manifest = unmarshalManifest(raw)
	return nil
}

func (p *Cargo) List() ([]compulsive.Package, error) {
	if err := p.loadManifest(); err != nil {
		return nil, fmt.Errorf("could not build package list: %s", err)
	}
	var pkgs []compulsive.Package
	for _, it := range p.manifest {
		pkg := compulsive.Package{
			Provider: p,
			Name:     it.name,
			Label:    it.name,
			Binaries: it.binaries,
			State:    compulsive.StateUnknown,
			Version:  it.version,
		}
		if err := fetchPkgInfo(it.uri, &pkg); err != nil {
			log.Printf("failed to fetch data for package %q: %s", it.name, err)
		}
		pkgs = append(pkgs, pkg)
	}
	return pkgs, nil
}

func (p *Cargo) UpdateCommand(pkgs ...compulsive.Package) string {
	var names []string
	for _, it := range pkgs {
		names = append(names, it.Name)
	}
	return "cargo install --force " + strings.Join(names, " ")
}

func NewCargo() compulsive.Provider {
	return &Cargo{}
}
