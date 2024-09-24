package compiler

import (
	"fmt"

	"git.woa.com/modnarshen/excelconfc/reader/xlsx"
	"git.woa.com/modnarshen/excelconfc/translator"
	"git.woa.com/modnarshen/excelconfc/writer/golang"
	"git.woa.com/modnarshen/excelconfc/writer/json"
	"git.woa.com/modnarshen/excelconfc/writer/protobuf"
)

type Compiler interface {
	Compile() error
}

type compiler struct {
	filePath  string
	sheetName string
	enumSheet string
	outDir    string
	goPackage string
	addEnum   bool
}

func newCompilerWithOptions(options *Options) Compiler {
	if options == nil {
		return nil
	}
	return &compiler{
		filePath:  options.filePath,
		sheetName: options.sheetName,
		enumSheet: options.enumSheet,
		outDir:    options.outDir,
		goPackage: options.goPackage,
		addEnum:   options.addEnum,
	}
}

func New(options ...Option) Compiler {
	var opt Options
	for _, o := range options {
		o(&opt)
	}
	return newCompilerWithOptions(&opt)
}

func (c *compiler) Compile() error {
	xlsxData, err := xlsx.ReadFile(c.filePath, c.sheetName, c.enumSheet)
	if err != nil {
		return fmt.Errorf("exec ReadExcel failed|filePath:%s|sheetName:%s -> %w", c.filePath, c.sheetName, err)
	}
	dataHolder, err := translator.NewDataHolder(
		translator.WithXlsxData(xlsxData),
	)
	if err != nil {
		return fmt.Errorf("NewDataHolder failed|file:%s|sheet:%s -> %w", c.filePath, c.sheetName, err)
	}
	if err := protobuf.WriteToFile(dataHolder, c.goPackage, c.outDir, c.addEnum); err != nil {
		return fmt.Errorf("exec WriteToProto failed|filePath:%s|sheet:%s|outDir:%s -> %w", c.filePath, c.sheetName, c.outDir, err)
	}
	if err := golang.WriteToFile(dataHolder, c.goPackage, c.outDir, c.addEnum); err != nil {
		return fmt.Errorf("exec golang.WriteToFile failed|file:%s|sheet:%s -> %w", dataHolder.FileName(), dataHolder.SheetName(), err)
	}
	if err := json.WriteToFile(dataHolder, c.outDir); err != nil {
		return fmt.Errorf("exec json.WriteToFile failed|filePath:%s|sheet:%s|outDir:%s -> %w", c.filePath, c.sheetName, c.outDir, err)
	}
	return nil
}
