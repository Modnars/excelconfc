package json

import (
	"fmt"
	"io"
	"strings"

	"git.woa.com/modnarshen/excelconfc/code/template"
	"git.woa.com/modnarshen/excelconfc/translator"
	"git.woa.com/modnarshen/excelconfc/util"
	"git.woa.com/modnarshen/excelconfc/writer"
)

const (
	outFileSuffix = ".ec.json"
)

func writeRowData(wr io.Writer, nodes []*translator.Node, rowData []string, indent int, isLastElem bool) error {
	fmt.Fprintf(wr, "%s{\n", util.IndentSpace(indent))
	indent++
	for idx, node := range nodes {
		if node.ColIdx >= len(rowData) {
			break
		}
		comma := ","
		if idx == len(nodes)-1 {
			comma = ""
		}
		if node.IsVectorDecl() {
			if len(node.SubNodes) <= 0 {
				fmt.Fprintf(wr, "%s\"%s\": [ ]\n%s", util.IndentSpace(indent), node.Name, comma)
			} else {
				fmt.Fprintf(wr, "%s\"%s\": [\n", util.IndentSpace(indent), node.Name)
				writeRowData(wr, node.SubNodes, rowData, indent+1, node.ColIdx >= len(rowData)-1)
				fmt.Fprintf(wr, "%s]%s\n", util.IndentSpace(indent), comma)
				continue
			}
		}
		fmt.Fprintf(wr, "%s\"%s\": %s%s\n", util.IndentSpace(indent), node.Name, CellVal(node, rowData[node.ColIdx]), comma)
		if len(node.SubNodes) > 0 {
			writeRowData(wr, node.SubNodes, rowData, indent+1, node.ColIdx >= len(rowData)-1)
		}
	}
	indent--
	comma := ","
	if isLastElem {
		comma = ""
	}
	fmt.Fprintf(wr, "%s}%s\n", util.IndentSpace(indent), comma)
	return nil
}

func writeFieldData1(wr io.Writer, field *Field, vals []string, indent int, isLastElem bool) bool {
	if field.ColIdx >= len(vals) {
		return true
	}
	if field.IsVectorDecl() {
		fmt.Fprintf(wr, "%s\"%s\": [\n", util.IndentSpace(indent), field.Name)
	} else if field.IsStructDecl() {
		if field.ColIdx == 0 { // fake struct decl
			fmt.Fprintf(wr, "%s{\n", util.IndentSpace(indent))
		} else {
			fmt.Fprintf(wr, "%s\"%s\": {\n", util.IndentSpace(indent), field.Name)
		}
	} else {
		if field.ColIdx < len(vals) {
			fmt.Fprintf(wr, "%s\"%s\": %s", util.IndentSpace(indent), field.Name, CellVal(field, vals[field.ColIdx]))
		}
	}

	right := false
	for idx, subField := range field.SubNodes {
		right = right || writeFieldData1(wr, subField, vals, indent+1, idx == len(field.SubNodes)-1)
	}

	// 不加 ',' 的情况有三种：
	// 1. 明确是最后一个可输出元素
	// 2. 父结点明确指定当前结点是最后一个元素（比如结构体定义结点指定最后一个字段元素是最后一个元素）
	// 3. 如果当前结点的子结点（递归下去）是最后一个可输出元素（由此推出根节点实际上也是不会输出 ',' 的）
	comma := ","
	if field.ColIdx == len(vals)-1 || isLastElem || right {
		comma = ""
	}
	if field.IsVectorDecl() {
		fmt.Fprintf(wr, "%s]%s\n", util.IndentSpace(indent), comma)
	} else if field.IsStructDecl() {
		fmt.Fprintf(wr, "%s}%s\n", util.IndentSpace(indent), comma)
	} else {
		if field.ColIdx < len(vals) {
			fmt.Fprintf(wr, "%s\n", comma)
		}
	}
	return field.ColIdx == len(vals)-1
}

func writeDrive(wr io.Writer, data *translator.DataHolder, indent int, isLastElem bool) error {
	fmt.Fprintf(wr, "%s\"data\": [\n", util.IndentSpace(indent))
	indent++
	for idx, rowData := range data.GetData() {
		writeFieldData1(wr, data.ASTRoot, rowData, indent, idx == len(data.GetData())-1)
		if idx != len(data.GetData())-1 {
			fmt.Fprintf(wr, "%s,\n", util.IndentSpace(indent))
		}
	}
	indent--
	comma := ","
	if isLastElem {
		comma = ""
	}
	fmt.Fprintf(wr, "%s]%s\n", util.IndentSpace(indent), comma)
	return nil
}

func writeFieldData(wr io.Writer, data *translator.DataHolder, indent int, isLastElem bool) error {
	fmt.Fprintf(wr, "%s\"data\": [\n", util.IndentSpace(indent))
	indent++
	for idx, rowData := range data.GetData() {
		if err := writeRowData(wr, data.ASTRoot.SubNodes, rowData, indent, idx == len(data.GetData())-1); err != nil {
			return err
		}
	}
	indent--
	comma := ","
	if isLastElem {
		comma = ""
	}
	fmt.Fprintf(wr, "%s]%s\n", util.IndentSpace(indent), comma)
	return nil
}

func WriteToFile(data *translator.DataHolder, outDir string) error {
	wr := &strings.Builder{}
	indent := 0

	fmt.Fprintf(wr, "{\n")
	indent++
	tmplParams := template.T{
		"Indentation": util.IndentSpace(indent),
		"File":        data.GetFileName(),
		"Sheet":       data.GetSheetName(),
		"OutDir":      outDir,
	}
	if err := template.ExecuteTemplate(wr, template.TmplJsonFields, tmplParams); err != nil {
		return fmt.Errorf("exectue template failed|tmplName:%s -> %w", template.TmplJsonFields, err)
	}
	// writeFieldData(wr, data, indent, true)
	writeDrive(wr, data, indent, true)
	fmt.Fprintf(wr, "}\n")
	return writer.WriteToFile(outDir, data.GetSheetName(), outFileSuffix, []byte(wr.String()))
}
