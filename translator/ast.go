/*
 * @Author: modnarshen
 * @Date: 2024.10.25 10:11:54
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package translator

import (
	"fmt"
	"regexp"

	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/types"
)

func NewASTNodes(name string, fieldType string, desc string, colIdx int) []mcc.ASTNode {
	res := []mcc.ASTNode{}
	groupFlag := 0
	isFound := false
	if name, isFound = removeAllMarks(name, "|S", "|s"); isFound {
		groupFlag = groupFlag | 0b01
	}
	if name, isFound = removeAllMarks(name, "|C", "|c"); isFound {
		groupFlag = groupFlag | 0b10
	}
	if name == "[]" {
		// res = append(res, &Node{Name: "", lexVal: "[]"})
		res = append(res, mcc.NewASTNode("[]"))
	} else {
		re := regexp.MustCompile(`(\[|\]|\{|\})`)

		// 使用正则表达式进行切分，并保留分隔符
		parts := re.Split(name, -1)
		matches := re.FindAllString(name, -1)
		for i, part := range parts {
			if part != "" {
				lexVal := "id"
				if types.IsIntType(fieldType) {
					lexVal = "int"
				} else if types.IsStringType(fieldType) {
					lexVal = "string"
				} else if desc == "E" {
					lexVal = "enum"
				}
				// res = append(res, &Node{Name: part, lexVal: lexVal, ColIdx: colIdx, GroupFlag: groupFlag})
				res = append(res, mcc.NewASTNode(lexVal))
			}
			if i < len(matches) {
				// res = append(res, &Node{Name: matches[i], lexVal: matches[i], ColIdx: colIdx, GroupFlag: groupFlag})
				res = append(res, mcc.NewASTNode(matches[i]))
			}
		}
	}
	return res
}

func TransToASTNodes(headers [][]string) ([]mcc.ASTNode, error) {
	if len(headers) < rules.ROW_HEAD_MAX {
		return nil, fmt.Errorf("invalid line count in headers|rowNum:%d", len(headers))
	}
	nodes := []mcc.ASTNode{}
	for colIdx := range headers[rules.ROW_IDX_NAME] {
		newNodes := NewASTNodes(
			headers[rules.ROW_IDX_NAME][colIdx],
			headers[rules.ROW_IDX_TYPE][colIdx],
			headers[rules.ROW_IDX_DESC][colIdx],
			colIdx,
		)
		if newNodes != nil {
			nodes = append(nodes, newNodes...)
		}
	}
	return nodes, nil
}
