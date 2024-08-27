package main

import (
	"flag"
	"fmt"
	"os"

	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/translator"
	"git.woa.com/modnarshen/excelconfc/util"
)

func showErrorAndUsage(errMsg string) {
	fmt.Fprintln(os.Stderr, errMsg)
	flag.Usage()
}

func main() {
	filePath := flag.String("excel", "", "target Excel config file path")
	sheetName := flag.String("sheet", "", "target Excel config sheet")
	goPackage := flag.String("go_package", "excelconf", "target protobuf option go_package value")
	outDir := flag.String("outdir", ".", "output directory")
	flag.BoolVar(&rules.DEBUG_MODE, "debug", false, "DEBUG mode allows invalid output")
	flag.BoolVar(&util.NO_COLORFUL_LOG, "ncl", false, "`ncl` makes no colorful log output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	if !flag.Parsed() {
		flag.Parse()
	}

	if *filePath == "" {
		showErrorAndUsage("Error: -excel is a required parameter")
		os.Exit(1)
	}
	if *sheetName == "" {
		showErrorAndUsage("Error: -sheet is a required parameter")
		os.Exit(1)
	}

	if err := translator.Translate(*filePath, *sheetName, "ENUM_DESC", *goPackage, *outDir); err != nil {
		util.LogError(err.Error())
	}
}
