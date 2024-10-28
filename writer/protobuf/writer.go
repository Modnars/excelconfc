package protobuf

import (
	"fmt"
	"io"
	"path"
	"strings"

	"git.woa.com/modnarshen/excelconfc/code/template"
	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
	"git.woa.com/modnarshen/excelconfc/types"
	"git.woa.com/modnarshen/excelconfc/util"
	"git.woa.com/modnarshen/excelconfc/writer"
)

const (
	outFileSuffix = ".ec.proto"
)

func writeFileComment(wr io.Writer, filePath string, sheetName string) error {
	tmplParams := template.T{
		"SourceFile":  path.Base(filePath),
		"SourceSheet": sheetName,
	}
	if err := template.ExecuteTemplate(wr, template.TmplProtoComments, tmplParams); err != nil {
		return fmt.Errorf("exectue template failed|tmplName:%s -> %w", template.TmplProtoComments, err)
	}
	return nil
}

func writeDeclaration(wr io.Writer, goPackage string) error {
	tmplParams := template.T{
		"PackageName": util.GetPackageName(goPackage),
		"GoPackage":   goPackage,
	}
	if err := template.ExecuteTemplate(wr, template.TmplProtoCodePackage, tmplParams); err != nil {
		return fmt.Errorf("exectue template failed|tmplName:%s -> %w", template.TmplProtoCodePackage, err)
	}
	return nil
}

func collectMessages(astNode mcc.ASTNode, result []mcc.ASTNode) []mcc.ASTNode {
	for _, subNode := range astNode.SubNodes() {
		result = collectMessages(subNode, result)
	}
	if astNode.LexVal() == types.MID_NODE_FIELDS && astNode.Type() != types.TOK_NONE {
		result = append(result, astNode)
	}
	return result
}

func writeMessage(wr io.Writer, rootNode mcc.ASTNode) error {
	messages := collectMessages(rootNode, nil)
	indent := 0
	doneMsgSet := util.NewSet[string]()
	for _, message := range messages {
		if doneMsgSet.Contains(message.Type()) {
			continue
		}
		fmt.Fprintf(wr, "\nmessage %s {\n", message.Type())
		indent++
		msgFieldNo := 0
		for _, field := range message.SubNodes() {
			msgFieldNo++
			if field.LexVal() == types.MID_NODE_VEC {
				fmt.Fprintf(wr, "%srepeated %s %s = %d;\n", util.IndentSpace(indent), field.Type(), field.Name(), msgFieldNo)
			} else {
				fmt.Fprintf(wr, "%s%s %s = %d;\n", util.IndentSpace(indent), field.Type(), field.Name(), msgFieldNo)
			}
		}
		indent--
		fmt.Fprintf(wr, "}\n")
		doneMsgSet.Add(message.Type())
	}
	return nil
}

func writeEnum(wr io.Writer, enumTypes []*types.EnumTypeSt) error {
	indent := 0
	for _, enumType := range enumTypes {
		fmt.Fprintf(wr, "\nenum %s {\n", enumType.Name)
		indent++
		for _, enumVal := range enumType.EnumVals {
			fmt.Fprintf(wr, "%s%s = %v;\n", util.IndentSpace(indent), enumVal.Name, enumVal.ID)
		}
		indent--
		fmt.Fprintf(wr, "}\n")
	}
	return nil
}

func WriteToFile(data types.DataHolder, goPackage string, outDir string, addEnum bool) error {
	wr := &strings.Builder{}

	if err := writeFileComment(wr, data.FileName(), data.SheetName()); err != nil {
		return fmt.Errorf("generate proto file comment failed|file:%s|sheet:%s -> %w", data.FileName(), data.SheetName(), err)
	}
	if err := writeDeclaration(wr, goPackage); err != nil {
		return fmt.Errorf("generate proto declaration failed|file:%s|sheet:%s -> %w", data.FileName(), data.SheetName(), err)
	}
	if addEnum { // 只有明确指明需要添加枚举定义时，才将枚举定义输出
		if err := writeEnum(wr, data.EnumTypes()); err != nil {
			return fmt.Errorf("generate proto message failed|file:%s|sheet:%s|enumTypes:{%+v} -> %w", data.FileName(), data.SheetName(), data.EnumTypes(), err)
		}
	}
	if err := writeMessage(wr, data.AST()); err != nil {
		return fmt.Errorf("generate proto message failed|file:%s|sheet:%s -> %w", data.FileName(), data.SheetName(), err)
	}

	return writer.WriteToFile(outDir, data.SheetName(), outFileSuffix, []byte(wr.String()))
}
