package translator

import (
	"fmt"
	"strings"

	"git.woa.com/modnarshen/excelconfc/rules"
	"git.woa.com/modnarshen/excelconfc/types"
	"git.woa.com/modnarshen/excelconfc/util"
)

type Node struct {
	Name     string  // 字段名（解析后）
	Type     string  // 字段类型（解析后，可供解析时直接取用的类型）
	RawType  string  // 原始类型（例如 Excel 中原始配置的类型）
	Desc     string  // 修饰符（比如 Excel 配置中使用 D 来修饰 string 类型为时间类型）
	Group    string  // 分组（区分前台客户端、后台服务器等）
	ColIdx   int     // 列坐标，用于索引源数据
	SubNodes []*Node // 子字段（用于嵌套定义）

	structLabel string
}

func (n *Node) AddSubNode(subNode *Node) {
	n.SubNodes = append(n.SubNodes, subNode)
}

type NodeOption = func(node *Node)

func removeLabel(str string, labels ...string) (string, bool) {
	removed := false
	for _, label := range labels {
		if strings.Contains(str, label) {
			str = strings.ReplaceAll(str, label, "")
			removed = true
		}
	}
	return str, removed
}

func WithName(name string) NodeOption {
	return func(node *Node) {
		isFound := false
		if name, isFound = removeLabel(name, "|S", "|s"); isFound {
			node.Group += "S"
		}
		if name, isFound = removeLabel(name, "|C", "|c"); isFound {
			node.Group += "C"
		}

		if name, isFound = removeLabel(name, "["); isFound {
			node.structLabel += "["
		}
		if name, isFound = removeLabel(name, "{"); isFound {
			node.structLabel += "{"
		}
		if name, isFound = removeLabel(name, "]"); isFound {
			node.structLabel += "]"
		}
		if name, isFound = removeLabel(name, "}"); isFound {
			node.structLabel += "}"
		}

		node.Name = name
		if node.Group == "" {
			node.Group = "All"
		}
	}
}

func WithType(tp string) NodeOption {
	return func(node *Node) {
		node.Type = tp
		node.RawType = tp
	}
}

func WithDesc(desc string) NodeOption {
	return func(node *Node) {
		node.Desc = desc
		if desc == "D" && node.RawType == "string" {
			node.Type = "DateTime"
		}
		if desc == "E" {
			node.Type = "Enum"
		}
		if desc == "" && !types.IsBasicType(node.RawType) {
			node.Type = "Struct"
		}
	}
}

func WithColIdx(colIdx int) NodeOption {
	return func(node *Node) {
		node.ColIdx = colIdx
	}
}

func NewNode(options ...NodeOption) *Node {
	node := &Node{}
	for _, option := range options {
		option(node)
	}
	return node
}

func TransToNodes(headers [][]string) ([]*Node, error) {
	nodes := []*Node{}
	for colIdx := range headers[rules.ROW_IDX_NAME] {
		newFiled := NewNode(
			WithName(headers[rules.ROW_IDX_NAME][colIdx]),
			WithType(headers[rules.ROW_IDX_TYPE][colIdx]),
			WithDesc(headers[rules.ROW_IDX_DESC][colIdx]),
			WithColIdx(colIdx),
		)
		if newFiled != nil {
			nodes = append(nodes, newFiled)
		}
	}
	return nodes, nil
}

func (n *Node) isDescMatch(label string) bool {
	if n != nil && n.Desc == label {
		return true
	}
	return false
}

func (n *Node) IsVectorDecl() bool {
	return n.isDescMatch(types.TOK_DESC_VECTOR)
}

func (n *Node) IsStructDecl() bool {
	// return !types.IsBasicType(n.Type) && n.Type != types.TOK_PARSED_TYPE_DATETIME // && !n.IsEnum() TODO make sure enum's definition rule
	return n.Type == "Struct"
}

func (n *Node) IsEnum() bool {
	return n.isDescMatch(types.TOK_DESC_ENUM)
}

func (n *Node) isStructLabelMatch(label string) bool {
	if n != nil && n.structLabel == label {
		return true
	}
	return false
}

func (n *Node) IsVecNodeBegin() bool {
	return n.isStructLabelMatch(types.TOK_LF_SQ_BRACKET)
}

func (n *Node) IsVecNodeEnd() bool {
	return n.isStructLabelMatch(types.TOK_RG_SQ_BRACKET)
}

func (n *Node) IsStructNodeBegin() bool {
	return n.isStructLabelMatch(types.TOK_LF_CR_BRACKET)
}

func (n *Node) IsStructNodeEnd() bool {
	return n.isStructLabelMatch(types.TOK_RG_CR_BRACKET)
}

func BuildNodeTree(nodes []*Node) *Node {
	root := &Node{Name: "root", Type: "Struct"}
	structStack := util.Stack[*Node]{}
	structStack.Push(root)
	vecNodeStack := util.Stack[*Node]{}

	for idx, node := range nodes {
		if node.IsVectorDecl() {
			if idx+1 < len(nodes) && nodes[idx+1].IsVecNodeBegin() {
				vecNodeStack.Push(node)
			}
		} else if node.IsStructDecl() {
			if idx+1 < len(nodes) && nodes[idx+1].IsStructNodeBegin() {
				structStack.PeekOrZero().AddSubNode(node)
				structStack.Push(node)
				continue
			}
		}
		if node.IsVecNodeBegin() {
			structNode := &Node{Name: vecNodeStack.PeekOrZero().Name, Type: "Struct"}
			vecNodeStack.PeekOrZero().AddSubNode(structNode)
			structStack.Push(structNode)
		}

		structStack.PeekOrZero().AddSubNode(node)

		if node.IsStructNodeEnd() {
			structStack.Pop()
		}
		if node.IsVecNodeEnd() {
			structStack.Pop()
			if idx+1 < len(nodes) && !nodes[idx+1].IsVecNodeBegin() {
				vecNodeStack.Pop()
			}
		}
	}
	return root
}

var (
	spaces = "          "
)

func indentSpace(indent int) string {
	return spaces[:indent*2]
}

func PrintNodes(headers []*Node, indent int) {
	for _, header := range headers {
		if header != nil {
			fmt.Printf("%s%+v\n", indentSpace(indent), *header)
			PrintNodes(header.SubNodes, indent+1)
		}
	}
}

func PrintTree(node *Node, depth int) {
	fmt.Printf("%s%+v\n", indentSpace(depth), *node)
	for _, subNode := range node.SubNodes {
		PrintTree(subNode, depth+1)
	}
}
