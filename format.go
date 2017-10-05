package compulsive

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

const pkgDescTpl = `Package: {{.Provider.Name}}/{{.Name}}
Name: {{.Label}}
Version: {{.Version}}{{if .NextVersion}}
Available: {{.NextVersion}}{{end}}
`

func FmtPkgDesc(pkg Package) string {
	t := template.Must(template.New("description").Parse(pkgDescTpl))
	buf := bytes.NewBufferString("")
	if err := t.Execute(buf, pkg); err != nil {
		panic(err)
	}
	return buf.String()
}

func FmtPkgLine(pkg Package) string {
	var line []string
	if pkg.Label() == pkg.Name() {
		line = append(line, pkg.Label())
	} else {
		label := fmt.Sprintf("%s - %s", pkg.Label(), pkg.Name())
		line = append(line, label)
	}
	if pkg.State() == StateOutdated {
		line = append(line, "("+pkg.Version()+" → "+pkg.NextVersion()+")")
	} else {
		line = append(line, "("+pkg.Version()+")")
	}
	return strings.Join(line, " ")
}
