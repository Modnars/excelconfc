package main

import (
	"flag"
	"fmt"
	"os"

	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/translator"
)

func showErrorAndUsage(errMsg string) {
	fmt.Fprintln(os.Stderr, errMsg)
	flag.Usage()
}

func main() {
	filePath := flag.String("excel", "", "target Excel config file path")
	sheetName := flag.String("sheet", "", "target Excel config sheet")
	goPackage := flag.String("go_package", "excelconf", "target protobuf option go_package value")
	debugMode := flag.Bool("debug", false, "DEBUG mode allows invalid output")
	outDir := flag.String("outdir", ".", "output directory")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *filePath == "" {
		showErrorAndUsage("Error: -excel is a required parameter")
		os.Exit(1)
	}
	if *sheetName == "" {
		showErrorAndUsage("Error: -sheet is a required parameter")
		os.Exit(1)
	}

	if *debugMode {
		rules.DEBUG_MODE = true
	}

	translator.Translate(*filePath, *sheetName, *goPackage, *outDir)
}
