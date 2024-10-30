package template

import (
	"embed"
	"io"
	"text/template"
)

type codeTmplType string
type T map[string]any

const (
	TmplGoCommentsHead    = codeTmplType("golang.comments.head.tmpl")
	TmplGoCommentsSource  = codeTmplType("golang.comments.source.tmpl")
	TmplGoCodePackage     = codeTmplType("golang.code.package.tmpl")
	TmplGoCodeImport      = codeTmplType("golang.code.import.tmpl")
	TmplGoCodeDefDateTime = codeTmplType("golang.code.def.datetime.tmpl")
	TmplGoCodeConfMap     = codeTmplType("golang.code.confmap.tmpl")

	TmplProtoCodePackage = codeTmplType("proto.code.package.tmpl")
	TmplProtoComments    = codeTmplType("proto.comments.tmpl")
)

var (
	//go:embed templates/*.tmpl
	embedFS        embed.FS
	globalTemplate *template.Template
)

func init() {
	var err error
	if globalTemplate, err = template.ParseFS(embedFS, "templates/*.tmpl"); err != nil {
		panic(err)
	}
}

func ExecuteTemplate(writer io.Writer, codeTmpl codeTmplType, tmplParams any) error {
	return globalTemplate.ExecuteTemplate(writer, string(codeTmpl), tmplParams)
}
