package translator

import (
	"fmt"

	"git.woa.com/modnarshen/excelconfc/reader"
	"git.woa.com/modnarshen/excelconfc/writer"
)

func Translate(filePath string, sheetName string, enumSheetName string, goPackage string, outDir string) error {
	xlsxData, err := reader.ReadExcel(filePath, sheetName, enumSheetName)
	if err != nil {
		return fmt.Errorf("exec ReadExcel failed|filePath:%s|sheetName:%s -> %w", filePath, sheetName, err)
	}
	fmt.Printf("len:%d, contents: %+v\n", len(xlsxData.GetEnumData()), xlsxData.GetEnumData())
	if err := writer.WriteToProtoFile(xlsxData.GetHeaders(), filePath, sheetName, goPackage, outDir); err != nil {
		return fmt.Errorf("exec WriteToProto failed|filePath:%s|sheet:%s|outDir:%s -> %w", filePath, sheetName, outDir, err)
	}
	if err := writer.WriteToGoFile(xlsxData.GetHeaders(), filePath, sheetName, goPackage, outDir); err != nil {
		return fmt.Errorf("WriteToGo failed|filePath:%s|sheet:%s|outDir:%s -> %w", filePath, sheetName, outDir, err)
	}
	if err := writer.WriteToJsonFile(xlsxData.GetHeaders(), xlsxData.GetData(), filePath, sheetName, outDir); err != nil {
		return fmt.Errorf("WriteToJson failed|filePath:%s|sheet:%s|outDir:%s -> %w", filePath, sheetName, outDir, err)
	}
	return nil
}
