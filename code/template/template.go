package template

import (
	"io"
	"text/template"
)

type codeTmplType string

const (
	templatesDir = "code/templates/"

	TmplGoCommentsHead    = codeTmplType("golang.comments.head.tmpl")
	TmplGoCommentsSource  = codeTmplType("golang.comments.source.tmpl")
	TmplGoCodePackage     = codeTmplType("golang.code.package.tmpl")
	TmplGoCodeImport      = codeTmplType("golang.code.import.tmpl")
	TmplGoCodeDefDateTime = codeTmplType("golang.code.def.datetime.tmpl")
	TmplGoCodeConfMap     = codeTmplType("golang.code.confmap.tmpl")
)

var (
	globalTemplate *template.Template
)

func init() {
	var err error
	if globalTemplate, err = template.ParseFiles(
		string(templatesDir+TmplGoCommentsHead),
		string(templatesDir+TmplGoCommentsSource),
		string(templatesDir+TmplGoCodePackage),
		string(templatesDir+TmplGoCodeImport),
		string(templatesDir+TmplGoCodeDefDateTime),
		string(templatesDir+TmplGoCodeConfMap),
	); err != nil {
		panic(err)
	}
}

func ExecuteTemplate(writer io.Writer, codeTmpl codeTmplType, tmplParams any) error {
	return globalTemplate.ExecuteTemplate(writer, string(codeTmpl), tmplParams)
}
