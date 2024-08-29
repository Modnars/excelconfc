package translator

import "git.woa.com/modnarshen/excelconfc/types"

type DataHolder struct {
	types.DataHolder
	ASTRoot *Node
}
