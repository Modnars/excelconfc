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
}
