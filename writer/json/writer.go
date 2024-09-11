package json

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"git.woa.com/modnarshen/excelconfc/code/template"
	"git.woa.com/modnarshen/excelconfc/translator"
	"git.woa.com/modnarshen/excelconfc/types"
	"git.woa.com/modnarshen/excelconfc/util"
	"git.woa.com/modnarshen/excelconfc/writer"
)

const (
	outFileSuffix = ".ec.json"
)

// 按照 AST 树结构生成单行 json 数据映射
func buildRowMap(m map[string]any, node *translator.Node, rowData []string, evm types.EVM) map[string]any {
	if len(node.SubNodes) > 0 {
		if node.IsVectorDecl() {
			vec := []any{}
			for _, subNode := range node.SubNodes {
				vec = append(vec, buildRowMap(make(map[string]any), subNode, rowData, evm))
			}
			m[node.Name] = vec
		} else if node.IsStructDecl() {
			subM := make(map[string]any)
			for _, subNode := range node.SubNodes {
				subM = buildRowMap(subM, subNode, rowData, evm)
			}
			if types.IsRealStruct(node.Type) {
				m[node.Name] = subM
			} else { // 如果是 VecStruct 或 RootSruct，其本身不是一个有意义的结点，其子结点才是有意义的
				return subM
			}
		}
	} else {
		if node.ColIdx < len(rowData) {
			if val, err := CellValue(node, rowData[node.ColIdx], evm); err != nil {
				util.LogError("Wrong CellValue|colIdx:%d", node.ColIdx)
			} else {
				m[node.Name] = val
			}
		}
	}
	return m
}

func writeDataRows(wr io.Writer, data *translator.DataHolder, indent int, isLastElem bool) error {
	fmt.Fprintf(wr, "%s\"data\": ", util.IndentSpace(indent))
	rowMaps := []map[string]any{}
	for _, rowData := range data.Data() {
		rowMaps = append(rowMaps, buildRowMap(nil, data.ASTRoot, rowData, data.EnumValMap()))
	}
	if b, err := json.MarshalIndent(rowMaps, util.IndentSpace(indent), "    "); err == nil {
		wr.Write(b)
	} else {
		return err
	}
	comma := ","
	if isLastElem {
		comma = ""
	}
	fmt.Fprintf(wr, "%s\n", comma)
	return nil
}

func WriteToFile(data *translator.DataHolder, outDir string) error {
	wr := &strings.Builder{}
	indent := 0

	fmt.Fprintf(wr, "{\n")
	indent++
	tmplParams := template.T{
		"Indentation": util.IndentSpace(indent),
		"File":        data.FileName(),
		"Sheet":       data.SheetName(),
		"OutDir":      outDir,
	}
	if err := template.ExecuteTemplate(wr, template.TmplJsonFields, tmplParams); err != nil {
		return fmt.Errorf("exectue template failed|tmplName:%s -> %w", template.TmplJsonFields, err)
	}
	writeDataRows(wr, data, indent, true)
	fmt.Fprintf(wr, "}\n")
	return writer.WriteToFile(outDir, data.SheetName(), outFileSuffix, []byte(wr.String()))
}
