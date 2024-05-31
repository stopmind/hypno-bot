package core

import (
	"html/template"
	"strings"
)

var templateFuncs = template.FuncMap{
	"replace": strings.Replace,
}
