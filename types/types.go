package types

import "fmt"

type EnumValSt struct {
	Name  string
	ID    string
	Notes string
}

type EnumTypeSt struct {
	Name     string
	Notes    string
	EnumVals []*EnumValSt
}

// Enum Value Map
type EVM map[string]*EnumValSt

type DataHolder interface {
	GetFileName() string
	GetSheetName() string
	GetHeaders() [][]string
	GetData() [][]string
	GetEnumTypes() []*EnumTypeSt
	GetEnumValMap() EVM
}

func (ev *EnumValSt) String() string {
	return fmt.Sprintf("%+v", *ev)
}

func (et *EnumTypeSt) String() string {
	return fmt.Sprintf("%+v", *et)
}
