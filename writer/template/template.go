package wrtmpl

import (
	"text/template"

	"github.com/hashicorp/go-multierror"
)

type WrTemplateType string
type tmplMap map[WrTemplateType]*template.Template

const (
	WrTmplGoXXConfMap         = WrTemplateType("GoXXConfMap")
	WrTmplGoFileCommentSource = WrTemplateType("GoFileCommentSource")
	WrTmplProtoFileComment    = WrTemplateType("ProtoFileComment")
	WrTmplProtoDecl           = WrTemplateType("ProtoDecl")
	WrTmplJsonFields          = WrTemplateType("JsonFields")
)

var (
	wrTemplates = make(tmplMap)
)

func GetWrTemplate(tmplType WrTemplateType) *template.Template {
	return wrTemplates[tmplType]
}

func init() {
	var errRes *multierror.Error

	// Go templates
	if tmpl, err := template.New(string(WrTmplGoXXConfMap)).Parse(tmplGoXXConfMapText); err == nil {
		wrTemplates[WrTmplGoXXConfMap] = tmpl
	} else {
		errRes = multierror.Append(errRes, err)
	}
	if tmpl, err := template.New(string(WrTmplGoFileCommentSource)).Parse(tmplGoFileCommentSourceText); err == nil {
		wrTemplates[WrTmplGoFileCommentSource] = tmpl
	} else {
		errRes = multierror.Append(errRes, err)
	}

	// Protobuf templates
	if tmpl, err := template.New(string(WrTmplProtoFileComment)).Parse(tmplProtoFileCommentText); err == nil {
		wrTemplates[WrTmplProtoFileComment] = tmpl
	} else {
		errRes = multierror.Append(errRes, err)
	}
	if tmpl, err := template.New(string(WrTmplProtoDecl)).Parse(tmplProtoDeclText); err == nil {
		wrTemplates[WrTmplProtoDecl] = tmpl
	} else {
		errRes = multierror.Append(errRes, err)
	}

	// Json templates
	if tmpl, err := template.New(string(WrTmplJsonFields)).Parse(tmplJsonFieldsText); err == nil {
		wrTemplates[WrTmplJsonFields] = tmpl
	} else {
		errRes = multierror.Append(errRes, err)
	}

	if errRes.ErrorOrNil() != nil {
		panic(errRes.ErrorOrNil())
	}
}
