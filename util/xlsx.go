/*
 * @Author: modnarshen
 * @Date: 2024.10.30 17:20:56
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package util

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func ColumnName(colIdx int) string {
	name, err := excelize.ColumnNumberToName(colIdx + 1)
	if err != nil {
		return fmt.Sprint(colIdx)
	}
	return name
}
