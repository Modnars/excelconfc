# Excel-Config-Compiler

## reader

### API: [excelize](https://xuri.me/excelize/zh-hans/)

通过 excelize 的 行迭代器 和 列迭代器 来获取表格中的数据，对于一些 excelize 与以往 Python 下 xlrd 表现的差异，在此做一个说明。

- 数据配置表（DataSheet）

excelize 的 `func (f *excelize.File) GetRows(sheet string, opts ...excelize.Options) ([][]string, error)` 函数将 Excel 表格数据直接映射成“一张二维字符串表”，这极大地遍历了后续处理数据的操作，但同时这个二维表本身的一些特征依然需要开发者留意。如果某一行中没有填入任何元素，那么这一行的值将会是一个 `[]`，因此如果想要通过行与行列数对齐的方式来随机访问行数据，这样的操作会导致 slice 访问越界。举个例子：

| act\_id | title | begin\_time | end\_time |
| :-: | :-: | :-: | :-: |
| | | | |
| uint32 | string | string | string |
| ID | 名称 | 开始时间 | 结束时间 |

此时的第二行配置数据为空，那么通过 excelize 读取得到的数据将会是

```text
[act_id title begin_time end_time]
[]
[uint32 string string string]
[ID 名称 开始时间 结束时间]
```

这几行配置作为表头信息配置（用于指示从第 5 行开始的数据以什么样的形式来进行解析），在解析数据时往往需要直接访问来确定表格内数据如何“变形”，比如同样是字符串，如果第二行标识了 `D`，就代表表格数据是一个时间（D，DateTime），如果没有标识，就单纯代表单元格数据是一个字符串，那么运行时读取 `headers[1][col_idx]` 就一定会得到一个越界访问的 panic！因此，对于此类需要支持“对齐项访问”的场景，需要对 API 进行一个封装修饰来满足实际的业务需求，因此在 `reader.readDataSheet` 中对表头 headers 进行了填充处理，使得表头的每一行的列数都与表头第一行的列数保持一致。

而对于下方的配置数据而言，如果通过正常的 range 迭代遍历得到的数据，就已经足够支持数据的生成，那么此时就不对下方的配置数据进行确保数据列数对齐的填充。

因此对于通过 `reader.readDataSheet` 获取得到的数据，往往可以这样去想象它们的样子：

```text
[a a a a a a a]
[b b b b b b b]
[c c c c c c c]
[d d d d d d d]
[1 1 1]
[2 2 2 2 2]
[3 3 3 3]
```

由此来进行开发可以规避一定场景下的数据访问失败的问题。

- 枚举定义配置表（EnumSheet）

按照 `{EnumType}EnumLabel` 的方式来定义一个枚举类，用 `{EnumLabel}EnumValueLabel` 开始的行来定义此枚举类中具体的枚举值。

| ... | ... | ... | ... |
| :-: | :-: | :-: | :-: |
| {ActType}活动类型 | | |
| [活动类型]签到活动 | 1 | ACT\_TYPE\_CHECK\_IN |

有以下几个约束来确保枚举值定义的简洁性、规范性。

1. 如果定义了枚举类，那么接下来的行中出现的枚举值定义都将尝试用于填充此枚举类。比如说上述表格中如果下方多了一行 `[任务类型]签到任务` 开始的行，那么**此行数据会被直接跳过**。
2. 直到出现另一个枚举类标签的定义或配置数据行结束，才会终止当前枚举类的枚举值获取。比如上述表格下方多了一行 `{TaskType}任务类型` 开始的行，那么从此行以下出现的以 `[任务类型]` 起始的行均被填充到 `TaskType` 枚举类中作为其枚举值定义。

因此，规范的配置方式应当是这样的：

| ... | ... | ... | ... |
| :-: | :-: | :-: | :-: |
| {ActType}活动类型 | | | |
| [活动类型]签到活动 | 1 | ACT\_TYPE\_CHECK\_IN |
| [活动类型]登录活动 | 1 | ACT\_TYPE\_LOG\_IN |
| | | | |
| {TaskType}任务类型 | | | |
| Label | ID | Name | NOTES |
| [任务类型]登录 | 1 | TASK\_TYPE\_LOG\_IN | 登录任务 |
| [任务类型]等级 | 1 | TASK\_TYPE\_LEVEL | 等级满足指定条件 |

