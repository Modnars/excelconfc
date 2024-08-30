package translator

import (
	"fmt"

	"git.woa.com/modnarshen/excelconfc/util"
)

func PrintNodes(nodes []*Node, indent int) {
	for _, node := range nodes {
		if node != nil {
			fmt.Printf("%s%+v\n", util.IndentSpace(indent), *node)
			PrintNodes(node.SubNodes, indent+1)
		}
	}
}

func PrintTree(node *Node, depth int) {
	fmt.Printf("%s%+v\n", util.IndentSpace(depth), *node)
	for _, subNode := range node.SubNodes {
		PrintTree(subNode, depth+1)
	}
}
