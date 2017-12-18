package lib

import (
	"bytes"
	"os"
	"text/template"

	"github.com/pkg/errors"
)

var (
	FuncMap = template.FuncMap{
		"env": os.Getenv,
	}
)

// TemplateField takes a field string and a state.
// With that it applies the state in the template and
// generates a response.
func TemplateField(field string, state *RenderState) (res string, err error) {
	var (
		tmpl   *template.Template
		output bytes.Buffer
	)

	tmpl, err = template.
		New("tmpl").
		Funcs(FuncMap).
		Parse(field)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to instantiate template for record '%s'",
			field)
		return
	}

	err = tmpl.Execute(&output, state)
	if err != nil {
		err = errors.Wrapf(err, "failed to execute template")
		return
	}

	res = output.String()

	return
}
