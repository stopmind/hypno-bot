package utils

import (
	"strings"
	"text/template"
)

func ExecuteTemplate(template *template.Template, data any) (string, error) {
	builder := &strings.Builder{}
	err := template.Execute(builder, data)
	return builder.String(), err
}
