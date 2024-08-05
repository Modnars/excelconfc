package wrtmpl

import (
	"text/template"
)

type wrTemplateType string
type tmplMap map[wrTemplateType]*template.Template

const (
	WrTmplGoXXConfMap         = wrTemplateType("GoXXConfMap")
	WrTmplGoFileCommentSource = wrTemplateType("GoFileCommentSource")
	WrTmplProtoFileComment    = wrTemplateType("ProtoFileComment")
	WrTmplProtoDecl           = wrTemplateType("ProtoDecl")
	WrTmplJsonFields          = wrTemplateType("JsonFields")
)

var (
	wrTemplates = make(tmplMap)
)

func GetWrTemplate(tmplType wrTemplateType) *template.Template {
	return wrTemplates[tmplType]
}

func init() {
	// Go templates
	wrTemplates[WrTmplGoXXConfMap] = template.Must(template.New(string(WrTmplGoXXConfMap)).Parse(tmplGoXXConfMapText))
	wrTemplates[WrTmplGoFileCommentSource] = template.Must(template.New(string(WrTmplGoFileCommentSource)).Parse(tmplGoFileCommentSourceText))

	// Protobuf templates
	wrTemplates[WrTmplProtoFileComment] = template.Must(template.New(string(WrTmplProtoFileComment)).Parse(tmplProtoFileCommentText))
	wrTemplates[WrTmplProtoDecl] = template.Must(template.New(string(WrTmplProtoDecl)).Parse(tmplProtoDeclText))

	// Json templates
	wrTemplates[WrTmplJsonFields] = template.Must(template.New(string(WrTmplJsonFields)).Parse(tmplJsonFieldsText))
}
