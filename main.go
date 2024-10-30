package main

import (
	"flag"
	"fmt"
	"os"

	"git.woa.com/modnarshen/excelconfc/compiler"
	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/types"
	"git.woa.com/modnarshen/excelconfc/util"
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
	outDir := flag.String("outdir", ".", "output directory")
	goPackage := flag.String("go_package", "excelconf", "target protobuf option go_package value")
	groupLabel := flag.String("group", "server", "filter fields with group label, the label could be 'server', 'client' or 'all'")
	addEnum := flag.Bool("add_enum", false, "add the enumeration values defined in the enumeration table to the current table output")
	flag.BoolVar(&rules.DEBUG_MODE, "debug", false, "DEBUG mode allows invalid output")
	flag.BoolVar(&util.NO_COLORFUL_LOG, "ncl", false, "`ncl` makes no colorful log output")
	flag.BoolVar(&util.VerboseMode, "verbose", false, "verbose mode show more debug information")

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
		groupFlag = groupFlag | types.GroupClient | types.GroupServer
	case "server":
		groupFlag = groupFlag | types.GroupServer
	case "client":
		groupFlag = groupFlag | types.GroupClient
	}

	if err := compiler.New(
		compiler.WithFilePath(*filePath),
		compiler.WithSheetName(*sheetName),
		compiler.WithEnumSheet(*enumSheet),
		compiler.WithOutDir(*outDir),
		compiler.WithGoPackage(*goPackage),
		compiler.WithAddEnum(*addEnum),
		compiler.WithGroupFlag(groupFlag),
	).Compile(); err != nil {
		util.LogError(err.Error())
	}
}
