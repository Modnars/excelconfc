package writer

import (
	"fmt"
	"io"
	"os"
	"strings"
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
		return fmt.Errorf("failed to create output dir|err:%w", err)
	}
	filePath := genOutFilePath(outDir, fileName, fileSuffix)
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file|err:%w", err)
	}
	defer outFile.Close()
	if _, err := outFile.Write(writeBytes); err != nil {
		return fmt.Errorf("write file failed|err:%w", err)
	}
	return nil
}
