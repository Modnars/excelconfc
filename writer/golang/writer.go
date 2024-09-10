package golang

import (
	"fmt"
	"go/format"
	"io"
	"strings"

	"git.woa.com/modnarshen/excelconfc/code/template"
	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/translator"
	"git.woa.com/modnarshen/excelconfc/types"
	"git.woa.com/modnarshen/excelconfc/util"
	"git.woa.com/modnarshen/excelconfc/writer"
)

const (
	// Output file suffix `.ec.go` means Excel Config Go code.
	outFileSuffix    = ".ec.go"
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

func writeDeclPackage(wr io.Writer, goPackage string) {
	fmt.Fprintf(wr, "\npackage %s\n", util.GetPackageName(goPackage))
}

func writeDeclImports(wr io.Writer, packages ...string) {
	if len(packages) <= 0 {
		return
	} else if len(packages) == 1 {
		fmt.Fprintf(wr, "\nimport \"%s\"\n", packages[0])
	}
	outLines := []string{}
	outLines = append(outLines, "\nimport (")
	for _, packageName := range packages {
		outLines = append(outLines, util.IndentSpace(1)+"\""+packageName+"\"")
	}
	outLines = append(outLines, ")")
	fmt.Fprintf(wr, strings.Join(outLines, "\n")+"\n")
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

func writeFileComments(wr io.Writer, fileName string, sheetName string) error {
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
		fmt.Fprintf(wr, "\n")
	}
	return nil
}

func writeDeclaration(wr io.Writer, goPackage string, importPkgs ...string) error {
	writeDeclPackage(wr, goPackage)
	writeDeclImports(wr, importPkgs...)
	return nil
}

func writeEnum(wr io.Writer, enumTypes []*types.EnumTypeSt) error {
	indent := 0
	fmt.Fprintf(wr, "\n")
	for _, enumType := range enumTypes {
		fmt.Fprintf(wr, "type %s int32\n", enumType.Name)
	}

	fmt.Fprintf(wr, "\nconst (")
	for _, enumType := range enumTypes {
		indent++
		fmt.Fprintf(wr, "\n")
		for _, enumVal := range enumType.EnumVals {
			fmt.Fprintf(wr, "%s%s %s = %s\n", util.IndentSpace(indent), enumVal.Name, enumType.Name, enumVal.ID)
		}
		indent--
	}
	fmt.Fprintf(wr, ")\n")

	fmt.Fprintf(wr, "\nvar (")
	for _, enumType := range enumTypes {
		indent++
		fmt.Fprintf(wr, "\n%s%s_name = map[int32]string{\n", util.IndentSpace(indent), enumType.Name)
		indent++
		for _, enumVal := range enumType.EnumVals {
			fmt.Fprintf(wr, "%s%s: \"%s\",\n", util.IndentSpace(indent), enumVal.ID, enumVal.Name)
		}
		indent--
		fmt.Fprintf(wr, "%s}\n", util.IndentSpace(indent))
		fmt.Fprintf(wr, "%s%s_value = map[string]int32{\n", util.IndentSpace(indent), enumType.Name)
		indent++
		for _, enumVal := range enumType.EnumVals {
			fmt.Fprintf(wr, "%s\"%s\": %s,\n", util.IndentSpace(indent), enumVal.Name, enumVal.ID)
		}
		indent--
		fmt.Fprintf(wr, "%s}\n", util.IndentSpace(indent))
		indent--
	}
	fmt.Fprintf(wr, ")\n")

	for _, enumType := range enumTypes {
		fmt.Fprintf(wr, "\nfunc(x %s) String() string {\n", enumType.Name)
		indent++
		fmt.Fprintf(wr, "return %s_name[int32(x)]\n", enumType.Name)
		indent--
		fmt.Fprintf(wr, "}\n")
	}
	return nil
}

func writeStructMap(wr io.Writer, headers [][]string, sheetName string) error {
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

func outputDefFile(goPackage string, outDir string) error {
	var wr strings.Builder
	if err := writeFileComments(&wr, "", ""); err != nil {
		return fmt.Errorf("write Go file comments failed|fineName:%s -> %w", outGoDefFileName, err)
	}
	if err := writeDeclaration(&wr, goPackage, "encoding/json", "time"); err != nil {
		return fmt.Errorf("write Go declaration failed -> %w", err)
	}
	if err := template.ExecuteTemplate(&wr, template.TmplGoCodeDefDateTime, nil); err != nil {
		return fmt.Errorf("exectue template failed|tmplName:%s -> %w", template.TmplGoCodeDefDateTime, err)
	}
	outBytes, err := toOutBytes(wr.String())
	if err != nil {
		return fmt.Errorf("to output bytes failed -> %w", err)
	}
	return writer.WriteToFile(outDir, outGoDefFileName, outFileSuffix, outBytes)
}

func collectStructFields(field writer.Field, structFields []writer.Field) []writer.Field {
	if field.IsStructDecl() || field.IsVectorDecl() {
		for _, subField := range field.SubNodes {
			structFields = collectStructFields(subField, structFields)
		}
		if field.IsStructDecl() {
			structFields = append(structFields, field)
		}
	}
	return structFields
}

func writeStruct(wr io.Writer, data *translator.DataHolder, sheetName string) error {
	rootNode := &translator.Node{
		Name:     sheetName,
		SubNodes: data.ASTRoot.SubNodes,
		Type:     types.TOK_TYPE_ROOT_STRUCT,
		DataType: sheetName,
	}
	structFields := []writer.Field{}
	indent := 0
	structFields = collectStructFields(rootNode, structFields)
	doneStructSet := util.NewSet[string]()
	for _, structField := range structFields {
		if doneStructSet.Contains(structField.DataType) {
			continue
		}
		fmt.Fprintf(wr, "\ntype %s struct {\n", structField.DataType)
		indent++
		for _, subField := range structField.SubNodes {
			if subField.IsVectorDecl() { // repeated
				fmt.Fprintf(wr, "%s%s []%s\n", util.IndentSpace(indent), util.SnakeToPascal(subField.Name), subField.DataType)
			} else if subField.IsStructDecl() {
				fmt.Fprintf(wr, "%s%s %s\n", util.IndentSpace(indent), util.SnakeToPascal(subField.Name), subField.DataType)
			} else if subField.IsEnum() {
				fmt.Fprintf(wr, "%s%s %s\n", util.IndentSpace(indent), util.SnakeToPascal(subField.Name), subField.DataType)
			} else {
				fmt.Fprintf(wr, "%s%s %s\n", util.IndentSpace(indent), util.SnakeToPascal(subField.Name), subField.Type)
			}
		}
		indent--
		fmt.Fprintf(wr, "}\n")
		doneStructSet.Add(structField.DataType)
	}
	return nil
}

func outputSrcFile(data *translator.DataHolder, goPackage string, outDir string, addEnum bool) error {
	wr := &strings.Builder{}

	if err := writeFileComments(wr, data.FileName(), data.SheetName()); err != nil {
		return fmt.Errorf("generate Go file comments failed|file:%s|sheet:%s -> %w", data.FileName(), data.SheetName(), err)
	}
	if err := writeDeclaration(wr, goPackage, "encoding/json", "os"); err != nil {
		return fmt.Errorf("generate Go declaration code failed|file:%s|sheet:%s -> %w", data.FileName(), data.SheetName(), err)
	}
	if addEnum {
		if err := writeEnum(wr, data.EnumTypes()); err != nil {
			return fmt.Errorf("generate Go enum code failed|file:%s|sheet:%s|enumTypes:%v -> %w", data.FileName(), data.SheetName(), data.EnumTypes(), err)
		}
	}
	if err := writeStruct(wr, data, data.SheetName()); err != nil {
		return fmt.Errorf("generate Conf Go code failed|file:%s|sheet:%s -> %w", data.FileName(), data.SheetName(), err)
	}
	if err := writeStructMap(wr, data.Headers(), data.SheetName()); err != nil {
		return fmt.Errorf("generate ConfMap Go code failed|file:%s|sheet:%s -> %w", data.FileName(), data.SheetName(), err)
	}

	outBytes, err := toOutBytes(wr.String())
	if err != nil {
		return fmt.Errorf("to output bytes failed -> %w", err)
	}
	return writer.WriteToFile(outDir, data.SheetName(), outFileSuffix, outBytes)
}

func WriteToFile(data *translator.DataHolder, goPackage string, outDir string, addEnum bool) error {
	if err := outputDefFile(goPackage, outDir); err != nil {
		return fmt.Errorf("generate go def file failed -> %w", err)
	}
	if err := outputSrcFile(data, goPackage, outDir, addEnum); err != nil {
		return fmt.Errorf("generate go code file failed -> %w", err)
	}
	return nil
}
