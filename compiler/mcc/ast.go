/*
 * @Author: modnarshen
 * @Date: 2024.10.24 10:05:40
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package mcc

import (
	"fmt"

	"git.woa.com/modnarshen/excelconfc/util"
)

type ASTNode interface {
	Lex
	AddSubNode(ASTNode) error
	SubNodes() []ASTNode
}

type astNode struct {
	lexVal   string
	subNodes []ASTNode
}

func (n *astNode) LexVal() string {
	return n.lexVal
}

func (n *astNode) AddSubNode(subNode ASTNode) error {
	n.subNodes = append(n.subNodes, subNode)
	return nil
}

func (n *astNode) SubNodes() []ASTNode {
	return n.subNodes
}

func (n *astNode) String() string {
	return n.lexVal
}

func NewASTNode(lexVal string) ASTNode {
	return &astNode{lexVal: lexVal}
}

var _ ASTNode = (*astNode)(nil)

func PrintTree(node ASTNode, depth int) {
	fmt.Printf("%s%+v\n", util.IndentSpace(depth), node)
	for _, subNode := range node.SubNodes() {
		PrintTree(subNode, depth+1)
	}
}

func PrintASTNodes(nodes []ASTNode, depth int) {
	for _, node := range nodes {
		if node != nil {
			fmt.Printf("%s%+v\n", util.IndentSpace(depth), node)
			PrintASTNodes(node.SubNodes(), depth+1)
		}
	}
}
