package writer

import (
	"fmt"
	"io"
	"path"
	"strings"

	"git.woa.com/modnarshen/excelconfc/code/template"
	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/types"
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
		return fmt.Errorf("exectue template failed|tmplName:%s -> %w", template.TmplProtoComments, err)
	}
	return nil
}

func writeProtoDecl(wr io.Writer, goPackage string) error {
	tmplParams := template.T{
		"PackageName": getPackageName(goPackage),
		"GoPackage":   goPackage,
	}
	if err := template.ExecuteTemplate(wr, template.TmplProtoCodePackage, tmplParams); err != nil {
		return fmt.Errorf("exectue template failed|tmplName:%s -> %w", template.TmplProtoCodePackage, err)
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

func writeProtoEnum(wr io.Writer, enumTypes []*types.EnumTypeSt) error {
	indent := 0
	for _, enumType := range enumTypes {
		wrf(wr, "\nenum %s {\n", enumType.Name)
		indent++
		for _, enumVal := range enumType.EnumVals {
			wrf(wr, "%s%s = %s;\n", indentSpace(indent), enumVal.Name, enumVal.ID)
		}
		indent--
		wrf(wr, "}\n")
	}
	return nil
}

func WriteToProtoFile(outData types.DataHolder, goPackage string, outDir string) error {
	var wr strings.Builder

	if err := writeProtoFileComment(&wr, outData.GetFileName(), outData.GetSheetName()); err != nil {
		return fmt.Errorf("generate proto file comment failed|file:%s|sheet:%s -> %w", outData.GetFileName(), outData.GetSheetName(), err)
	}
	if err := writeProtoDecl(&wr, goPackage); err != nil {
		return fmt.Errorf("generate proto declaration failed|file:%s|sheet:%s -> %w", outData.GetFileName(), outData.GetSheetName(), err)
	}
	if err := writeProtoEnum(&wr, outData.GetEnumTypes()); err != nil {
		return fmt.Errorf("generate proto message failed|file:%s|sheet:%s|enumTypes:{%+v} -> %w", outData.GetFileName(), outData.GetSheetName(), outData.GetEnumTypes(), err)
	}
	if err := writeProtoMessage(&wr, outData.GetHeaders(), outData.GetSheetName()); err != nil {
		return fmt.Errorf("generate proto message failed|file:%s|sheet:%s|headers:{%+v} -> %w", outData.GetFileName(), outData.GetSheetName(), outData.GetHeaders(), err)
	}

	return WriteToFile(outDir, outData.GetSheetName(), outProtoFileSuffix, []byte(wr.String()))
}
