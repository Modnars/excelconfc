package lex

import (
	"fmt"
	"strconv"
	"strings"

	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
)

func extractCellVal(cell string, asType string) (any, error) {
	switch asType {
	case TOK_TYPE_STRING, TOK_TYPE_DATETIME:
		return cell, nil
	case TOK_TYPE_BOOL:
		if cell == "" || cell == "0" || cell == TOK_VAL_FALSE {
			return false, nil
		} else {
			return true, nil
		}
	case TOK_TYPE_INT32:
		if cell == "" {
			return int32(0), nil
		}
		return strconv.ParseInt(cell, 10, 32)
	case TOK_TYPE_UINT32:
		if cell == "" {
			return uint32(0), nil
		}
		return strconv.ParseUint(cell, 10, 32)
	case TOK_TYPE_INT64:
		if cell == "" {
			return int64(0), nil
		}
		return strconv.ParseInt(cell, 10, 64)
	case TOK_TYPE_UINT64:
		if cell == "" {
			return uint64(0), nil
		}
		return strconv.ParseUint(cell, 10, 64)
	}
	return cell, nil
}

func CellValue(astNode mcc.ASTNode, cell string, evm EVM) (any, error) {
	if astNode.LexVal() == LEX_ENUM {
		if evm[cell] == nil {
			return nil, fmt.Errorf("enum label %s not found", cell)
		}
		return evm[cell].ID, nil
	} else if astNode.LexVal() == LEX_ARRAY {
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
