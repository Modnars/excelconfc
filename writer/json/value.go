package json

import (
	"fmt"

	"git.woa.com/modnarshen/excelconfc/translator"
	"git.woa.com/modnarshen/excelconfc/types"
)

type Field = translator.Node

func CellVal(field *Field, cell string, evm types.EVM) string {
	switch field.Type {
	case types.TOK_TYPE_STRING, types.TOK_TYPE_DATETIME:
		return fmt.Sprintf(`"%s"`, cell)
	case types.TOK_TYPE_BOOL:
		if cell == "" || cell == "0" || cell == types.MARK_VAL_FALSE {
			return types.TOK_VAL_FALSE
		} else {
			return types.TOK_VAL_TRUE
		}
	case types.TOK_TYPE_INT32:
		return cell
	case types.TOK_TYPE_UINT32:
		return cell
	case types.TOK_TYPE_INT64:
		return cell
	case types.TOK_TYPE_UINT64:
		return cell
	case types.TOK_TYPE_ENUM:
		return evm[cell].ID
	}
	return cell
}
