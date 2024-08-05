package writer

import (
	"io"
	"os"
	"path"
	"strings"

	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/util"
	wrtmpl "git.woa.com/modnarshen/excelconfc/writer/template"
)

const (
	outProtoFileSuffix = ".ec.proto"
)

func writeProtoFileComment(wr io.Writer, filePath string, sheetName string) error {
	tmplParams := struct {
		SourceFile  string
		SourceSheet string
	}{
		SourceFile:  path.Base(filePath),
		SourceSheet: sheetName,
	}
	if err := wrtmpl.GetWrTemplate(wrtmpl.WrTmplProtoFileComment).Execute(wr, tmplParams); err != nil {
		util.LogError("exectue template failed|tmplName:%s", wrtmpl.WrTmplProtoFileComment)
		return err
	}
	return nil
}

func writeProtoDecl(wr io.Writer, goPackage string) error {
	tmplParams := struct {
		PackageName string
		GoPackage   string
	}{
		PackageName: getPackageName(goPackage),
		GoPackage:   goPackage,
	}
	if err := wrtmpl.GetWrTemplate(wrtmpl.WrTmplProtoDecl).Execute(wr, tmplParams); err != nil {
		util.LogError("exectue template failed|tmplName:%s", wrtmpl.WrTmplProtoDecl)
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

	os.WriteFile(genOutFilePath(outDir, sheetName, outProtoFileSuffix), []byte(wr.String()), outFilePerm)
	return nil
}
