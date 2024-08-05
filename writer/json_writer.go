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
	outJsonFileSuffix = ".ec.json"
)

func getCellValByType(cell string, tp string) string {
	if isIntType(tp) {
		if cell == "" {
			return "0"
		}
		return cell
	} else if isStringType(tp) {
		return fmt.Sprintf("\"%s\"", cell)
	}
	return fmt.Sprintf("\"%s\"", cell)
}

func writeJsonRowsData(wr io.Writer, headers [][]string, excelRows [][]string, indent int, isLastElem bool) error {
	maxColNum := len(headers[rules.ROW_IDX_NAME])

	wrf(wr, "%s\"data\": [\n", indentSpace(indent))
	indent += 1
	for i, row := range excelRows {
		// 确保只会解析 headers[rules.ROW_IDX_NAME] 界定列范围内的元素
		row := row[:min(len(row), maxColNum)]
		wrf(wr, "%s{\n", indentSpace(indent))

		indent += 1
		for j, cell := range row {
			fieldName := headers[rules.ROW_IDX_NAME][j]
			fieldType := headers[rules.ROW_IDX_TYPE][j]
			if j == len(row)-1 {
				wrf(wr, "%s\"%s\": %s\n", indentSpace(indent), fieldName, getCellValByType(cell, fieldType))
			} else {
				wrf(wr, "%s\"%s\": %s,\n", indentSpace(indent), fieldName, getCellValByType(cell, fieldType))
			}
		}
		indent -= 1

		if i == len(excelRows)-1 {
			wrf(wr, "%s}\n", indentSpace(indent))
		} else {
			wrf(wr, "%s},\n", indentSpace(indent))
		}
	}
	indent -= 1
	if isLastElem {
		wrf(wr, "%s]\n", indentSpace(indent))
	} else {
		wrf(wr, "%s],\n", indentSpace(indent))
	}
	return nil
}

func WriteToJsonFile(headers [][]string, excelRows [][]string, filePath string, sheetName string, outDir string) error {
	indent := 0
	var wr strings.Builder

	wr.WriteString("{\n")
	indent += 1
	tmplParams := struct {
		IndentSpace string
		Filepath    string
		Basename    string
		Sheet       string
		Outdir      string
	}{
		IndentSpace: indentSpace(indent),
		Filepath:    filePath,
		Basename:    path.Base(filePath),
		Sheet:       sheetName,
		Outdir:      outDir,
	}
	if err := wrtmpl.GetWrTemplate(wrtmpl.WrTmplJsonFields).Execute(&wr, tmplParams); err != nil {
		util.LogError("exectue template failed|tmplName:%s", wrtmpl.WrTmplJsonFields)
		return err
	}

	// 传入一个 isLastElem 的参数来指示写入的内容是否是最后一个元素。如果是，就省略最后的 `,`，否则添加上 `,`
	if err := writeJsonRowsData(&wr, headers, excelRows, indent, true); err != nil {
		util.LogError("parse excel rows to JSON failed|file:%s|sheet:%s", filePath, sheetName)
		return err
	}
	indent -= 1
	wrf(&wr, "}\n")

	os.WriteFile(genOutFilePath(outDir, sheetName, outJsonFileSuffix), []byte(wr.String()), outFilePerm)
	return nil
}
