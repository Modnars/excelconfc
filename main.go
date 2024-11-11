package main

import (
	"flag"
	"fmt"
	"os"

	"git.woa.com/modnarshen/excelconfc/compiler"
	"git.woa.com/modnarshen/excelconfc/lex"
	"git.woa.com/modnarshen/excelconfc/reader/xlsx"
	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/util"
	"git.woa.com/modnarshen/excelconfc/writer/golang"
	"git.woa.com/modnarshen/excelconfc/writer/json"
	"git.woa.com/modnarshen/excelconfc/writer/protobuf"
	"git.woa.com/modnarshen/excelconfc/writer/xml"
)

func makeSureArgv(isPass bool, errMsg string) {
	if !isPass {
		fmt.Fprintln(os.Stderr, errMsg)
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	filePath := flag.String("excel", "", "target Excel config file path")
	sheetName := flag.String("sheet", "", "target Excel config sheet")
	enumSheet := flag.String("enum_sheet", rules.DEFAULT_ENUM_SHEET_NAME, "enum definition sheet")
	containerLabel := flag.String("container", "map", "out data container type")
	protoOutDir := flag.String("proto_out", "", "Generate Proto source file.")
	goOutDir := flag.String("go_out", "", "Generate Go source file.")
	jsonOutDir := flag.String("json_out", "", "Generate Json source file.")
	xmlOutDir := flag.String("xml_out", "", "Generate XML source file.")
	goPackage := flag.String("go_package", "excelconf", "target protobuf option go_package value")
	groupLabel := flag.String("group", "server", "filter fields with group label, the label could be 'server', 'client' or 'all'")
	addEnum := flag.Bool("add_enum", false, "add the enumeration values defined in the enumeration table to the current table output")
	flag.BoolVar(&rules.DEBUG_MODE, "debug", false, "DEBUG mode allows invalid output")
	flag.BoolVar(&util.NO_COLORFUL_LOG, "ncl", false, "`ncl` makes no colorful log output")
	flag.BoolVar(&util.VerboseMode, "verbose", false, "verbose mode show more details information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	if !flag.Parsed() {
		flag.Parse()
	}

	makeSureArgv(*filePath != "", "Error: -excel is a required parameter")
	makeSureArgv(*sheetName != "", "Error: -sheet is a required parameter")

	groupFlag := uint8(0)
	switch *groupLabel {
	case "all":
		groupFlag = groupFlag | lex.GroupClient | lex.GroupServer
	case "server":
		groupFlag = groupFlag | lex.GroupServer
	case "client":
		groupFlag = groupFlag | lex.GroupClient
	}

	xlsxData, err := xlsx.ReadFile(*filePath, *sheetName, *enumSheet)
	if err != nil {
		util.LogError("exec ReadExcel failed|filePath:%s|sheetName:%s -> %v", *filePath, *sheetName, err)
		os.Exit(1)
	}

	switch *containerLabel {
	case "map":
		xlsxData.SetContainerType(rules.CONTAINER_TYPE_MAP)
	case "vec", "vector":
		xlsxData.SetContainerType(rules.CONTAINER_TYPE_VECTOR)
	}

	if err := compiler.New(
		compiler.WithFileName(*filePath),
		compiler.WithSheetName(*sheetName),
		compiler.WithGroupFlag(groupFlag),
	).Compile(xlsxData); err != nil {
		util.LogError(err.Error())
		os.Exit(1)
	}

	if *protoOutDir != "" {
		if err := protobuf.WriteToFile(xlsxData, *goPackage, *protoOutDir, *addEnum); err != nil {
			util.LogError("exec WriteToProto failed|filePath:%s|sheet:%s|outDir:%s -> %w", *filePath, *sheetName, *protoOutDir, err)
		}
	}

	if *goOutDir != "" {
		if err := golang.WriteToFile(xlsxData, *goPackage, *goOutDir, *addEnum); err != nil {
			util.LogError("exec golang.WriteToFile failed|file:%s|sheet:%s -> %w", xlsxData.FileName(), xlsxData.SheetName(), err)
		}
	}

	if *jsonOutDir != "" {
		if err := json.WriteToFile(xlsxData, *jsonOutDir); err != nil {
			util.LogError("exec json.WriteToFile failed|filePath:%s|sheet:%s|outDir:%s -> %w", *filePath, *sheetName, *jsonOutDir, err)
		}
	}

	if *xmlOutDir != "" {
		if err := xml.WriteToFile(xlsxData, *xmlOutDir); err != nil {
			util.LogError("exec xml.WriteToFile failed|filePath:%s|sheet:%s|outDir:%s -> %w", *filePath, *sheetName, *xmlOutDir, err)
		}
	}

}
