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
		if subNode.ColIdx() >= len(rowData) {
			return nil
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
			val, err := lex.CellValue(subNode, rowData[subNode.ColIdx()], evm)
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
			val, err := lex.CellValue(subNode, rowData[subNode.ColIdx()], evm)
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
	fmt.Fprintf(wr, "<%sMap>\n<all_infos>\n", data.SheetName())
	indent++
	for i, rowData := range data.Data() {
		fmt.Fprintf(wr, "%s<item>\n", util.IndentSpace(indent))
		if err := writeLineData(wr, data.AST(), rowData, data.EnumValMap(), indent+1); err != nil {
			return fmt.Errorf("row:%d,%w", rules.ROW_HEAD_MAX+1+i, err)
		}
		fmt.Fprintf(wr, "%s</item>\n", util.IndentSpace(indent))
	}
	indent--
	fmt.Fprintf(wr, "</all_infos>\n</%sMap>\n", data.SheetName())
	return nil
}

func WriteToFile(data lex.DataHolder, outDir string) error {
	var strBuilder strings.Builder
	if err := writeAllLineData(&strBuilder, data); err != nil {
		return err
	}
	return writer.WriteToFile(outDir, data.SheetName(), outFileSuffix, []byte(strBuilder.String()))
}
