package types

import (
	"fmt"

	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
)

type EnumValSt struct {
	Name  string
	ID    any
	Notes string
}

type EnumTypeSt struct {
	Name     string
	Notes    string
	EnumVals []*EnumValSt
}

func (ev *EnumValSt) String() string {
	return fmt.Sprintf("%+v", *ev)
}

func (et *EnumTypeSt) String() string {
	return fmt.Sprintf("%+v", *et)
}

// Enum Value Map
type EVM map[string]*EnumValSt

type DataHolder interface {
	FileName() string
	SheetName() string
	Headers() [][]string
	Data() [][]string
	EnumTypes() []*EnumTypeSt
	EnumValMap() EVM
	AST() mcc.ASTNode
	SetAST(mcc.ASTNode)
}

type dataHolder struct {
	fileName   string
	sheetName  string
	headers    [][]string
	data       [][]string
	enumTypes  []*EnumTypeSt
	enumValMap EVM
	astRoot    mcc.ASTNode
}

func (h *dataHolder) FileName() string {
	return h.fileName
}

func (h *dataHolder) SheetName() string {
	return h.sheetName
}

func (h *dataHolder) Headers() [][]string {
	return h.headers
}

func (h *dataHolder) Data() [][]string {
	return h.data
}

func (h *dataHolder) EnumTypes() []*EnumTypeSt {
	return h.enumTypes
}

func (h *dataHolder) EnumValMap() EVM {
	return h.enumValMap
}

func (h *dataHolder) AST() mcc.ASTNode {
	return h.astRoot
}

func (h *dataHolder) SetAST(astRoot mcc.ASTNode) {
	h.astRoot = astRoot
}

var _ DataHolder = (*dataHolder)(nil)

func NewDataHolder(fileName, sheetName string, headers, data [][]string, enumTypes []*EnumTypeSt, enumValMap EVM) DataHolder {
	return &dataHolder{fileName: fileName,
		sheetName:  sheetName,
		headers:    headers,
		data:       data,
		enumTypes:  enumTypes,
		enumValMap: enumValMap}
}
