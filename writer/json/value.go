package json

import (
	"fmt"

	"git.woa.com/modnarshen/excelconfc/translator"
	"git.woa.com/modnarshen/excelconfc/types"
)

type Field = translator.Node

func CellVal(field *Field, cell string) string {
	switch field.Type {
	case types.TOK_TYPE_STRING, types.TOK_PARSED_TYPE_DATETIME:
		return fmt.Sprintf(`"%s"`, cell)
	case types.TOK_TYPE_BOOL:
		if cell == "" || cell == "0" || cell == "false" {
			return "false"
		} else {
			return "true"
		}
	case types.TOK_TYPE_INT32:
		return cell
	case types.TOK_TYPE_UINT32:
		return cell
	case types.TOK_TYPE_INT64:
		return cell
	case types.TOK_TYPE_UINT64:
		return cell
	}
	return cell
}
