package translator

import (
	"fmt"

	"git.woa.com/modnarshen/excelconfc/reader"
	"git.woa.com/modnarshen/excelconfc/writer"
)

func Translate(filePath string, sheetName string, goPackage string, outDir string) error {
	headers, excelRows, err := reader.ReadExcel(filePath, sheetName)
	if err != nil {
		return fmt.Errorf("exec ReadExcel failed|filePath:%s|sheetName:%s|err:%w", filePath, sheetName, err)
	}
	if err := writer.WriteToProtoFile(headers, filePath, sheetName, goPackage, outDir); err != nil {
		return fmt.Errorf("exec WriteToProto failed|filePath:%s|sheet:%s|outDir:%s|err:%w", filePath, sheetName,
			outDir, err)
	}
	if err := writer.WriteToGoFile(headers, filePath, sheetName, goPackage, outDir); err != nil {
		return fmt.Errorf("WriteToGo failed|filePath:%s|sheet:%s|outDir:%s|err:%w", filePath, sheetName, outDir, err)
	}
	if err := writer.WriteToJsonFile(headers, excelRows, filePath, sheetName, outDir); err != nil {
		return fmt.Errorf("WriteToJson failed|filePath:%s|sheet:%s|outDir:%s|err:%w", filePath, sheetName, outDir, err)
	}
	return nil
}
