package translator

import "git.woa.com/modnarshen/excelconfc/types"

type DataHolder struct {
	types.DataHolder
	ASTRoot *Node
}

func (dh *DataHolder) GetAST() *Node {
	return dh.ASTRoot
}

type NewDataHolderOption func(*DataHolder) error

func WithXlsxData(xlsxData types.DataHolder) NewDataHolderOption {
	return func(dh *DataHolder) error {
		dh.DataHolder = xlsxData
		if nodes, err := TransToNodes(xlsxData.Headers()); err != nil {
			return err
		} else {
			dh.ASTRoot = BuildNodeTree(nodes)
		}
		return nil
	}
}

func NewDataHolder(options ...NewDataHolderOption) (*DataHolder, error) {
	dataHolder := &DataHolder{}
	for _, option := range options {
		if err := option(dataHolder); err != nil {
			return nil, err
		}
	}
	return dataHolder, nil
}
