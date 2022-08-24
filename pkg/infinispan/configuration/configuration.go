package configuration

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed infinispan-14.xml
var ispnTpl string

type Spec struct {
}

// Generate the Infinispan server configuration
func Generate(spec *Spec) (string, error) {
	tpl, err := template.New("config").Parse(ispnTpl)
	if err != nil {
		return "", err
	}

	buffIspn := new(bytes.Buffer)
	if err = tpl.Execute(buffIspn, spec); err != nil {
		return "", err
	}
	return buffIspn.String(), nil
}
