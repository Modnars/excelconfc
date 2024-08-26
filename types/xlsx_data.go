package types

type XlsxDataHolder struct {
	FileName    string
	SheetName   string
	DataHeaders [][]string
	DataRows    [][]string
	EnumTypes   []*EnumTypeSt
	EnumValMap  map[string]*EnumValSt
}

func (xdh *XlsxDataHolder) GetFileName() string {
	return xdh.FileName
}

func (xdh *XlsxDataHolder) GetSheetName() string {
	return xdh.SheetName
}

func (xdh *XlsxDataHolder) GetHeaders() [][]string {
	return xdh.DataHeaders
}

func (xdh *XlsxDataHolder) GetData() [][]string {
	return xdh.DataRows
}

func (xdh *XlsxDataHolder) GetEnumTypes() []*EnumTypeSt {
	return xdh.EnumTypes
}

func (xdh *XlsxDataHolder) GetEnumValMap() map[string]*EnumValSt {
	return xdh.EnumValMap
}

var _ OutDataHolder = (*XlsxDataHolder)(nil)
