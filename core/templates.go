package core

import (
	"html/template"
	"strings"
)

var templateFuncs = template.FuncMap{
	"replace":    strings.Replace,
	"memberName": MemberName,
	"sum":        func(a int, b int) int { return a + b },
}
