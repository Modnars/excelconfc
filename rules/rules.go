package rules

const (
	ROW_IDX_NAME = 0 // 名字定义在 Excel 第一行
	ROW_IDX_DESC = 1 // 修饰符定义在 Excel 第一行
	ROW_IDX_TYPE = 2 // 类型定义在 Excel 第三行
	ROW_IDX_NOTE = 4 // 注释定义在 Excel 第四行
	ROW_HEAD_MAX = 4 // headers 的最大行数

	DEFAULT_ENUM_SHEET_NAME = "ENUM_DESC" // 默认枚举定义表名
)
