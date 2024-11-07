# excelconfc

excelconfc (`Excel-Config-Compiler`) 是一个以《编译原理》为理论依据的配置文件生成工具。

## 安装及使用

- 克隆本项目到本地

```bash
git clone https://github.com/Modnars/excelconfc
```

- 编译项目

```bash
go build
```

- 使用工具

```bash
./excelconfc
```

- 安装到本地

```bash
sudo mv ./excelconfc /usr/local/bin
```

## 用法说明

使用 `./excelconfc` 或 `./excelconfc -help` 来查看使用说明。

```text
Usage: ./excelconfc [options]
Options:
  -add_enum
    	add the enumeration values defined in the enumeration table to the current table output
  -debug
    	DEBUG mode allows invalid output
  -enum_sheet string
    	enum definition sheet (default "ENUM_DESC")
  -excel string
    	target Excel config file path
  -go_out string
    	Generate Go source file.
  -go_package string
    	target protobuf option go_package value (default "excelconf")
  -group string
    	filter fields with group label, the label could be 'server', 'client' or 'all' (default "server")
  -json_out string
    	Generate Json source file.
  -ncl ncl
    	ncl makes no colorful log output
  -proto_out string
    	Generate Proto source file.
  -sheet string
    	target Excel config sheet
  -verbose
    	verbose mode show more details information
  -xml_out string
    	Generate XML source file.
```

其中，使用时必须提供 `-excel` 和 `-sheet` 参数，这两个参数分别指定需要导出的配置文件及文件中的配置表页签。在编译时，会自动寻找 `-enum_sheet` 指定的枚举表页签指定的枚举值定义并对其进行加载（如果没有此页签，目前的表现会直接报错）。

如果需要导出相关的枚举值定义到输出文件中（比如希望导出的 proto 文件中包含指定的枚举值定义，或导出的 Go 文件中是包含相关的枚举值以便在业务中直接使用这些枚举值定义），可以添加 `-add_enum` 参数，此时在导出相关的配置表的同时会将上述指定的定义的枚举值同时写入到导出文件中，因此这里需要注意的是，如果对于一个 Excel 文件中的多个 Sheet 页签进行导出，只需要在一个页签导出时指定 `-add_enum` 即可，否则会在多个导出文件中重复定义相同的枚举值。

在导出 Go 源码文件时，使用了 `go/format` 包来对生成的 Go 源码自动排版对齐，如果生成的 Go 源码包含语法错误（比如某些命名与语言本身保留字冲突等等），此时就可能会触发排版对齐报错。为了方便用户更快定位错误位置，使用 `-debug` 可以将直接生成的代码进行输出（无需进一步来排版对齐），由此来方便问题的定位，正式使用时，请不要指定此参数。同时，如果希望观察导表解析阶段的词法分析动作，可以指定 `-verbose` 参数来查看配置表结构的“规约”分析详情，同时，在接收后，也会输出其对应的语法分析树，由此更便于了解此工具是如何分析配置表结构的。

简单地说，一个配置表结构如果能够通过语法分析，那么证明此配置表结构能够被正确地解析成一棵语法分析树（注意，是“一棵”树，因此只会收敛到一个根结点），而这个根结点的词素值（LexVal，用于文法标识匹配）是 FIELDS 类型。有了这样的一个语法树的概念，后续的诸多内容就很方便理解了。

`-group` 选项允许用户进行分组导出。举个常见的生产实践的场景：某些配置表字段，我作为服务器开发，是不关心一些类似文字描述的客户端资源数据的。那么此时如果配置中存在一个名为 `xxxdesc` 的字段，我可能就希望导出的配置中忽略这个字段，以便我导出的配置文件更加精简，带来的好处就是在不影响我的业务需求的前提下，程序加载配置的效率更高。那么对于这样的字段，我仅需在字段名后添加 `|C`（或 `|c`）即可，C 表示 Client，这样的标识就表示这个字段仅需要给客户端导出时使用。此时，为了告诉工具我是一位服务器开发人员，在使用工具时添加 `-group=server` 即可，此时工具就会自动过滤掉那些不是 `server` 这个组的字段，从而保证分组功能的实现。值得一提的是，如果我在上述的语法分析树的中间结点使用了这个标识，那么以此中间结点为根结点的子树都会被裁剪掉，因此无需对这些子结点再进一步标识分组标记（这将极大地提高配置表的简洁性、直观性）。

`-go_package` 选项允许用户指定导出的 Go 源码、proto 文件的 package，默认为 `excelconf`。因为 protobuf 在生成 Go 源码的场景下，需要指定 `go_package` 选项，因此这里相当于直接将此选项的值写入到生成的 protobuf 文件中。对于形如 `go_package;package_name` 这样的选项值，直接生成的 Go 源码（`.ec.go` 文件）会截取 `package_name` 来作为包名。比如导表时选项：`-go_package=git.woa.com/modnarshen/uasvr-go/configs/excelconf;excelconf`，生成的 proto 文件中是 

```protobuf
option go_package = "git.woa.com/modnarshen/uasvr-go/configs/excelconf;excelconf";
```

生成的 Go 文件中是

```go
package excelconf
```

此工具作为一个命令行工具，默认用户会在终端环境下进行使用，因此对于错误信息等进行了颜色输出，这是通过 Linux 的终端字符颜色标记来进行实现的，如果需要将工具输出重定向到其他位置，不希望有颜色标记字符输出，可使用 `-ncl` 来显示指定不启用色彩标记。

最后，也是最重要的，对于形如 `-xxx_out` 标记而言，其支持导出数据到不同文件格式，目前已经支持的文件格式包括 `.proto`、`.go`、`.json` 和 `.xml`。其中前两者是作为结构定义而导出的。对于结构定义文件而言，其中不包括具体的配置数据，而只是必要的结构体定义代码。而 JSON 和 XML 则是作为数据格式文件进行输出，用于实际业务场景的数据加载。这几种文件的输出，都是采用上述的语法树结构来进行遍历，进一步分析语义动作，从而输出相应代码的，由此可以确保同一棵语法树生成的结构总是一致的。对于后续的输出结构，也会通过这个语法树来进行结构遍历输出。

## 其他文档

1. [设计原理](./docs/principle.md)
