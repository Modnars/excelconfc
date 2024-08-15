package writer

import (
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"os"
	"path"

	"git.woa.com/modnarshen/excelconfc/util"
)

const (
	spaceStr    = "                                        " // len(spaceStr) == 40
	tabSpaceNum = 4
	outFilePerm = fs.FileMode(0644)
)

var (
	DEBUG_MODE = false

	intTypes    = util.NewSet("int32", "uint32")
	stringTypes = util.NewSet("string")
)

func isIntType(tp string) bool {
	return intTypes[tp]
}

func isStringType(tp string) bool {
	return stringTypes[tp]
}

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

func toOutBytes(output string) ([]byte, error) {
	var outBytes []byte
	if DEBUG_MODE {
		outBytes = []byte(output)
	} else {
		var err error
		outBytes, err = format.Source([]byte(output))
		if err != nil {
			util.LogError("format %s failed|err:{%+v}", outGoDefFileFullName, err)
			return nil, err
		}
	}
	return outBytes, nil
}

func WriteToFile(filePath string, writeBytes []byte) error {
	if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
		util.LogError("Failed to create file: %v", err)
		return err
	}
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
