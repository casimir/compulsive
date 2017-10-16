package compulsive

import (
	"bytes"
	"strings"
	"text/template"
)

const pkgDescTpl = `Package: {{.Provider.Name}}/{{.Name}}
Name: {{.Label}}{{if .Summary}}
Summary: {{.Summary}}{{end}}{{if .Binaries}}
Binaries: {{StringsJoin .Binaries ", "}}{{end}}
Version: {{.Version}}{{if .NextVersion}}
Available: {{.NextVersion}}{{end}}
`

func FmtPkgDesc(pkg Package) string {
	t := template.New("description")
	t.Funcs(template.FuncMap{"StringsJoin": strings.Join})
	t = template.Must(t.Parse(pkgDescTpl))
	buf := bytes.NewBufferString("")
	if err := t.Execute(buf, pkg); err != nil {
		panic(err)
	}
	return buf.String()
}

func FmtPkgLine(pkg Package) string {
	var line []string
	line = append(line, pkg.Provider.Name()+"/"+pkg.Name)
	if pkg.Label != pkg.Name {
		line = append(line, "-")
		line = append(line, pkg.Label)
	}
	if pkg.State == StateOutdated {
		line = append(line, "("+pkg.Version+" → "+pkg.NextVersion+")")
	} else {
		line = append(line, "("+pkg.Version+")")
	}
	return strings.Join(line, " ")
}
