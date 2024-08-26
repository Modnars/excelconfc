package types

type XlsxDataHolder struct {
	FileName    string
	SheetName   string
	DataHeaders [][]string
	DataRows    [][]string
	EnumData    []*EnumTypeSt
	EnumMap     map[string]*EnumValSt
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

func (xdh *XlsxDataHolder) GetEnumData() []*EnumTypeSt {
	return xdh.EnumData
}

func (xdh *XlsxDataHolder) GetEnumMap() map[string]*EnumValSt {
	return xdh.EnumMap
}

var _ OutDataHolder = (*XlsxDataHolder)(nil)
