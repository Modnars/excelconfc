/*
 * @Author: modnarshen
 * @Date: 2024.10.25 10:11:54
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package lex

import (
	"fmt"
	"regexp"
	"strings"

	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
	"git.woa.com/modnarshen/excelconfc/rules"
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

func getLexInfo(rawType, desc string) (LexMark, Token, error) {
	switch rawType {
	case TOK_TYPE_BOOL:
		if desc == LEX_ARRAY {
			return LEX_ARRAY, rawType, nil
		}
		return LEX_BOOL, rawType, nil

	case TOK_TYPE_INT32, TOK_TYPE_UINT32, TOK_TYPE_INT64, TOK_TYPE_UINT64:
		if desc == LEX_ARRAY {
			return LEX_ARRAY, rawType, nil
		}
		return LEX_INT, rawType, nil

	case TOK_TYPE_STRING, TOK_TYPE_FSTRING, TOK_TYPE_FTEXT:
		if desc == TOK_DESC_DATETIME {
			return LEX_STRING, LEX_STRING, nil
		} else if desc == LEX_ARRAY {
			return LEX_ARRAY, rawType, nil
		}
		return LEX_STRING, LEX_STRING, nil

	case TOK_TYPE_DATETIME: // built-in type
		if desc == LEX_ARRAY {
			return TOK_NONE, TOK_NONE, fmt.Errorf("datetime can not use `array`")
		}
		return LEX_STRING, rawType, nil

	default:
		if desc == TOK_DESC_ENUM {
			return LEX_ENUM, rawType, nil
		}
		return LEX_ID, rawType, nil
	}
	// return TOK_NONE, TOK_NONE, fmt.Errorf("unknown lex node|type:%v|desc:%v", rawType, desc)
}

func NewASTNodes(name string, fieldType string, desc string, colIdx int) ([]mcc.ASTNode, error) {
	res := []mcc.ASTNode{}
	groupFlag := uint8(0)
	isFound := false
	if name, isFound = removeAllMarks(name, "|S", "|s"); isFound {
		groupFlag = groupFlag | GroupServer
	}
	if name, isFound = removeAllMarks(name, "|C", "|c"); isFound {
		groupFlag = groupFlag | GroupClient
	}
	if groupFlag == 0 { // if not set, set all flag
		groupFlag = groupFlag | GroupServer | GroupClient
	}
	if name == "[]" {
		res = append(res, mcc.NewASTNode("[]", name, fieldType, desc, colIdx, groupFlag))
	} else {
		// 使用正则表达式进行切分，并保留分隔符
		parts := bracketRegexp.Split(name, -1)
		matches := bracketRegexp.FindAllString(name, -1)
		for i, part := range parts {
			if part != "" {
				lexVal, nodeType, err := getLexInfo(fieldType, desc)
				if err != nil {
					return nil, err
				}
				res = append(res, mcc.NewASTNode(lexVal, part, nodeType, desc, colIdx, groupFlag))
			}
			if i < len(matches) {
				res = append(res, mcc.NewASTNode(matches[i], matches[i], "", "", colIdx, groupFlag))
			}
		}
	}
	return res, nil
}

func TransToASTNodes(headers [][]string) ([]mcc.ASTNode, error) {
	if len(headers) < rules.ROW_HEAD_MAX {
		return nil, fmt.Errorf("invalid line count in headers|rowNum:%d", len(headers))
	}
	nodes := []mcc.ASTNode{}
	for colIdx := range headers[rules.ROW_IDX_NAME] {
		newNodes, err := NewASTNodes(
			headers[rules.ROW_IDX_NAME][colIdx],
			headers[rules.ROW_IDX_TYPE][colIdx],
			headers[rules.ROW_IDX_DESC][colIdx],
			colIdx,
		)
		if err != nil {
			return nil, err
		}
		if newNodes != nil {
			nodes = append(nodes, newNodes...)
		}
	}
	return nodes, nil
}
