package writer

import (
	"fmt"
	"go/format"
	"io"
	"path"
	"strings"

	"git.woa.com/modnarshen/excelconfc/code/template"
	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/util"
)

const (
	// Output file suffix `.ec.go` means Excel Config Go code.
	outGoFileSuffix      = ".ec.go"
	outGoDefFileName     = "excelconf.def"
	outGoDefFileFullName = outGoDefFileName + outGoFileSuffix
)

func toOutBytes(output string) ([]byte, error) {
	var outBytes []byte
	if rules.DEBUG_MODE {
		outBytes = []byte(output)
	} else {
		var err error
		outBytes, err = format.Source([]byte(output))
		if err != nil {
			return nil, fmt.Errorf("format %s failed|err:%w", outGoDefFileFullName, err)
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

func writeGoFileComments(wr io.Writer, filePath string, sheetName string) error {
	template.ExecuteTemplate(wr, template.TmplGoCommentsHead, nil)
	// 仅当 filePath 和 sheetName 有意义时才生成文件注释代码
	if filePath != "" && sheetName != "" {
		tmplParams := template.T{
			"SourceFile":  path.Base(filePath),
			"SourceSheet": sheetName,
		}
		if err := template.ExecuteTemplate(wr, template.TmplGoCommentsSource, tmplParams); err != nil {
			return fmt.Errorf("exectue template failed|tmplName:%s|err:%w", template.TmplGoCommentsSource, err)
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
		return fmt.Errorf("exectue template failed|tmplName:%s|err:%w", template.TmplGoCodeConfMap, err)
	}
	return nil
}

func outputGoDefFile(goPackage string, outDir string) error {
	var wr strings.Builder
	if err := writeGoFileComments(&wr, "", ""); err != nil {
		return fmt.Errorf("write Go file comments failed|fineName:%s|err:%w", outGoDefFileName, err)
	}
	if err := writeGoDeclaration(&wr, goPackage, "encoding/json", "time"); err != nil {
		return fmt.Errorf("write Go declaration failed|err:%w", err)
	}
	if err := template.ExecuteTemplate(&wr, template.TmplGoCodeDefDateTime, nil); err != nil {
		return fmt.Errorf("exectue template failed|tmplName:%s|err:%w", template.TmplGoCodeDefDateTime, err)
	}
	outBytes, err := toOutBytes(wr.String())
	if err != nil {
		return fmt.Errorf("to output bytes failed|err:%w", err)
	}
	return WriteToFile(outDir, "excelconf.def", outGoFileSuffix, outBytes)
}

func outputGoFile(headers [][]string, filePath string, sheetName string, goPackage string, outDir string) error {
	var wr strings.Builder

	if err := writeGoFileComments(&wr, filePath, sheetName); err != nil {
		return fmt.Errorf("generate Go file comments failed|fileName:%s|sheetName:%s|err:%w",
			path.Base(filePath), sheetName, err)
	}
	if err := writeGoDeclaration(&wr, goPackage, "encoding/json", "os"); err != nil {
		return fmt.Errorf("generate Go declaration code failed|fileName:%s|sheetName:%s|err:%w",
			path.Base(filePath), sheetName, err)
	}
	if err := writeGoConfStruct(&wr, headers, sheetName); err != nil {
		return fmt.Errorf("generate Conf Go code failed|fileName:%s|sheetName:%s|headers:{%+v}|err:%w",
			path.Base(filePath), sheetName, headers, err)
	}
	if err := writeConfMapStruct(&wr, headers, sheetName); err != nil {
		return fmt.Errorf("generate ConfMap Go code failed|fileName:%s|sheetName:%s|headers:{%+v}|err:%w",
			path.Base(filePath), sheetName, headers, err)
	}

	outBytes, err := toOutBytes(wr.String())
	if err != nil {
		return fmt.Errorf("to output bytes failed|err:%w", err)
	}
	return WriteToFile(outDir, sheetName, outGoFileSuffix, outBytes)
}

func WriteToGoFile(headers [][]string, filePath string, sheetName string, goPackage string, outDir string) error {
	if err := outputGoDefFile(goPackage, outDir); err != nil {
		return fmt.Errorf("generate go def file failed|err:%w", err)
	}
	if err := outputGoFile(headers, filePath, sheetName, goPackage, outDir); err != nil {
		return fmt.Errorf("generate go code file failed|err:%w", err)
	}
	return nil
}
