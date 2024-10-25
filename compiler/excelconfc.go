/*
 * @Author: modnarshen
 * @Date: 2024.10.24 09:49:45
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package compiler

import (
	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
)

var OnReduce mcc.ReduceCallback = func(production *mcc.Production, nodeStack []mcc.ASTNode) ([]mcc.ASTNode, error) {
	switch production.Number {
	case 1: // START -> FIELDS

	case 2: // FIELDS -> FIELDS FIELD
		nodeStack[len(nodeStack)-2].AddSubNode(nodeStack[len(nodeStack)-1])
		nodeStack = nodeStack[:len(nodeStack)-1]

	case 3: // FIELDS -> ε
		newASTNode := mcc.NewASTNode("Node@FIELDS")
		nodeStack = append(nodeStack, newASTNode)

	case 4: // FIELD -> BDT

	case 5: // FIELD -> ARRAY

	case 6: // FIELD -> STRUCT

	case 7: // FIELD -> VEC

	case 8, 9, 10, 11: // BDT -> int/float/string/enum
		newASTNode := mcc.NewASTNode("Node@BDT")
		newASTNode.AddSubNode(nodeStack[len(nodeStack)-1])
		nodeStack[len(nodeStack)-1] = newASTNode

	case 12: // ARRAY -> array
		newASTNode := mcc.NewASTNode("Node@ARRAY")
		newASTNode.AddSubNode(nodeStack[len(nodeStack)-1])
		nodeStack = nodeStack[:len(nodeStack)-1]
		nodeStack = append(nodeStack, newASTNode)

	case 13: // STRUCT -> ADT { FIELDS }
		newASTNode := mcc.NewASTNode("Node@STRUCT")
		newASTNode.AddSubNode(nodeStack[len(nodeStack)-4])
		newASTNode.AddSubNode(nodeStack[len(nodeStack)-2])
		nodeStack = nodeStack[:len(nodeStack)-len(production.Right)]
		nodeStack = append(nodeStack, newASTNode)

	case 14: // VEC -> ADT VEC_ADT_ITEMS
		newASTNode := mcc.NewASTNode("Node@VEC")
		newASTNode.AddSubNode(nodeStack[len(nodeStack)-2])
		newASTNode.AddSubNode(nodeStack[len(nodeStack)-1])
		nodeStack = nodeStack[:len(nodeStack)-len(production.Right)]
		nodeStack = append(nodeStack, newASTNode)

	case 15: // VEC -> BDT VEC_BDT_ITEMS
		newASTNode := mcc.NewASTNode("Node@VEC")
		newASTNode.AddSubNode(nodeStack[len(nodeStack)-2])
		newASTNode.AddSubNode(nodeStack[len(nodeStack)-1])
		nodeStack = nodeStack[:len(nodeStack)-len(production.Right)]
		nodeStack = append(nodeStack, newASTNode)

	case 16: // ADT -> id
		newASTNode := mcc.NewASTNode("Node@ADT")
		newASTNode.AddSubNode(nodeStack[len(nodeStack)-1])
		nodeStack[len(nodeStack)-1] = newASTNode

	case 17: // VEC_ADT_ITEMS -> VEC_ADT_ITEMS [ FIELDS ]
		nodeStack[len(nodeStack)-4].AddSubNode(nodeStack[len(nodeStack)-2])
		nodeStack = nodeStack[:len(nodeStack)-3]

	case 18: // VEC_ADT_ITEMS -> [ FIELDS ]
		newASTNode := mcc.NewASTNode("Node@VEC_ADT_ITEMS")
		newASTNode.AddSubNode(nodeStack[len(nodeStack)-2])
		nodeStack = nodeStack[:len(nodeStack)-len(production.Right)]
		nodeStack = append(nodeStack, newASTNode)

	case 19: // VEC_BDT_ITEMS -> VEC_BDT_ITEMS [ VEC_BDT_ITEMS ]
		nodeStack[len(nodeStack)-4].AddSubNode(nodeStack[len(nodeStack)-2])
		nodeStack = nodeStack[:len(nodeStack)-3]

	case 20: // VEC_BDT_ITEMS -> [ VEC_BDT_ITEMS ]
		top := nodeStack[len(nodeStack)-2]
		nodeStack = nodeStack[:len(nodeStack)-len(production.Right)]
		nodeStack = append(nodeStack, top)

	case 21: // VEC_BDT_ITEMS -> VEC_BDT_ITEMS []
		nodeStack[len(nodeStack)-2].AddSubNode(nodeStack[len(nodeStack)-1])
		nodeStack = nodeStack[:len(nodeStack)-1]

	case 22: // VEC_BDT_ITEMS -> []
		newASTNode := mcc.NewASTNode("Node@VEC_BDT_ITEMS")
		newASTNode.AddSubNode(nodeStack[len(nodeStack)-1])
		nodeStack[len(nodeStack)-1] = newASTNode

	}
	return nodeStack, nil
}
