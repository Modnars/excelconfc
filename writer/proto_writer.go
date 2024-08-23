package writer

import (
	"fmt"
	"io"
	"path"
	"strings"

	"git.woa.com/modnarshen/excelconfc/code/template"
	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/util"
)

const (
	outProtoFileSuffix = ".ec.proto"
)

func writeProtoFileComment(wr io.Writer, filePath string, sheetName string) error {
	tmplParams := template.T{
		"SourceFile":  path.Base(filePath),
		"SourceSheet": sheetName,
	}
	if err := template.ExecuteTemplate(wr, template.TmplProtoComments, tmplParams); err != nil {
		return fmt.Errorf("exectue template failed|tmplName:%s", template.TmplProtoComments)
	}
	return nil
}

func writeProtoDecl(wr io.Writer, goPackage string) error {
	tmplParams := template.T{
		"PackageName": getPackageName(goPackage),
		"GoPackage":   goPackage,
	}
	if err := template.ExecuteTemplate(wr, template.TmplProtoCodePackage, tmplParams); err != nil {
		util.LogError("exectue template failed|tmplName:%s", template.TmplProtoCodePackage)
		return err
	}
	return nil
}

func writeProtoMessage(wr io.Writer, headers [][]string, sheetName string) error {
	indent := 0
	msgLabelNumber := 0
	wrf(wr, "\nmessage %s {\n", sheetName)
	indent += 1
	for i, name := range headers[rules.ROW_IDX_NAME] {
		msgLabelNumber += 1
		wrf(wr, "%s%s %s = %d;\n", indentSpace(indent), headers[rules.ROW_IDX_TYPE][i], name, msgLabelNumber)
	}
	indent -= 1
	wrf(wr, "}\n")
	return nil
}

func WriteToProtoFile(headers [][]string, filePath string, sheetName string, goPackage string, outDir string) error {
	var wr strings.Builder

	if err := writeProtoFileComment(&wr, filePath, sheetName); err != nil {
		util.LogError("generate proto file comment failed|fileName:%s|sheetName:%s|headers:{%+v}|err:{%+v}", path.Base(filePath), sheetName, headers, err)
		return err
	}
	if err := writeProtoDecl(&wr, goPackage); err != nil {
		util.LogError("generate proto declaration failed|fileName:%s|sheetName:%s|headers:{%+v}|err:{%+v}", path.Base(filePath), sheetName, headers, err)
		return err
	}
	if err := writeProtoMessage(&wr, headers, sheetName); err != nil {
		util.LogError("generate proto message failed|fileName:%s|sheetName:%s|headers:{%+v}|err:{%+v}", path.Base(filePath), sheetName, headers, err)
		return err
	}

	return WriteToFile(outDir, sheetName, outProtoFileSuffix, []byte(wr.String()))
}
