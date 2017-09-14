package golang

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/casimir/compulsive"
)

func Available() bool {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return false
	}
	matched, err := regexp.Match(`\bgo1.\d+(.\d+)?\b`, out)
	return err == nil && matched
}

type pkgInfo struct {
	ImportPath  string
	Name        string
	Target      string
	Stale       bool
	StaleReason string
}

func loadPackages() []pkgInfo {
	out, err := exec.Command("go", "list", "-json", "all").Output()
	if err != nil {
		log.Printf("error while fetching packages: %s\n", err)
	}
	var pkgs []pkgInfo
	dec := json.NewDecoder(bytes.NewReader(out))
	for dec.More() {
		var pkg pkgInfo
		if err := dec.Decode(&pkg); err != nil {
			log.Printf("failed to decode package info: %s", err)
		} else {
			pkgs = append(pkgs, pkg)
		}
	}
	return pkgs
}

func loadMainPackages() []pkgInfo {
	var pkgs []pkgInfo
	for _, it := range loadPackages() {
		if it.Name == "main" {
			pkgs = append(pkgs, it)
		}
	}
	return pkgs
}

type binary struct {
	command string
	modTime time.Time
	info    pkgInfo
}

func (b binary) Name() string {
	return b.info.ImportPath
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
	path     string
	binaries []binary
}

func (m Golang) List() []compulsive.Package {
	var pkgs []compulsive.Package
	for _, it := range m.binaries {
		pkgs = append(pkgs, it)
	}
	return pkgs
}

func (m Golang) UpdateCommand(pkgs ...compulsive.Package) string {
	var commands []string
	for _, it := range pkgs {
		commands = append(commands, "go get "+it.Name())
	}
	return strings.Join(commands, "\n") + "\n"
}

func New() compulsive.Provider {
	manager := Golang{
		path: os.Getenv("GOPATH"),
	}
	binPath := filepath.Join(manager.path, "bin")
	binaries, _ := ioutil.ReadDir(binPath)
	pkgsInfo := loadMainPackages()
	for _, it := range binaries {
		filename := it.Name()
		bin := binary{
			command: filename,
			modTime: it.ModTime(),
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
		manager.binaries = append(manager.binaries, bin)
	}
	return manager
}
