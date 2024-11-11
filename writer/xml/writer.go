/*
 * @Author: modnarshen
 * @Date: 2024.10.30 17:17:47
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package xml

import (
	"fmt"
	"io"
	"strings"

	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
	"git.woa.com/modnarshen/excelconfc/lex"
	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/util"
	"git.woa.com/modnarshen/excelconfc/writer"
)

const (
	outFileSuffix = ".ec.xml"
)

func writeLineData(wr io.Writer, astNode mcc.ASTNode, rowData []string, evm lex.EVM, indent int) error {
	if astNode.LexVal() != lex.MID_NODE_FIELDS {
		return nil
	}

	for _, subNode := range astNode.SubNodes() {
		cell := ""
		if subNode.ColIdx() < len(rowData) {
			cell = rowData[subNode.ColIdx()]
		}

		switch subNode.LexVal() {
		case lex.MID_NODE_FIELDS:
			fmt.Fprintf(wr, "%s<%s>\n", util.IndentSpace(indent), subNode.Name())
			indent++
			if err := writeLineData(wr, subNode, rowData, evm, indent); err != nil {
				return err
			}
			indent--
			fmt.Fprintf(wr, "%s</%s>\n", util.IndentSpace(indent), subNode.Name())

		case lex.MID_NODE_VEC:
			fmt.Fprintf(wr, "%s<%s>\n", util.IndentSpace(indent), subNode.Name())
			indent++
			if lex.IsBasicType(subNode.Type()) {
				for _, ssubNode := range subNode.SubNodes() {
					if writer.CanBeOmitted(ssubNode, rowData) {
						continue
					}
					val, err := lex.CellValue(ssubNode, rowData[ssubNode.ColIdx()], evm)
					if err != nil {
						return err
					}
					fmt.Fprintf(wr, "%s<item>%v</item>\n", util.IndentSpace(indent), val)
				}
			} else {
				for _, ssubNode := range subNode.SubNodes() {
					if writer.CanBeOmitted(ssubNode, rowData) {
						continue
					}
					fmt.Fprintf(wr, "%s<item>\n", util.IndentSpace(indent))
					indent++
					if err := writeLineData(wr, ssubNode, rowData, evm, indent); err != nil {
						return err
					}
					indent--
					fmt.Fprintf(wr, "%s</item>\n", util.IndentSpace(indent))
				}
			}
			indent--
			fmt.Fprintf(wr, "%s</%s>\n", util.IndentSpace(indent), subNode.Name())

		case lex.LEX_ARRAY:
			val, err := lex.CellValue(subNode, cell, evm)
			if err != nil {
				return err
			}
			if arrayVal, ok := val.([]any); ok {
				fmt.Fprintf(wr, "%s<%s>\n", util.IndentSpace(indent), subNode.Name())
				for _, item := range arrayVal {
					fmt.Fprintf(wr, "%s<item>%v</item>\n", util.IndentSpace(indent+1), item)
				}
				fmt.Fprintf(wr, "%s</%s>\n", util.IndentSpace(indent), subNode.Name())
			}

		default:
			val, err := lex.CellValue(subNode, cell, evm)
			if err != nil {
				return err
			}
			fmt.Fprintf(wr, "%s<%s>%v</%s>\n", util.IndentSpace(indent), subNode.Name(), val, subNode.Name())
		}
	}
	return nil
}

func writeAllLineData(wr io.Writer, data lex.DataHolder) error {
	indent := 0
	headLabel := data.SheetName()

	if len(data.AST().SubNodes()) <= 0 {
		return fmt.Errorf("there is no header fields")
	}

	switch data.ContainerType() {
	case rules.CONTAINER_TYPE_MAP:
		headLabel += "Map"
	case rules.CONTAINER_TYPE_VECTOR:
		headLabel += "Vector"
	}

	needCheckKey := data.ContainerType() == rules.CONTAINER_TYPE_MAP
	keyIdxes := []int{}
	// 只有 map 类型的底层容器才需要检查重复索引
	if needCheckKey {
		keyIdxes = lex.GetKeyFieldIdxes(data.AST())
	}
	confKeys := util.NewSet[string]()
	errMsgs := []string{}

	fmt.Fprintf(wr, "<%s>\n<all_infos>\n", headLabel)
	indent++
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
		fmt.Fprintf(wr, "%s<item>\n", util.IndentSpace(indent))
		if err := writeLineData(wr, data.AST(), rowData, data.EnumValMap(), indent+1); err != nil {
			return fmt.Errorf("row:%d,%w", rules.ROW_HEAD_MAX+1+i, err)
		}
		fmt.Fprintf(wr, "%s</item>\n", util.IndentSpace(indent))
	}
	indent--
	fmt.Fprintf(wr, "</all_infos>\n</%s>\n", headLabel)

	if len(errMsgs) > 0 {
		for _, errMsg := range errMsgs {
			util.LogError(errMsg)
		}
		return fmt.Errorf("config key error, please fix and try again")
	}

	return nil
}

func WriteToFile(data lex.DataHolder, outDir string) error {
	var strBuilder strings.Builder
	if err := writeAllLineData(&strBuilder, data); err != nil {
		return err
	}
	return writer.WriteToFile(outDir, data.SheetName(), outFileSuffix, []byte(strBuilder.String()))
}
