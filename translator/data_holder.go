package translator

import (
	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
	"git.woa.com/modnarshen/excelconfc/types"
)

type DataHolder interface {
	types.DataHolder
	AST() mcc.ASTNode
}

type mDataHolder struct {
	types.DataHolder
	astRoot mcc.ASTNode
}

func (dh *mDataHolder) AST() mcc.ASTNode {
	return dh.astRoot
}

var _ DataHolder = (*mDataHolder)(nil)

type NewDataHolderOption func(*mDataHolder) error

func WithXlsxData(xlsxData types.DataHolder) NewDataHolderOption {
	return func(dh *mDataHolder) error {
		dh.DataHolder = xlsxData
		return nil
	}
}

func WithAST(astRoot mcc.ASTNode) NewDataHolderOption {
	return func(dh *mDataHolder) error {
		dh.astRoot = astRoot
		return nil
	}
}

func NewDataHolder(options ...NewDataHolderOption) (DataHolder, error) {
	dataHolder := &mDataHolder{}
	for _, option := range options {
		if err := option(dataHolder); err != nil {
			return nil, err
		}
	}
	return dataHolder, nil
}
