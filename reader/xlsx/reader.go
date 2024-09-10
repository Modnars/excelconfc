package xlsx

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/types"
	"git.woa.com/modnarshen/excelconfc/util"
	"github.com/xuri/excelize/v2"
)

const (
	patternEnumDeclType = `^\{([^\s]+)\}([^\s]+)$`
)

var (
	reEnumDeclType *regexp.Regexp
)

func readXlsxDataSheet(xlsxFile *excelize.File, dataSheetName string) ([][]string, [][]string, error) {
	rows, err := xlsxFile.GetRows(dataSheetName)
	if err != nil {
		return nil, nil, fmt.Errorf("get data sheet rows failed|sheet:%s -> %w", dataSheetName, err)
	}
	if len(rows) < rules.ROW_HEAD_MAX {
		return nil, nil, fmt.Errorf("invalid excel configuration format")
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

func readXlsxEnumSheet(xlsxFile *excelize.File, enumSheetName string) ([]*types.EnumTypeSt, map[string]*types.EnumValSt, error) {
	rows, err := xlsxFile.GetRows(enumSheetName)
	if err != nil {
		return nil, nil, fmt.Errorf("get enum sheet rows failed|sheet:%s -> %w", enumSheetName, err)
	}
	allEnumInfos := []*types.EnumTypeSt{}
	enumValMap := make(map[string]*types.EnumValSt)
	currEnumInfo := &types.EnumTypeSt{}
	currLabel := ""
	for i := 0; i < len(rows); i++ {
		for len(rows[i]) <= 0 {
			i++
		}
		matched := reEnumDeclType.FindStringSubmatch(rows[i][0])
		if len(matched) == 3 {
			allEnumInfos = append(allEnumInfos, currEnumInfo)
			currLabel = matched[2]
			currEnumInfo = &types.EnumTypeSt{Name: matched[1]}
			continue
		}
		if strings.HasPrefix(rows[i][0], "["+currLabel+"]") {
			newEnumVal := &types.EnumValSt{Name: rows[i][2], ID: rows[i][1]}
			currEnumInfo.EnumVals = append(currEnumInfo.EnumVals, newEnumVal)
			enumValMap[rows[i][0]] = newEnumVal
			continue
		}
	}
	allEnumInfos = append(allEnumInfos, currEnumInfo)
	return allEnumInfos[1:], enumValMap, nil
}

func ReadFile(filePath string, sheetName string, enumSheetName string) (types.DataHolder, error) {
	xlsxFile, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("excelize open file failed|filePath:%s|sheet:%s -> %w", filePath, sheetName, err)
	}
	defer func() {
		if err := xlsxFile.Close(); err != nil {
			util.LogError("close file failed|filePath:%s|sheet:%s -> {%+v}", filePath, sheetName, err)
		}
	}()

	headers, data, err := readXlsxDataSheet(xlsxFile, sheetName)
	if err != nil {
		return nil, fmt.Errorf("read data sheet failed|file:%s|sheet:%s -> %w", filepath.Base(filePath), sheetName, err)
	}
	enumTypes, enumValMap, err := readXlsxEnumSheet(xlsxFile, enumSheetName)
	if err != nil {
		return nil, fmt.Errorf("read enum sheet failed|file:%s|sheet:%s -> %w", filepath.Base(filePath), sheetName, err)
	}
	return NewDataHolder(
		WithFileName(filepath.Base(filePath)),
		WithSheetName(sheetName),
		WithHeaders(headers),
		WithData(data),
		WithEnumTypes(enumTypes),
		WithEnumValMap(enumValMap),
	), nil
}

func init() {
	reEnumDeclType = regexp.MustCompile(patternEnumDeclType)
}
