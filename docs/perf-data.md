# excelconfc 性能数据

## 生成数据

Book1.xlsx 数据行数：102399 行

| 导出数据文件类型 | 用时 | 执行命令 |
| :-: | :-: | :-: |
| JSON | 64.95s user 3.94s system 117% cpu 58.665 total | `time ./excelconfc -excel=./testdata/Book1.xlsx -sheet=Sheet1 -go_package="git.woa.com/modnarshen/uasvr-go/configs/excelconf;excelconf" --group=server --verbose -go_out=./tmp -json_out=./tmp -proto_out=./tmp -add_enum` |
| XML | 65.62s user 4.18s system 116% cpu 59.725 total | `time ./excelconfc -excel=./testdata/Book1.xlsx -sheet=Sheet1 -go_package="git.woa.com/modnarshen/uasvr-go/configs/excelconf;excelconf" --group=server --verbose -go_out=./tmp -xml_out=./tmp -proto_out=./tmp -add_enum` |

## 加载数据

Book1.xlsx 数据行数：102399 行，重复 Key 行数

| 加载数据文件类型 | 用时 | 执行命令 |
| :-: | :-: | :-: |
| JSON | PASS: TestLoadFromHugeJson (1.97s) | `go test -run TestLoadFromHugeJson -v` |
| XML | PASS: TestLoadFromHugeXml (8.57s) | `go test -run TestLoadFromHugeXml -v` |

由此可见，Go 程序加载 JSON 格式的配置文件，确实比 XML 格式的配置文件更加高效
