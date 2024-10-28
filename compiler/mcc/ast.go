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

type NodeInfo interface {
	Name() string
	SetName(string) NodeInfo
	Type() string
	SetType(string) NodeInfo
	ColIdx() int
}

type ASTNode interface {
	Lex
	NodeInfo
	AddSubNode(ASTNode) ASTNode
	SubNodes() []ASTNode
}

type astNode struct {
	lexVal   string
	subNodes []ASTNode
	name     string
	nodeType string
	colIdx   int
}

func (n *astNode) LexVal() string {
	return n.lexVal
}

func (n *astNode) Name() string {
	return n.name
}

func (n *astNode) SetName(name string) NodeInfo {
	n.name = name
	return n
}

func (n *astNode) Type() string {
	return n.nodeType
}

func (n *astNode) SetType(nodeType string) NodeInfo {
	n.nodeType = nodeType
	return n
}

func (n *astNode) ColIdx() int {
	return n.colIdx
}

func (n *astNode) AddSubNode(subNode ASTNode) ASTNode {
	n.subNodes = append(n.subNodes, subNode)
	return n
}

func (n *astNode) SubNodes() []ASTNode {
	return n.subNodes
}

func (n *astNode) String() string {
	return fmt.Sprintf("%s name:%s type:%s colIdx:%d", n.lexVal, n.name, n.nodeType, n.colIdx)
}

func NewASTNode(lexVal string, name string, nodeType string, colIdx int) ASTNode {
	return &astNode{lexVal: lexVal, name: name, nodeType: nodeType, colIdx: colIdx}
}

func NewMiddleASTNode(lexVal string) ASTNode {
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
