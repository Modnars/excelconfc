package reader

import (
	"fmt"

	"github.com/xuri/excelize/v2"

	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/util"
)

func ReadExcel(filePath string, sheetName string) ([][]string, [][]string, error) {
	excelFile, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("excelize open file failed|filePath:%s|sheet:%s|err:%w", filePath, sheetName, err)
	}

	defer func() {
		if err := excelFile.Close(); err != nil {
			util.LogError("close file failed|filePath:%s|sheet:%s|err:{%+v}", filePath, sheetName, err)
		}
	}()

	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return nil, nil, fmt.Errorf("get excel file rows failed|filePath:%s|sheet:%s|err:%w", filePath, sheetName, err)
	}
	if len(rows) < rules.ROW_HEAD_MAX {
		return nil, nil, fmt.Errorf("invalid excel configuration format|filePath:%s|sheetName:%s", filePath, sheetName)
	}
	// 确保 headers 中的每一行的元素个数都与第一行（名字行）元素个数相同
	headers := rows[:rules.ROW_HEAD_MAX]
	maxColNum := len(headers[rules.ROW_IDX_NAME])
	for i := 0; i < rules.ROW_HEAD_MAX; i += 1 {
		if len(headers[i]) < maxColNum {
			appendElementNum := maxColNum - len(headers[i])
			appendElements := make([]string, appendElementNum)
			for j := 0; j < appendElementNum; j += 1 {
				appendElements = append(appendElements, "")
			}
			headers[i] = append(headers[i], appendElements...)
		} else if len(headers[i]) > maxColNum {
			headers[i] = headers[i][:maxColNum]
		}
	}

	return headers, rows[rules.ROW_HEAD_MAX:], nil
}
