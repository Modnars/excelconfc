package json

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"git.woa.com/modnarshen/excelconfc/code/template"
	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
	"git.woa.com/modnarshen/excelconfc/translator"
	"git.woa.com/modnarshen/excelconfc/types"
	"git.woa.com/modnarshen/excelconfc/util"
	"git.woa.com/modnarshen/excelconfc/writer"
)

const (
	outFileSuffix = ".ec.json"
)

func buildLineData(astNode mcc.ASTNode, rowData []string, evm types.EVM) map[string]any {
	if astNode.LexVal() != types.MID_NODE_FIELDS {
		return nil
	}

	result := make(map[string]any)
	for _, subNode := range astNode.SubNodes() {
		if subNode.LexVal() == "Node@VEC" {
			for _, ssubNode := range subNode.SubNodes() {
				if ssubNode.LexVal() == types.MID_NODE_VEC_ADT_ITEMS {
					vec := []any{}
					for _, sssubNode := range ssubNode.SubNodes() {
						vec = append(vec, buildLineData(sssubNode, rowData, evm))
					}
					result[ssubNode.Name()] = vec
				}
			}
		} else if subNode.LexVal() == "Node@STRUCT" {
			for _, ssubNode := range subNode.SubNodes() {
				if ssubNode.LexVal() != types.MID_NODE_FIELDS {
					continue
				}
				result[subNode.Name()] = buildLineData(ssubNode, rowData, evm)
				break
			}
		} else if subNode.LexVal() == "Node@BDT" {
			ssubNode := subNode.SubNodes()[0]
			if ssubNode.ColIdx() >= len(rowData) {
				continue
			}
			if val, err := CellValue(ssubNode, rowData[ssubNode.ColIdx()], evm); err != nil {
				util.LogError("Wrong CellValue|colIdx:%d", ssubNode.ColIdx())
			} else {
				result[ssubNode.Name()] = val
			}
		}
	}
	return result
}

func writeDataRows(wr io.Writer, data translator.DataHolder, indent int, isLastElem bool) error {
	fmt.Fprintf(wr, "%s\"data\": ", util.IndentSpace(indent))
	rowMaps := []map[string]any{}
	for _, rowData := range data.Data() {
		rowMaps = append(rowMaps, buildLineData(data.AST(), rowData, data.EnumValMap()))
	}
	if b, err := json.MarshalIndent(rowMaps, util.IndentSpace(indent), util.IndentSpace(1)); err == nil {
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

func WriteToFile(data translator.DataHolder, outDir string) error {
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
