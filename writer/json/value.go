package json

import (
	"fmt"
	"strconv"
	"strings"

	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
	"git.woa.com/modnarshen/excelconfc/translator"
	"git.woa.com/modnarshen/excelconfc/types"
)

type Field = translator.Node

func extractCellVal(cell string, asType string) (any, error) {
	switch asType {
	case types.TOK_TYPE_STRING, types.TOK_TYPE_DATETIME:
		return cell, nil
	case types.TOK_TYPE_BOOL:
		if cell == "" || cell == "0" || cell == types.MARK_VAL_FALSE {
			return false, nil
		} else {
			return true, nil
		}
	case types.TOK_TYPE_INT32:
		return strconv.ParseInt(cell, 10, 32)
	case types.TOK_TYPE_UINT32:
		return strconv.ParseUint(cell, 10, 32)
	case types.TOK_TYPE_INT64:
		return strconv.ParseInt(cell, 10, 64)
	case types.TOK_TYPE_UINT64:
		return strconv.ParseUint(cell, 10, 64)
	}
	return cell, nil
}

func CellValue(astNode mcc.ASTNode, cell string, evm types.EVM) (any, error) {
	if astNode.LexVal() == types.LEX_ENUM {
		if evm[cell] == nil {
			return nil, fmt.Errorf("enum label %s not found", cell)
		}
		return evm[cell].ID, nil
	} else if astNode.LexVal() == types.LEX_ARRAY {
		parts := strings.Split(cell, ";")
		result := []any{}
		for _, part := range parts {
			if val, err := extractCellVal(part, astNode.Type()); err != nil {
				return nil, err
			} else {
				result = append(result, val)
			}
		}
		return result, nil
	}
	return extractCellVal(cell, astNode.Type())
}
