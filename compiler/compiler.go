package compiler

import (
	"fmt"

	"git.woa.com/modnarshen/excelconfc/reader"
	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/translator"
	"git.woa.com/modnarshen/excelconfc/writer/golang"
	"git.woa.com/modnarshen/excelconfc/writer/json"
	"git.woa.com/modnarshen/excelconfc/writer/protobuf"
)

type Compiler struct {
	filePath  string
	sheetName string
	enumSheet string
	outDir    string
	goPackage string
	addEnum   bool
}

type CompileOption func(*Compiler)

func WithFilePath(filePath string) CompileOption {
	return func(c *Compiler) {
		c.filePath = filePath
	}
}

func WithSheetName(sheetName string) CompileOption {
	return func(c *Compiler) {
		c.sheetName = sheetName
	}
}

func WithEnumSheet(enumSheet string) CompileOption {
	return func(c *Compiler) {
		c.enumSheet = enumSheet
	}
}

func WithOutDir(outDir string) CompileOption {
	return func(c *Compiler) {
		c.outDir = outDir
	}
}

func WithGoPackage(goPackage string) CompileOption {
	return func(c *Compiler) {
		c.goPackage = goPackage
	}
}

func WithAddEnum(addEnum bool) CompileOption {
	return func(c *Compiler) {
		c.addEnum = addEnum
	}
}

func New(options ...CompileOption) *Compiler {
	c := &Compiler{
		enumSheet: rules.DEFAULT_ENUM_SHEET_NAME,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

func (c *Compiler) Compile() error {
	xlsxData, err := reader.ReadXlsx(c.filePath, c.sheetName, c.enumSheet)
	if err != nil {
		return fmt.Errorf("exec ReadExcel failed|filePath:%s|sheetName:%s -> %w", c.filePath, c.sheetName, err)
	}
	headers, err := translator.TransToNodes(xlsxData.GetHeaders())
	if err != nil {
		return fmt.Errorf("exec TransToNodes failed -> %w", err)
	}
	dataHolder := &translator.DataHolder{
		DataHolder: xlsxData,
		ASTRoot:    translator.BuildNodeTree(headers),
	}
	if err := json.WriteToFile(dataHolder, c.outDir); err != nil {
		return fmt.Errorf("exec json.WriteToFile failed|filePath:%s|sheet:%s|outDir:%s -> %w", c.filePath, c.sheetName, c.outDir, err)
	}
	if err := golang.WriteToFile(dataHolder, c.goPackage, c.outDir); err != nil {
		return fmt.Errorf("exec golang.WriteToFile failed|file:%s|sheet:%s -> %w", dataHolder.GetFileName(), dataHolder.GetSheetName(), err)
	}
	if err := protobuf.WriteToFile(dataHolder, c.goPackage, c.outDir); err != nil {
		return fmt.Errorf("exec WriteToProto failed|filePath:%s|sheet:%s|outDir:%s -> %w", c.filePath, c.sheetName, c.outDir, err)
	}
	return nil
}
