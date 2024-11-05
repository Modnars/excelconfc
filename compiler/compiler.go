package compiler

import (
	"fmt"

	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
	"git.woa.com/modnarshen/excelconfc/lex"
	"git.woa.com/modnarshen/excelconfc/util"
)

type Compiler interface {
	Compile(lex.DataHolder) error
}

type compiler struct {
	fileName  string
	sheetName string
	groupFlag uint8
	parser    mcc.Parser
}

func newCompilerWithOptions(options *Options) *compiler {
	if options == nil {
		return nil
	}
	return &compiler{
		fileName:  options.fileName,
		sheetName: options.sheetName,
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

func (c *compiler) Compile(data lex.DataHolder) error {
	nodes, err := lex.TransToASTNodes(data.Headers())
	if err != nil {
		return fmt.Errorf("exec TransToASTNodes failed|file:%s|sheet:%s -> %w", c.fileName, c.sheetName, err)
	}

	astRoot, err := c.parser.BuildAST(nodes, OnReduce)
	if err != nil {
		return fmt.Errorf("exec BuildAST failed|file:%s|sheet:%s -> %w", c.fileName, c.sheetName, err)
	}

	astRoot.SetType(c.sheetName)
	data.SetAST(mcc.FilterAST(astRoot, func(node mcc.ASTNode) bool {
		return node.GroupFlag()&c.groupFlag != 0
	}))

	if util.VerboseMode {
		mcc.PrintAST(data.AST(), 0)
	}

	return nil
}
