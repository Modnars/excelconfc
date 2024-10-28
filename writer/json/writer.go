package json

import (
	"encoding/json"
	"fmt"

	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
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
		if subNode.LexVal() == types.MID_NODE_VEC {
			if !types.IsBDT(subNode.Type()) {
				vec := []any{}
				for _, ssubNode := range subNode.SubNodes() {
					vec = append(vec, buildLineData(ssubNode, rowData, evm))
				}
				result[subNode.Name()] = vec

			}
		} else if subNode.LexVal() == types.MID_NODE_FIELDS {
			result[subNode.Name()] = buildLineData(subNode, rowData, evm)
		} else {
			if subNode.ColIdx() >= len(rowData) {
				break
			}
			if val, err := CellValue(subNode, rowData[subNode.ColIdx()], evm); err != nil {
				util.LogError("Wrong CellValue|colIdx:%d", subNode.ColIdx())
			} else {
				result[subNode.Name()] = val
			}
		}
	}
	return result
}

func buildAllLineData(data types.DataHolder) ([]map[string]any, error) {
	allLineData := []map[string]any{}
	for _, rowData := range data.Data() {
		allLineData = append(allLineData, buildLineData(data.AST(), rowData, data.EnumValMap()))
	}
	return allLineData, nil
}

func WriteToFile(data types.DataHolder, outDir string) error {
	allLineData, err := buildAllLineData(data)
	if err != nil {
		return fmt.Errorf("build line data failed -> %w", err)
	}

	jsonBytes, err := json.MarshalIndent(map[string]any{"data": allLineData}, "", util.IndentSpace(1))
	if err != nil {
		return fmt.Errorf("json.MarshalIndent failed -> %w", err)
	}

	return writer.WriteToFile(outDir, data.SheetName(), outFileSuffix, jsonBytes)
}
