package wrtmpl

const (
	tmplJsonFieldsText = `
{{- .IndentSpace }}"filepath": "{{.Filepath}}",
{{.IndentSpace}}"basename": "{{.Basename}}",
{{.IndentSpace}}"sheet": "{{.Sheet}}",
{{.IndentSpace}}"outdir": "{{.Outdir}}",
`
)
