package types

import "fmt"

// Enum Value Map
type EVM map[string]*EnumValSt

type FieldSt struct {
	Name       string     // 字段名（解析后）
	Type       string     // 字段类型（解析后，可供解析时直接取用的类型）
	RawType    string     // 原始类型（例如 Excel 中原始配置的类型）
	SubFields  []*FieldSt // 子字段（用于嵌套定义）
	Descriptor string     // 修饰符（比如 Excel 配置中使用 D 来修饰 string 类型为时间类型）
	Group      string     // 分组（区分前台客户端、后台服务器等）
	ColIdx     int        // 列坐标，用于索引源数据
}

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

type OutDataHolder interface {
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
