package translator

import (
	"git.woa.com/modnarshen/excelconfc/reader"
	"git.woa.com/modnarshen/excelconfc/util"
	"git.woa.com/modnarshen/excelconfc/writer"
)

func Translate(filePath string, sheetName string, goPackage string, outDir string) error {
	headers, excelRows, err := reader.ReadExcel(filePath, sheetName)
	if err != nil {
		util.LogError("ReadExcel failed|filePath:%s|sheetName:%s", filePath, sheetName)
		return err
	}
	if err := writer.WriteToProtoFile(headers, filePath, sheetName, goPackage, outDir); err != nil {
		util.LogError("WriteToProto failed|filePath:%s|sheet:%s|outDir:%s", filePath, sheetName, outDir)
		return err
	}
	if err := writer.WriteToGoFile(headers, filePath, sheetName, goPackage, outDir); err != nil {
		util.LogError("WriteToGo failed|filePath:%s|sheet:%s|outDir:%s", filePath, sheetName, outDir)
		return err
	}
	if err := writer.WriteToJsonFile(headers, excelRows, filePath, sheetName, outDir); err != nil {
		util.LogError("WriteToJson failed|filePath:%s|sheet:%s|outDir:%s", filePath, sheetName, outDir)
		return err
	}
	return nil
}
