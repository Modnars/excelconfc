package writer

import (
	"fmt"
	"go/format"
	"io"
	"strings"

	"git.woa.com/modnarshen/excelconfc/code/template"
	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/types"
	"git.woa.com/modnarshen/excelconfc/util"
)

const (
	// Output file suffix `.ec.go` means Excel Config Go code.
	outGoFileSuffix  = ".ec.go"
	outGoDefFileName = "excelconf.def"
)

func toOutBytes(output string) ([]byte, error) {
	var outBytes []byte
	if rules.DEBUG_MODE {
		outBytes = []byte(output)
	} else {
		var err error
		outBytes, err = format.Source([]byte(output))
		if err != nil {
			return nil, fmt.Errorf("format code failed -> %w", err)
		}
	}
	return outBytes, nil
}

func writeGoDeclPackage(wr io.Writer, goPackage string) {
	wrf(wr, "\npackage %s\n", getPackageName(goPackage))
}

func writeGoDeclImport(wr io.Writer, packages ...string) {
	if len(packages) <= 0 {
		return
	} else if len(packages) == 1 {
		wrf(wr, "\nimport \"%s\"\n", packages[0])
	}
	outLines := []string{}
	outLines = append(outLines, "\nimport (")
	for _, packageName := range packages {
		outLines = append(outLines, indentSpace(1)+"\""+packageName+"\"")
	}
	outLines = append(outLines, ")")
	wrf(wr, strings.Join(outLines, "\n")+"\n")
}

func genGoStructFieldType(tp string, desc string) string {
	if tp == "string" {
		if desc == "D" {
			return "DateTime"
		}
	}
	return tp
}

func genGoConfKeyInfo(headers [][]string) (string, string) {
	keyIndex := 0
	for i := range headers[rules.ROW_IDX_NAME] {
		if strings.Contains(headers[rules.ROW_IDX_DESC][i], "K") {
			keyIndex = i
			break
		}
	}
	return genGoStructFieldType(headers[rules.ROW_IDX_TYPE][keyIndex], headers[rules.ROW_IDX_DESC][keyIndex]),
		util.SnakeToPascal(headers[rules.ROW_IDX_NAME][keyIndex])
}

func writeGoFileComments(wr io.Writer, fileName string, sheetName string) error {
	template.ExecuteTemplate(wr, template.TmplGoCommentsHead, nil)
	// 仅当 fileName 和 sheetName 有意义时才生成文件注释代码
	if fileName != "" && sheetName != "" {
		tmplParams := template.T{
			"SourceFile":  fileName,
			"SourceSheet": sheetName,
		}
		if err := template.ExecuteTemplate(wr, template.TmplGoCommentsSource, tmplParams); err != nil {
			return fmt.Errorf("exectue template failed|tmplName:%s -> %w", template.TmplGoCommentsSource, err)
		}
	} else {
		// 此处显式添加换行，这是因为 Title 本身是不会有换行的，以此来保证注释代码的连贯性。
		wrf(wr, "\n")
	}
	return nil
}

func writeGoDeclaration(wr io.Writer, goPackage string, importPkgs ...string) error {
	writeGoDeclPackage(wr, goPackage)
	writeGoDeclImport(wr, importPkgs...)
	return nil
}

func writeGoEnum(wr io.Writer, enumTypes []*types.EnumTypeSt) error {
	indent := 0
	wrf(wr, "\n")
	for _, enumType := range enumTypes {
		wrf(wr, "type %s int32\n", enumType.Name)
	}

	wrf(wr, "\nconst (")
	for _, enumType := range enumTypes {
		indent++
		wrf(wr, "\n")
		for _, enumVal := range enumType.EnumVals {
			wrf(wr, "%s%s %s = %s\n", indentSpace(indent), enumVal.Name, enumType.Name, enumVal.ID)
		}
		indent--
	}
	wrf(wr, ")\n")

	wrf(wr, "\nvar (")
	for _, enumType := range enumTypes {
		indent++
		wrf(wr, "\n%s%s_name = map[int32]string{\n", indentSpace(indent), enumType.Name)
		indent++
		for _, enumVal := range enumType.EnumVals {
			wrf(wr, "%s%s: \"%s\",\n", indentSpace(indent), enumVal.ID, enumVal.Name)
		}
		indent--
		wrf(wr, "%s}\n", indentSpace(indent))
		wrf(wr, "%s%s_value = map[string]int32{\n", indentSpace(indent), enumType.Name)
		indent++
		for _, enumVal := range enumType.EnumVals {
			wrf(wr, "%s\"%s\": %s,\n", indentSpace(indent), enumVal.Name, enumVal.ID)
		}
		indent--
		wrf(wr, "%s}\n", indentSpace(indent))
		indent--
	}
	wrf(wr, ")\n")

	for _, enumType := range enumTypes {
		wrf(wr, "\nfunc(x %s) String() string {\n", enumType.Name)
		indent++
		wrf(wr, "return %s_name[int32(x)]\n", enumType.Name)
		indent--
		wrf(wr, "}\n")
	}
	return nil
}

