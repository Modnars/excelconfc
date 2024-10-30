package writer

import (
	"fmt"
	"os"
)

func genOutFilePath(outDir string, fileName string, fileSuffix string) string {
	if outDir[len(outDir)-1] == '/' {
		outDir = outDir[:len(outDir)-1]
	}
	return fmt.Sprintf("%s/%s%s", outDir, fileName, fileSuffix)
}

func WriteToFile(outDir string, fileName string, fileSuffix string, writeBytes []byte) error {
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output dir -> %w", err)
	}
	filePath := genOutFilePath(outDir, fileName, fileSuffix)
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file -> %w", err)
	}
	defer outFile.Close()
	if _, err := outFile.Write(writeBytes); err != nil {
		return fmt.Errorf("write file failed -> %w", err)
	}
	return nil
}
