package golang

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

type pkgInfo struct {
	ImportPath  string
	Name        string
	Target      string
	Stale       bool
	StaleReason string
}

func loadPackages(p *Golang) ([]pkgInfo, error) {
	out, err := exec.Command("go", "list", "-json", "all").Output()
	if err != nil {
		return nil, fmt.Errorf("error while fetching packages: %s\n", err)
	}
	var pkgs []pkgInfo
	dec := json.NewDecoder(bytes.NewReader(out))
	for dec.More() {
		var pkg pkgInfo
		if err := dec.Decode(&pkg); err != nil {
			return nil, fmt.Errorf("failed to decode package info: %s", err)
		}
		pkgs = append(pkgs, pkg)

	}
	return pkgs, nil
}

func loadMainPackages(p *Golang) ([]pkgInfo, error) {
	allPkgs, err := loadPackages(p)
	if err != nil {
		return nil, err
	}
	var pkgs []pkgInfo
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
	info     pkgInfo
}

func (b binary) Provider() compulsive.Provider {
	return b.provider
}

func (b binary) Name() string {
	if b.info.ImportPath != "" {
		return b.info.ImportPath
	}
	return b.command
}

func (b binary) Label() string {
	return b.command
}

func (b binary) State() compulsive.PackageState {
	if b.info.ImportPath == "" {
		return compulsive.StateUnknown
	} else if b.info.Stale {
		return compulsive.StateOutdated
	} else {
		return compulsive.StateUpToDate
	}
}

func (b binary) Version() string {
	return b.modTime.Format("2006-01-02")
}

func (b binary) NextVersion() string {
	return time.Now().Format("2006-01-02")
}

type Golang struct {
	path string
}

func (p *Golang) Name() string {
	return "go"
}

func (p *Golang) IsAvailable() bool {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return false
	}
	matched, err := regexp.Match(`\bgo1.\d+(.\d+)?\b`, out)
	return err == nil && matched

}

func (p *Golang) Sync() error {
	return nil
}

func (p *Golang) List() ([]compulsive.Package, error) {
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
		pkgs = append(pkgs, bin)
	}
	return pkgs, nil
}

func (p *Golang) UpdateCommand(pkgs ...compulsive.Package) string {
	var commands []string
	for _, it := range pkgs {
		commands = append(commands, "go get "+it.Name())
	}
	return strings.Join(commands, "\n")
}

func New() compulsive.Provider {
	return &Golang{path: os.Getenv("GOPATH")}
}
