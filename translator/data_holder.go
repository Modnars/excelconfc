package translator

import "git.woa.com/modnarshen/excelconfc/types"

type DataHolder interface {
	types.DataHolder
	AST() *Node
}

type mDataHolder struct {
	types.DataHolder
	astRoot *Node
}

func (dh *mDataHolder) AST() *Node {
	return dh.astRoot
}

var _ DataHolder = (*mDataHolder)(nil)

type NewDataHolderOption func(*mDataHolder) error

func WithXlsxData(xlsxData types.DataHolder) NewDataHolderOption {
	return func(dh *mDataHolder) error {
		dh.DataHolder = xlsxData
		if nodes, err := TransToNodes(xlsxData.Headers()); err != nil {
			return err
		} else {
			dh.astRoot = BuildNodeTree(nodes)
		}
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
