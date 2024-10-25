/*
 * @Author: modnarshen
 * @Date: 2024.10.25 15:55:26
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package protobuf

import "git.woa.com/modnarshen/excelconfc/compiler/mcc"

func collectStruct() {
}

func visitASTNode(astNode mcc.ASTNode) error {
	for _, subNode := range astNode.SubNodes() {
		visitASTNode(subNode)
	}
	return nil
}
