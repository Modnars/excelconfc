package xlsx

import (
	"git.woa.com/modnarshen/excelconfc/types"
)

type DataHolder struct {
	fileName   string
	sheetName  string
	headers    [][]string
	data       [][]string
	enumTypes  []*types.EnumTypeSt
	enumValMap types.EVM
}

func (h *DataHolder) FileName() string {
	return h.fileName
}

func (h *DataHolder) SheetName() string {
	return h.sheetName
}

func (h *DataHolder) Headers() [][]string {
	return h.headers
}

func (h *DataHolder) Data() [][]string {
	return h.data
}

func (h *DataHolder) EnumTypes() []*types.EnumTypeSt {
	return h.enumTypes
}

func (h *DataHolder) EnumValMap() types.EVM {
	return h.enumValMap
}

// var _ types.DataHolder = (*DataHolder)(nil)

type NewDataHolderOption func(*DataHolder)

func WithFileName(fileName string) NewDataHolderOption {
	return func(dh *DataHolder) {
		dh.fileName = fileName
	}
}

func WithSheetName(sheetName string) NewDataHolderOption {
	return func(dh *DataHolder) {
		dh.sheetName = sheetName
	}
}

func WithHeaders(headers [][]string) NewDataHolderOption {
	return func(dh *DataHolder) {
		dh.headers = headers
	}
}

func WithData(data [][]string) NewDataHolderOption {
	return func(dh *DataHolder) {
		dh.data = data
	}
}

func WithEnumTypes(enumTypes []*types.EnumTypeSt) NewDataHolderOption {
	return func(dh *DataHolder) {
		dh.enumTypes = enumTypes
	}
}

func WithEnumValMap(enumValMap types.EVM) NewDataHolderOption {
	return func(dh *DataHolder) {
		dh.enumValMap = enumValMap
	}
}

func NewDataHolder(options ...NewDataHolderOption) *DataHolder {
	xdh := &DataHolder{}
	for _, opt := range options {
		opt(xdh)
	}
	return xdh
}
