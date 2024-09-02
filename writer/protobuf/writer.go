package protobuf

import (
	"fmt"
	"io"
	"path"
	"strings"

	"git.woa.com/modnarshen/excelconfc/code/template"
	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/translator"
	"git.woa.com/modnarshen/excelconfc/types"
	"git.woa.com/modnarshen/excelconfc/util"
	"git.woa.com/modnarshen/excelconfc/writer"
)

const (
	outProtoFileSuffix = ".ec.proto"
)

func writeProtoFileComment(wr io.Writer, filePath string, sheetName string) error {
	tmplParams := template.T{
		"SourceFile":  path.Base(filePath),
		"SourceSheet": sheetName,
	}
	if err := template.ExecuteTemplate(wr, template.TmplProtoComments, tmplParams); err != nil {
		return fmt.Errorf("exectue template failed|tmplName:%s -> %w", template.TmplProtoComments, err)
	}
	return nil
}

func writeProtoDecl(wr io.Writer, goPackage string) error {
	tmplParams := template.T{
		"PackageName": util.GetPackageName(goPackage),
		"GoPackage":   goPackage,
	}
	if err := template.ExecuteTemplate(wr, template.TmplProtoCodePackage, tmplParams); err != nil {
		return fmt.Errorf("exectue template failed|tmplName:%s -> %w", template.TmplProtoCodePackage, err)
	}
	return nil
}

func writeProtoMessage(wr io.Writer, headers [][]string, sheetName string) error {
	indent := 0
	msgLabelNumber := 0
	fmt.Fprintf(wr, "\nmessage %s {\n", sheetName)
	indent += 1
	for i, name := range headers[rules.ROW_IDX_NAME] {
		msgLabelNumber += 1
		fmt.Fprintf(wr, "%s%s %s = %d;\n", util.IndentSpace(indent), headers[rules.ROW_IDX_TYPE][i], name, msgLabelNumber)
	}
	indent -= 1
	fmt.Fprintf(wr, "}\n")
	return nil
}

func collectMessages(field writer.Field, messages []writer.Field) []writer.Field {
	if field.IsStructDecl() || field.IsVectorDecl() {
		for _, subField := range field.SubNodes {
			messages = collectMessages(subField, messages)
		}
		if field.IsStructDecl() {
			messages = append(messages, field)
		}
	}
	return messages
}

func writeMessage(wr io.Writer, data *translator.DataHolder, sheetName string) error {
	rootMessage := &translator.Node{
		Name:     sheetName,
		SubNodes: data.ASTRoot.SubNodes,
		Type:     types.TOK_TYPE_ROOT_STRUCT,
		RawType:  sheetName,
	}
	messages := collectMessages(rootMessage, nil)
	indent := 0
	doneMessageSet := util.NewSet[string]()
	for _, message := range messages {
		if doneMessageSet.Contains(message.RawType) {
			continue
		}
		fmt.Fprintf(wr, "\nmessage %s {\n", message.RawType)
		indent++
		msgLabelNumber := 0
		for _, subMessage := range message.SubNodes {
			msgLabelNumber++
			if subMessage.IsVectorDecl() {
				fmt.Fprintf(wr, "%srepeated %s %s = %d;\n", util.IndentSpace(indent), subMessage.RawType, subMessage.Name, msgLabelNumber)
			} else if subMessage.IsStructDecl() {
				fmt.Fprintf(wr, "%s%s %s = %d;\n", util.IndentSpace(indent), subMessage.RawType, subMessage.Name, msgLabelNumber)
			} else if subMessage.IsEnum() {
				fmt.Fprintf(wr, "%s%s %s = %d;\n", util.IndentSpace(indent), subMessage.RawType, subMessage.Name, msgLabelNumber)
			} else {
				fmt.Fprintf(wr, "%s%s %s = %d;\n", util.IndentSpace(indent), subMessage.RawType, subMessage.Name, msgLabelNumber)
			}
		}
		indent--
		fmt.Fprintf(wr, "}\n")
		doneMessageSet.Add(message.RawType)
	}
	return nil
}

func writeProtoEnum(wr io.Writer, enumTypes []*types.EnumTypeSt) error {
	indent := 0
	for _, enumType := range enumTypes {
		fmt.Fprintf(wr, "\nenum %s {\n", enumType.Name)
		indent++
		for _, enumVal := range enumType.EnumVals {
			fmt.Fprintf(wr, "%s%s = %s;\n", util.IndentSpace(indent), enumVal.Name, enumVal.ID)
		}
		indent--
		fmt.Fprintf(wr, "}\n")
	}
	return nil
}

func WriteToFile(data *translator.DataHolder, goPackage string, outDir string) error {
	wr := &strings.Builder{}

	if err := writeProtoFileComment(wr, data.GetFileName(), data.GetSheetName()); err != nil {
		return fmt.Errorf("generate proto file comment failed|file:%s|sheet:%s -> %w", data.GetFileName(), data.GetSheetName(), err)
	}
	if err := writeProtoDecl(wr, goPackage); err != nil {
		return fmt.Errorf("generate proto declaration failed|file:%s|sheet:%s -> %w", data.GetFileName(), data.GetSheetName(), err)
	}
	if err := writeProtoEnum(wr, data.GetEnumTypes()); err != nil {
		return fmt.Errorf("generate proto message failed|file:%s|sheet:%s|enumTypes:{%+v} -> %w", data.GetFileName(), data.GetSheetName(), data.GetEnumTypes(), err)
	}
	if err := writeMessage(wr, data, data.GetSheetName()); err != nil {
		return fmt.Errorf("generate proto message failed|file:%s|sheet:%s -> %w", data.GetFileName(), data.GetSheetName(), err)
	}

	return writer.WriteToFile(outDir, data.GetSheetName(), outProtoFileSuffix, []byte(wr.String()))
}
