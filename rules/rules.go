package rules

import (
	"git.woa.com/modnarshen/excelconfc/util"
)

const (
	ROW_IDX_NAME = 0 // 名字定义在 Excel 第一行
	ROW_IDX_DESC = 1 // 修饰符定义在 Excel 第一行
	ROW_IDX_TYPE = 2 // 类型定义在 Excel 第三行
	ROW_IDX_NOTE = 4 // 注释定义在 Excel 第四行
	ROW_HEAD_MAX = 4 // headers 的最大行数
)

var (
	intTypes    = util.NewSet("int32", "uint32")
	stringTypes = util.NewSet("string")
)

func IsIntType(tp string) bool {
	return intTypes[tp]
}

func IsStringType(tp string) bool {
	return stringTypes[tp]
}