func writeGoConfStruct(wr io.Writer, headers [][]string, sheetName string) error {
	indent := 0
	wrf(wr, "\ntype %s struct {\n", sheetName)
	indent += 1
	for i, name := range headers[rules.ROW_IDX_NAME] {
		wrf(wr, "%s%s %s `json:\"%s,omitempty\"`\n",
			indentSpace(indent), util.SnakeToPascal(name),
			genGoStructFieldType(headers[rules.ROW_IDX_TYPE][i], headers[rules.ROW_IDX_DESC][i]),
			name)
	}
	indent -= 1
	wrf(wr, "}\n")
	return nil
}

func writeConfMapStruct(wr io.Writer, headers [][]string, sheetName string) error {
	confKeyType, confKeyField := genGoConfKeyInfo(headers)
	tmplParams := template.T{
		"XXConf":         sheetName,
		"XXConfMap":      sheetName + "Map",
		"XXConfKeyType":  confKeyType,
		"XXConfKeyField": confKeyField,
	}

	if err := template.ExecuteTemplate(wr, template.TmplGoCodeConfMap, tmplParams); err != nil {
		return fmt.Errorf("exectue template failed|tmplName:%s -> %w", template.TmplGoCodeConfMap, err)
	}
	return nil
}

func outputGoDefFile(goPackage string, outDir string) error {
	var wr strings.Builder
	if err := writeGoFileComments(&wr, "", ""); err != nil {
		return fmt.Errorf("write Go file comments failed|fineName:%s -> %w", outGoDefFileName, err)
	}
	if err := writeGoDeclaration(&wr, goPackage, "encoding/json", "time"); err != nil {
		return fmt.Errorf("write Go declaration failed -> %w", err)
	}
	if err := template.ExecuteTemplate(&wr, template.TmplGoCodeDefDateTime, nil); err != nil {
		return fmt.Errorf("exectue template failed|tmplName:%s -> %w", template.TmplGoCodeDefDateTime, err)
	}
	outBytes, err := toOutBytes(wr.String())
	if err != nil {
		return fmt.Errorf("to output bytes failed -> %w", err)
	}
	return WriteToFile(outDir, outGoDefFileName, outGoFileSuffix, outBytes)
}

func outputGoFile(outData types.DataHolder, goPackage string, outDir string) error {
	var wr strings.Builder

	if err := writeGoFileComments(&wr, outData.GetFileName(), outData.GetSheetName()); err != nil {
		return fmt.Errorf("generate Go file comments failed|file:%s|sheet:%s -> %w", outData.GetFileName(), outData.GetSheetName(), err)
	}
	if err := writeGoDeclaration(&wr, goPackage, "encoding/json", "os"); err != nil {
		return fmt.Errorf("generate Go declaration code failed|file:%s|sheet:%s -> %w", outData.GetFileName(), outData.GetSheetName(), err)
	}
	if err := writeGoEnum(&wr, outData.GetEnumTypes()); err != nil {
		return fmt.Errorf("generate Go enum code failed|file:%s|sheet:%s|enumTypes:%v -> %w", outData.GetFileName(), outData.GetSheetName(), outData.GetEnumTypes(), err)
	}
	if err := writeGoConfStruct(&wr, outData.GetHeaders(), outData.GetSheetName()); err != nil {
		return fmt.Errorf("generate Conf Go code failed|file:%s|sheet:%s|headers:{%+v} -> %w", outData.GetFileName(), outData.GetSheetName(), outData.GetHeaders(), err)
	}
	if err := writeConfMapStruct(&wr, outData.GetHeaders(), outData.GetSheetName()); err != nil {
		return fmt.Errorf("generate ConfMap Go code failed|file:%s|sheet:%s|headers:{%+v} -> %w", outData.GetFileName(), outData.GetSheetName(), outData.GetHeaders(), err)
	}

	outBytes, err := toOutBytes(wr.String())
	if err != nil {
		return fmt.Errorf("to output bytes failed -> %w", err)
	}
	return WriteToFile(outDir, outData.GetSheetName(), outGoFileSuffix, outBytes)
}

func WriteToGoFile(outData types.DataHolder, goPackage string, outDir string) error {
	if err := outputGoDefFile(goPackage, outDir); err != nil {
		return fmt.Errorf("generate go def file failed -> %w", err)
	}
	if err := outputGoFile(outData, goPackage, outDir); err != nil {
		return fmt.Errorf("generate go code file failed -> %w", err)
	}
	return nil
}
