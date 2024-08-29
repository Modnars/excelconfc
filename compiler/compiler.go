package compiler

import (
	"fmt"

	"git.woa.com/modnarshen/excelconfc/reader"
	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/translator"
	"git.woa.com/modnarshen/excelconfc/writer"
	"git.woa.com/modnarshen/excelconfc/writer/json"
)

func Compile(filePath string, sheetName string, enumSheetName string, goPackage string, outDir string) error {
	xlsxData, err := reader.ReadXlsx(filePath, sheetName, enumSheetName)
	if err != nil {
		return fmt.Errorf("exec ReadExcel failed|filePath:%s|sheetName:%s -> %w", filePath, sheetName, err)
	}
	headers, err := translator.TransToNodes(xlsxData.GetHeaders())
	if err != nil {
		return fmt.Errorf("exec ParseHeaders failed -> %w", err)
	}
	root := translator.BuildNodeTree(headers)
	translator.PrintTree(root, 0)

	if err := json.WriteToFile(&translator.DataHolder{DataHolder: xlsxData, ASTRoot: root}, outDir); err != nil {
		return fmt.Errorf("exec json.WriteToFile failed|filePath:%s|sheet:%s|outDir:%s -> %w", filePath, sheetName, outDir, err)
	}
	if err := writer.WriteToProtoFile(xlsxData, goPackage, outDir); err != nil {
		return fmt.Errorf("exec WriteToProto failed|filePath:%s|sheet:%s|outDir:%s -> %w", filePath, sheetName, outDir, err)
	}
	if rules.DEBUG_MODE {
		if err := writer.WriteToGoFile(xlsxData, goPackage, outDir); err != nil {
			return fmt.Errorf("WriteToGo failed|filePath:%s|sheet:%s|outDir:%s -> %w", filePath, sheetName, outDir, err)
		}
		if err := writer.WriteToJsonFile(xlsxData, outDir); err != nil {
			return fmt.Errorf("WriteToJson failed|filePath:%s|sheet:%s|outDir:%s -> %w", filePath, sheetName, outDir, err)
		}
	}
	return nil
}
