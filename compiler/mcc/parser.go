/*
 * @Author: modnarshen
 * @Date: 2024.10.16 14:08:47
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package mcc

import "io"

type ReduceCallback func(*Production, []ASTNode) ([]ASTNode, error)

type Parser interface {
	Analyze([]Lex) error
	AnalyzeString(io.Reader) error
	BuildAST([]ASTNode, ReduceCallback) (ASTNode, error)
}
