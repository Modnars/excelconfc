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
	Desc() string
	ColIdx() int
	GroupFlag() uint8
	SetGroupFlag(uint8) NodeInfo
}

type ASTNode interface {
	Lex
	NodeInfo
	AddSubNode(ASTNode) ASTNode
	SubNodes() []ASTNode
}

type astNode struct {
	lexVal    string
	subNodes  []ASTNode
	name      string
	nodeType  string
	desc      string
	colIdx    int
	groupFlag uint8
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

func (n *astNode) Desc() string {
	return n.desc
}

func (n *astNode) ColIdx() int {
	return n.colIdx
}

func (n *astNode) GroupFlag() uint8 {
	return n.groupFlag
}

func (n *astNode) SetGroupFlag(groupFlag uint8) NodeInfo {
	n.groupFlag = groupFlag
	return n
}

func (n *astNode) AddSubNode(subNode ASTNode) ASTNode {
	n.subNodes = append(n.subNodes, subNode)
	return n
}

func (n *astNode) SubNodes() []ASTNode {
	return n.subNodes
}

func (n *astNode) String() string {
	return fmt.Sprintf("%s name:%s type:%s colIdx:%d group:0b%b", n.lexVal, n.name, n.nodeType, n.colIdx, n.groupFlag)
}

func NewASTNode(lexVal string, name string, nodeType, desc string, colIdx int, groupFlag uint8) ASTNode {
	return &astNode{
		lexVal:    lexVal,
		name:      name,
		nodeType:  nodeType,
		desc:      desc,
		colIdx:    colIdx,
		groupFlag: groupFlag,
	}
}

func NewMiddleASTNode(lexVal string) ASTNode {
	return &astNode{lexVal: lexVal}
}

var (
	InputEndASTNode ASTNode = &astNode{lexVal: EndMark}
)

func PrintAST(node ASTNode, depth int) {
	fmt.Printf("%s%+v\n", util.IndentSpace(depth), node)
	for _, subNode := range node.SubNodes() {
		PrintAST(subNode, depth+1)
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
