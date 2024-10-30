/*
 * @Author: modnarshen
 * @Date: 2024.10.24 09:49:45
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package compiler

import (
	"fmt"

	"git.woa.com/modnarshen/excelconfc/compiler/mcc"
	"git.woa.com/modnarshen/excelconfc/types"
)

var OnReduce mcc.ReduceCallback = func(production *mcc.Production, nodeStack []mcc.ASTNode) ([]mcc.ASTNode, error) {
	stackTop := len(nodeStack)

	switch production.Number {
	case 1: // START -> FIELDS

	case 2: // FIELDS -> FIELDS FIELD
		nodeStack[stackTop-2].AddSubNode(nodeStack[stackTop-1])
		nodeStack = nodeStack[:stackTop-1]

	case 3: // FIELDS -> ε
		newASTNode := mcc.NewMiddleASTNode(types.MID_NODE_FIELDS)
		newASTNode.SetGroupFlag(types.GroupServer | types.GroupClient) // default flag is `all`
		nodeStack = append(nodeStack, newASTNode)

	case 4, 5, 6, 7: // FIELD -> BDT | ARRAY | STRUCT | VEC
		// use the leaf node as the filed node directly

	case 8, 9, 10, 11: // BDT -> int | float | string | enum
		// use the leaf node as the filed node directly

	case 12: // ARRAY -> array
		// treat the array leaf node as a `BDT`

	case 13: // STRUCT -> ADT { FIELDS }
		newASTNode := nodeStack[stackTop-2]
		subNode := nodeStack[stackTop-4]
		newASTNode.SetName(subNode.Name()).SetType(subNode.Type()).SetGroupFlag(subNode.GroupFlag())
		nodeStack = nodeStack[:stackTop-len(production.Right)]
		nodeStack = append(nodeStack, newASTNode)

	case 14: // VEC -> ADT VEC_ADT_ITEMS
		subNode := nodeStack[stackTop-2]
		newASTNode := mcc.NewMiddleASTNode(types.MID_NODE_VEC)
		newASTNode.SetName(subNode.Name()).SetType(subNode.Type()).SetGroupFlag(subNode.GroupFlag())
		for i, ssubNode := range nodeStack[stackTop-1].SubNodes() {
			ssubNode.SetName(fmt.Sprintf("%s[%d]", subNode.Name(), i)).SetType(subNode.Type())
			newASTNode.AddSubNode(ssubNode)
		}
		nodeStack = nodeStack[:stackTop-len(production.Right)]
		nodeStack = append(nodeStack, newASTNode)

	case 15: // VEC -> BDT VEC_BDT_ITEMS
		subNode := nodeStack[stackTop-2]
		newASTNode := mcc.NewMiddleASTNode(types.MID_NODE_VEC)
		newASTNode.SetName(subNode.Name()).SetType(subNode.Type()).SetGroupFlag(subNode.GroupFlag())
		for i, ssubNode := range nodeStack[stackTop-1].SubNodes() {
			ssubNode.SetName(fmt.Sprintf("%s[%d]", subNode.Name(), i)).SetType(subNode.Type())
			newASTNode.AddSubNode(ssubNode)
		}
		nodeStack = nodeStack[:stackTop-len(production.Right)]
		nodeStack = append(nodeStack, newASTNode)

	case 16: // ADT -> id
		subNode := nodeStack[stackTop-1]
		newASTNode := mcc.NewMiddleASTNode(types.MID_NODE_ADT)
		newASTNode.SetName(subNode.Name()).SetType(subNode.Type()).SetGroupFlag(subNode.GroupFlag())
		newASTNode.AddSubNode(subNode)
		nodeStack[stackTop-1] = newASTNode

	case 17: // VEC_ADT_ITEMS -> VEC_ADT_ITEMS [ FIELDS ]
		subNode := nodeStack[stackTop-2]
		nodeStack[stackTop-4].AddSubNode(subNode)
		nodeStack = nodeStack[:stackTop-3]

	case 18: // VEC_ADT_ITEMS -> [ FIELDS ]
		subNode := nodeStack[stackTop-2]
		newASTNode := mcc.NewMiddleASTNode(types.MID_NODE_VEC_ADT_ITEMS)
		newASTNode.AddSubNode(subNode)
		nodeStack = nodeStack[:stackTop-len(production.Right)]
		nodeStack = append(nodeStack, newASTNode)

	case 19: // VEC_BDT_ITEMS -> VEC_BDT_ITEMS []
		nodeStack[stackTop-2].AddSubNode(nodeStack[stackTop-1])
		nodeStack = nodeStack[:stackTop-1]

	case 20: // VEC_BDT_ITEMS -> []
		top := nodeStack[stackTop-1]
		newASTNode := mcc.NewMiddleASTNode(types.MID_NODE_VEC_BDT_ITEMS)
		newASTNode.SetType(top.Type())
		newASTNode.AddSubNode(top)
		nodeStack[stackTop-1] = newASTNode

	}

	return nodeStack, nil
}
