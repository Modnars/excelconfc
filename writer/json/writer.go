package json

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"

	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
	"git.woa.com/modnarshen/excelconfc/lex"
	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/util"
	"git.woa.com/modnarshen/excelconfc/writer"
)

const (
	outFileSuffix = ".ec.json"
)

func buildLineData(astNode mcc.ASTNode, rowData []string, evm lex.EVM) (map[string]any, error) {
	if astNode.LexVal() != lex.MID_NODE_FIELDS {
		return nil, nil
	}

	result := make(map[string]any)
	for _, subNode := range astNode.SubNodes() {
		// if subNode.ColIdx() >= len(rowData) {
		// 	break
		// }

		switch subNode.LexVal() {
		case lex.MID_NODE_FIELDS:
			val, err := buildLineData(subNode, rowData, evm)
			if err != nil {
				return nil, err
			}
			result[subNode.Name()] = val

		case lex.MID_NODE_VEC:
			vec := []any{}
			if lex.IsBasicType(subNode.Type()) {
				for _, ssubNode := range subNode.SubNodes() {
					if writer.CanBeOmitted(ssubNode, rowData) {
						continue
					}
					val, err := lex.CellValue(ssubNode, rowData[ssubNode.ColIdx()], evm)
					if err != nil {
						return nil, fmt.Errorf("col:%s -> %w", util.ColumnName(ssubNode.ColIdx()), err)
					}
					vec = append(vec, val)
				}
			} else {
				for _, ssubNode := range subNode.SubNodes() {
					if writer.CanBeOmitted(ssubNode, rowData) {
						continue
					}
					val, err := buildLineData(ssubNode, rowData, evm)
					if err != nil {
						return nil, err
					}
					vec = append(vec, val)
				}
			}
			result[subNode.Name()] = vec

		default:
			cell := ""
			if subNode.ColIdx() < len(rowData) {
				cell = rowData[subNode.ColIdx()]
			}
			val, err := lex.CellValue(subNode, cell, evm)
			if err != nil {
				return nil, fmt.Errorf("col:%s -> %w", util.ColumnName(subNode.ColIdx()), err)
			}
			result[subNode.Name()] = val
		}
	}
	return result, nil
}

func buildAllLineData(data lex.DataHolder) ([]map[string]any, error) {
	if len(data.AST().SubNodes()) <= 0 {
		return nil, fmt.Errorf("there is no header fields")
	}

	needCheckKey := data.ContainerType() == rules.CONTAINER_TYPE_MAP
	keyIdxes := []int{}
	if needCheckKey {
		keyIdxes = lex.GetKeyFieldIdxes(data.AST())
	}
	confKeys := util.NewSet[string]()
	errMsgs := []string{}

	allLineData := []map[string]any{}
	for i, rowData := range data.Data() {
		rowKey, err := lex.GenConfKey(keyIdxes, rowData)
		if err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("row:%d|%s", rules.ROW_HEAD_MAX+1+i, err.Error()))
		}
		if confKeys.Contains(rowKey) {
			errMsgs = append(errMsgs, fmt.Sprintf("row:%d|found a repeated key|key:%s", rules.ROW_HEAD_MAX+1+i, rowKey))
		}
		if needCheckKey {
			confKeys.Add(rowKey)
		}
		lineData, err := buildLineData(data.AST(), rowData, data.EnumValMap())
		if err != nil {
			return nil, fmt.Errorf("row:%d,%w", rules.ROW_HEAD_MAX+1+i, err)
		}
		allLineData = append(allLineData, lineData)
	}

	if len(errMsgs) > 0 {
		for _, errMsg := range errMsgs {
			util.LogError(errMsg)
		}
		return nil, fmt.Errorf("config key error, please fix and try again")
	}

	return allLineData, nil
}

func WriteToFile(data lex.DataHolder, outDir string) error {
	allLineData, err := buildAllLineData(data)
	if err != nil {
		return fmt.Errorf("invalid data -> %w", err)
	}

	jsonBytes, err := jsoniter.ConfigFastest.MarshalIndent(map[string]any{"data": allLineData}, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent failed -> %w", err)
	}

	return writer.WriteToFile(outDir, data.SheetName(), outFileSuffix, jsonBytes)
}
