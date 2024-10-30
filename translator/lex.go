/*
 * @Author: modnarshen
 * @Date: 2024.10.25 10:11:54
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package translator

import (
	"fmt"
	"regexp"
	"strings"

	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/types"
)

var (
	bracketRegexp = regexp.MustCompile(`(\[|\]|\{|\})`)
)

// 删除传入的所有标记，如果成功删除任意标记，返回成功删除标记后的 s 和 true，否则返回 s 本身和 false
func removeAllMarks(str string, marks ...string) (string, bool) {
	removed := false
	for _, mark := range marks {
		if strings.Contains(str, mark) {
			str = strings.ReplaceAll(str, mark, "")
			removed = true
		}
	}
	return str, removed
}

func getLexVal(rawType, desc string) string {
	if !types.IsBasicType(rawType) {
		if desc == "E" {
			return types.LEX_ENUM
		}
		return types.LEX_ID
	}
	if desc == types.MARK_DESC_ARRAY && types.IsBasicType(rawType) {
		return types.LEX_ARRAY
	}
	if types.IsIntType(rawType) {
		return types.LEX_INT
	}
	if types.IsStringType(rawType) {
		return types.LEX_STRING
	}
	return types.TOK_NONE
}

func NewASTNodes(name string, fieldType string, desc string, colIdx int) []mcc.ASTNode {
	res := []mcc.ASTNode{}
	groupFlag := uint8(0)
	isFound := false
	if name, isFound = removeAllMarks(name, "|S", "|s"); isFound {
		groupFlag = groupFlag | types.GroupServer
	}
	if name, isFound = removeAllMarks(name, "|C", "|c"); isFound {
		groupFlag = groupFlag | types.GroupClient
	}
	if groupFlag == 0 { // if not set, set all flag
		groupFlag = groupFlag | types.GroupServer | types.GroupClient
	}
	if name == "[]" {
		res = append(res, mcc.NewASTNode("[]", name, fieldType, colIdx, groupFlag))
	} else {
		// 使用正则表达式进行切分，并保留分隔符
		parts := bracketRegexp.Split(name, -1)
		matches := bracketRegexp.FindAllString(name, -1)
		for i, part := range parts {
			if part != "" {
				lexVal := getLexVal(fieldType, desc)
				res = append(res, mcc.NewASTNode(lexVal, part, fieldType, colIdx, groupFlag))
			}
			if i < len(matches) {
				res = append(res, mcc.NewASTNode(matches[i], matches[i], "", colIdx, groupFlag))
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
