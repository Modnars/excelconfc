package writer

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/util"
	wrtmpl "git.woa.com/modnarshen/excelconfc/writer/template"
)

const (
	// Output file suffix `.ec.go` means Excel Config Go code.
	outGoFileSuffix      = ".ec.go"
	outGoDefFileName     = "excelconf.def"
	outGoDefFileFullName = outGoDefFileName + outGoFileSuffix
)

func getPackageName(goPackage string) string {
	splitCh := ';'
	index := 0
	if strings.ContainsRune(goPackage, splitCh) {
		index = strings.IndexRune(goPackage, splitCh) + 1
	} else {
		index = strings.IndexRune(goPackage, '/') + 1
	}
	return goPackage[index:]
}

func genGoPackageDeclaration(goPackage string) string {
	return fmt.Sprintf("\npackage %s\n", getPackageName(goPackage))
}

func genGoImportDeclaration(packages ...string) string {
	if len(packages) <= 0 {
		return ""
	} else if len(packages) == 1 {
		return fmt.Sprintf("\nimport \"%s\"\n", packages[0])
	}
	outLines := []string{}
	outLines = append(outLines, "\nimport (")
	for _, packageName := range packages {
		outLines = append(outLines, indentSpace(1)+"\""+packageName+"\"")
	}
	outLines = append(outLines, ")")
	return strings.Join(outLines, "\n") + "\n"
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

func writeGoFileComment(wr io.Writer, filePath string, sheetName string) error {
	wrf(wr, wrtmpl.TextGoFileCommentTitle)
	// 仅当 filePath 和 sheetName 有意义时才生成文件注释代码
	if filePath != "" && sheetName != "" {
		tmplParams := struct {
			SourceFile  string
			SourceSheet string
		}{
			SourceFile:  path.Base(filePath),
			SourceSheet: sheetName,
		}
		if err := wrtmpl.GetWrTemplate(wrtmpl.WrTmplGoFileCommentSource).Execute(wr, tmplParams); err != nil {
			util.LogError("exectue template failed|tmplName:%s", wrtmpl.WrTmplGoFileCommentSource)
			return err
		}
	} else {
		// 此处显式添加换行，这是因为 Title 本身是不会有换行的，以此来保证注释代码的连贯性。
		wrf(wr, "\n")
	}
	return nil
}

func writeGoDeclaration(wr io.Writer, goPackage string, importPkgs ...string) error {
	wrf(wr, genGoPackageDeclaration(goPackage))
	wrf(wr, genGoImportDeclaration(importPkgs...))
	return nil
}

func writeGoConfStruct(wr io.Writer, headers [][]string, sheetName string) error {
	indent := 0
	wrf(wr, "\ntype %s struct {\n", sheetName)
	indent += 1
	for i, name := range headers[rules.ROW_IDX_NAME] {
		wrf(wr, "%s%s %s `json:\"%s\"`\n",
			indentSpace(indent), util.SnakeToPascal(name),
			genGoStructFieldType(headers[rules.ROW_IDX_TYPE][i], headers[rules.ROW_IDX_DESC][i]),
			name)
	}
	indent -= 1
	wrf(wr, "}\n")
	return nil
}

func writeGoConfMapStruct(wr io.Writer, headers [][]string, sheetName string) error {
	confKeyType, confKeyField := genGoConfKeyInfo(headers)
	tmplParams := struct {
		XXConf         string
		XXConfMap      string
		XXConfKeyType  string
		XXConfKeyField string
	}{
		XXConf:         sheetName,
		XXConfMap:      sheetName + "Map",
		XXConfKeyType:  confKeyType,
		XXConfKeyField: confKeyField,
	}

	if err := wrtmpl.GetWrTemplate(wrtmpl.WrTmplGoXXConfMap).Execute(wr, tmplParams); err != nil {
		util.LogError("exectue template failed|tmplName:%s", wrtmpl.WrTmplGoXXConfMap)
		return err
	}
	return nil
}

func outputGoDefFile(goPackage string, outDir string) error {
	var wr strings.Builder

	if err := writeGoFileComment(&wr, "", ""); err != nil {
		util.LogError("generate Go File comment failed|fineName:%s|err:{%+v}", outGoDefFileName)
		return err
	}
	if err := writeGoDeclaration(&wr, goPackage, "encoding/json", "time"); err != nil {
		util.LogError("generate Go Declaration code failed|fineName:%s|err:{%+v}", outGoDefFileName)
		return err
	}

	wrf(&wr, wrtmpl.TextGoDateTimeTypeDef)

	outBytes, err := toOutBytes(wr.String())
	if err != nil {
		return err
	}
	os.WriteFile(genOutFilePath(outDir, "excelconf.def", outGoFileSuffix), outBytes, outFilePerm)

	return nil
}

func outputGoFile(headers [][]string, filePath string, sheetName string, goPackage string, outDir string) error {
	var wr strings.Builder

	if err := writeGoFileComment(&wr, filePath, sheetName); err != nil {
		util.LogError("generate Go File comment failed|fileName:%s|sheetName:%s|headers:{%+v}|err:{%+v}", path.Base(filePath), sheetName, headers, err)
		return err
	}
	if err := writeGoDeclaration(&wr, goPackage, "encoding/json", "os"); err != nil {
		util.LogError("generate Go Declaration code failed|fileName:%s|sheetName:%s|headers:{%+v}|err:{%+v}", path.Base(filePath), sheetName, headers, err)
		return err
	}
	if err := writeGoConfStruct(&wr, headers, sheetName); err != nil {
		util.LogError("generate Conf Go code failed|fileName:%s|sheetName:%s|headers:{%+v}|err:{%+v}", path.Base(filePath), sheetName, headers, err)
		return err
	}
	if err := writeGoConfMapStruct(&wr, headers, sheetName); err != nil {
		util.LogError("generate ConfMap Go code failed|fileName:%s|sheetName:%s|headers:{%+v}|err:{%+v}", path.Base(filePath), sheetName, headers, err)
		return err
	}

	outBytes, err := toOutBytes(wr.String())
	if err != nil {
		return err
	}
	os.WriteFile(genOutFilePath(outDir, sheetName, outGoFileSuffix), outBytes, outFilePerm)

	return nil
}

func WriteToGoFile(headers [][]string, filePath string, sheetName string, goPackage string, outDir string) error {
	if err := outputGoDefFile(goPackage, outDir); err != nil {
		util.LogError("generate go def file failed|err:{%+v}", err)
		return err
	}
	if err := outputGoFile(headers, filePath, sheetName, goPackage, outDir); err != nil {
		util.LogError("generate go code file failed|err:{%+v}", err)
		return err
	}
	return nil
}
