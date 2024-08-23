package writer

import (
	"fmt"
	"io"
	"os"
	"strings"

	"git.woa.com/modnarshen/excelconfc/util"
)

const (
	spaceStr    = "                                        " // len(spaceStr) == 40
	tabSpaceNum = 4
)

func indentSpace(indent int) string {
	return spaceStr[:indent*tabSpaceNum]
}

func wrf(wr io.Writer, format string, args ...any) {
	fmt.Fprintf(wr, format, args...)
}

func genOutFilePath(outDir string, fileName string, fileSuffix string) string {
	if outDir[len(outDir)-1] == '/' {
		outDir = outDir[:len(outDir)-1]
	}
	return fmt.Sprintf("%s/%s%s", outDir, fileName, fileSuffix)
}

func getPackageName(goPackage string) string {
	splitCh := ';'
	index := 0
	if strings.ContainsRune(goPackage, splitCh) {
		index = strings.IndexRune(goPackage, splitCh) + 1
	} else {
		index = strings.IndexRune(goPackage, '/') + 1
	}
	return goPackage[index:]
}

func WriteToFile(outDir string, fileName string, fileSuffix string, writeBytes []byte) error {
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		util.LogError("Failed to create file: %v", err)
		return err
	}
	filePath := genOutFilePath(outDir, fileName, fileSuffix)
	outFile, err := os.Create(filePath)
	if err != nil {
		util.LogError("Failed to create file: %v", err)
		return err
	}
	defer outFile.Close()
	if _, err := outFile.Write(writeBytes); err != nil {
		return err
	}
	return nil
}
