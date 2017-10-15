package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/casimir/compulsive"
)

type goPkgInfo struct {
	ImportPath  string
	Name        string
	Target      string
	Stale       bool
	StaleReason string
}

func loadPackages(p *Go) ([]goPkgInfo, error) {
	out, err := exec.Command("go", "list", "-json", "all").Output()
	if err != nil {
		return nil, fmt.Errorf("error while fetching packages: %s", err)
	}
	var pkgs []goPkgInfo
	dec := json.NewDecoder(bytes.NewReader(out))
	for dec.More() {
		var pkg goPkgInfo
		if err := dec.Decode(&pkg); err != nil {
			return nil, fmt.Errorf("failed to decode package info: %s", err)
		}
		pkgs = append(pkgs, pkg)

	}
	return pkgs, nil
}

func loadMainPackages(p *Go) ([]goPkgInfo, error) {
	allPkgs, err := loadPackages(p)
	if err != nil {
		return nil, err
	}
	var pkgs []goPkgInfo
	for _, it := range allPkgs {
		if it.Name == "main" {
			pkgs = append(pkgs, it)
		}
	}
	return pkgs, nil
}

type binary struct {
	provider compulsive.Provider
	command  string
	modTime  time.Time
	info     goPkgInfo
}

func (b binary) name() string {
	if b.info.ImportPath != "" {
		return b.info.ImportPath
	}
	return b.command
}

type Go struct {
	path string
}

func (p *Go) Name() string {
	return "go"
}

func (p *Go) IsAvailable() bool {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return false
	}
	matched, err := regexp.Match(`\bgo1.\d+(.\d+)?\b`, out)
	return err == nil && matched

}

func (p *Go) Sync() error {
	return nil
}

func (p *Go) List() ([]compulsive.Package, error) {
	binPath := filepath.Join(p.path, "bin")
	binaries, _ := ioutil.ReadDir(binPath)
	pkgsInfo, err := loadMainPackages(p)
	if err != nil {
		return nil, err
	}
	var pkgs []compulsive.Package
	for _, it := range binaries {
		filename := it.Name()
		bin := binary{
			provider: p,
			command:  filename,
			modTime:  it.ModTime(),
		}
		if runtime.GOOS == "windows" {
			bin.command = bin.command[:len(bin.command)-4]
		}
		target := filepath.Join(binPath, filename)
		for _, pi := range pkgsInfo {
			if pi.Target == target {
				bin.info = pi
				break
			}
		}
		pkg := compulsive.Package{
			Provider:    p,
			Name:        bin.name(),
			Label:       bin.command,
			State:       compulsive.StateUpToDate,
			Version:     bin.modTime.Format("2006-01-02"),
			NextVersion: time.Now().Format("2006-01-02"),
		}
		if bin.info.ImportPath == "" {
			pkg.State = compulsive.StateUnknown
		} else if bin.info.Stale {
			pkg.State = compulsive.StateOutdated
		}
		pkgs = append(pkgs, pkg)
	}
	return pkgs, nil
}

func (p *Go) UpdateCommand(pkgs ...compulsive.Package) string {
	var commands []string
	for _, it := range pkgs {
		commands = append(commands, "go get "+it.Name)
	}
	return strings.Join(commands, "\n")
}

func NewGo() compulsive.Provider {
	return &Go{path: os.Getenv("GOPATH")}
}
