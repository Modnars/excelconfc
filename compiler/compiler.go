package compiler

import (
	"fmt"

	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
	"git.woa.com/modnarshen/excelconfc/reader/xlsx"
	"git.woa.com/modnarshen/excelconfc/translator"
	"git.woa.com/modnarshen/excelconfc/util"
	"git.woa.com/modnarshen/excelconfc/writer/golang"
	"git.woa.com/modnarshen/excelconfc/writer/json"
	"git.woa.com/modnarshen/excelconfc/writer/protobuf"
	"git.woa.com/modnarshen/excelconfc/writer/xml"
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
	groupFlag uint8
	parser    mcc.Parser
}

func newCompilerWithOptions(options *Options) *compiler {
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
		groupFlag: options.groupFlag,
	}
}

func New(options ...Option) Compiler {
	var opt Options
	for _, o := range options {
		o(&opt)
	}
	newCompiler := newCompilerWithOptions(&opt)
	newCompiler.parser = mcc.NewLRParser(mcc.NewGrammar(mcc.Productions))
	return newCompiler
}

func (c *compiler) Compile() error {
	xlsxData, err := xlsx.ReadFile(c.filePath, c.sheetName, c.enumSheet)
	if err != nil {
		return fmt.Errorf("exec ReadExcel failed|filePath:%s|sheetName:%s -> %w", c.filePath, c.sheetName, err)
	}

	nodes, err := translator.TransToASTNodes(xlsxData.Headers())
	if err != nil {
		return fmt.Errorf("exec NewTransToNodes failed|filePath:%s|sheet:%s|outDir:%s -> %w", c.filePath, c.sheetName, c.outDir, err)
	}
	astRoot, err := c.parser.BuildAST(nodes, OnReduce)
	if err != nil {
		return fmt.Errorf("exec BuildAST failed|filePath:%s|sheet:%s|outDir:%s -> %w", c.filePath, c.sheetName, c.outDir, err)
	}
	astRoot.SetType(c.sheetName)
	xlsxData.SetAST(mcc.FilterAST(astRoot, func(node mcc.ASTNode) bool {
		return node.GroupFlag()&c.groupFlag != 0
	}))
	if util.VerboseMode {
		mcc.PrintAST(xlsxData.AST(), 0)
	}

	if err := protobuf.WriteToFile(xlsxData, c.goPackage, c.outDir, c.addEnum); err != nil {
		return fmt.Errorf("exec WriteToProto failed|filePath:%s|sheet:%s|outDir:%s -> %w", c.filePath, c.sheetName, c.outDir, err)
	}
	if err := golang.WriteToFile(xlsxData, c.goPackage, c.outDir, c.addEnum); err != nil {
		return fmt.Errorf("exec golang.WriteToFile failed|file:%s|sheet:%s -> %w", xlsxData.FileName(), xlsxData.SheetName(), err)
	}
	if err := json.WriteToFile(xlsxData, c.outDir); err != nil {
		return fmt.Errorf("exec json.WriteToFile failed|filePath:%s|sheet:%s|outDir:%s -> %w", c.filePath, c.sheetName, c.outDir, err)
	}
	if err := xml.WriteToFile(xlsxData, c.outDir); err != nil {
		return fmt.Errorf("exec xml.WriteToFile failed|filePath:%s|sheet:%s|outDir:%s -> %w", c.filePath, c.sheetName, c.outDir, err)
	}

	return nil
}
