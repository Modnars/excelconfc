/*
 * @Author: modnarshen
 * @Date: 2024.10.16 11:31:04
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package mcc

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"git.woa.com/modnarshen/excelconfc/util"
)

type LRParser struct {
	grammar *Grammar
}

func (p *LRParser) Parse() error {
	return nil
}

func (p *LRParser) AnalyzeString(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)
		tokens = append(tokens, EndMark)
		lexVals := []Lex{}
		for _, token := range tokens {
			lexVals = append(lexVals, NewStringLex(token))
		}
		p.Analyze(lexVals)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner meets an error -> %w", err)
	}
	return nil
}

func (p *LRParser) Analyze(input []Lex) error {
	stateStack := util.Stack[int]{}
	stateStack.Push(0)
	var err error
	idx := 0
	for idx < len(input) {
		val := 0
		ok := true
		if val, ok = LRActionTable[stateStack.PeekOrZero()][input[idx].LexVal()]; !ok {
			err = fmt.Errorf("[SHIFT] (%d, %s)", stateStack.PeekOrZero(), input[idx].LexVal())
			break
		}
		if val >= 0 {
			fmt.Printf("[SHIFT] (%d, %s) -> %d\n", stateStack.PeekOrZero(), input[idx].LexVal(), val)
			stateStack.Push(val)
			idx++
		} else {
			reduceProduction := p.grammar.Production(-val)
			fmt.Printf("[REDUCE] (%d, %s) -> %d(%v)\n", stateStack.PeekOrZero(), input[idx].LexVal(), val, reduceProduction)
			if val == -1 && input[idx].LexVal() == EndMark {
				break
			}
			if reduceProduction.Right[0] != NilMark {
				for range reduceProduction.Right {
					stateStack.Pop()
				}
			}
			if val, ok = LRActionTable[stateStack.PeekOrZero()][reduceProduction.Left]; !ok {
				err = fmt.Errorf("[GOTO] (%d, %s)", stateStack.PeekOrZero(), reduceProduction.Left)
				break
			} else {
				fmt.Printf("[GOTO] (%d, %s) -> %d\n", stateStack.PeekOrZero(), reduceProduction.Left, val)
				stateStack.Push(val)
			}
		}
	}
	if err != nil {
		util.LogError("ANALYZE FAILED : %s", err.Error())
	} else {
		util.LogInfo("ACCEPT")
	}
	return err
}

func (p *LRParser) BuildAST(input []ASTNode, onReduce ReduceCallback) (ASTNode, error) {
	nodeStack := []ASTNode{}
	stateStack := util.Stack[int]{}
	stateStack.Push(0)
	var err error
	idx := 0
	input = append(input, NewMiddleASTNode(EndMark))
	util.LogTrace("input: %v", input)
	for idx < len(input) {
		val := 0
		ok := true
		if val, ok = LRActionTable[stateStack.PeekOrZero()][input[idx].LexVal()]; !ok {
			err = fmt.Errorf("[SHIFT] (%d, %s)", stateStack.PeekOrZero(), input[idx].LexVal())
			break
		}
		if val >= 0 {
			// fmt.Printf("[SHIFT] (%d, %s) -> %d\n", stateStack.PeekOrZero(), input[idx].LexVal(), val)
			stateStack.Push(val)
			nodeStack = append(nodeStack, input[idx])
			idx++
		} else {
			reduceProduction := p.grammar.Production(-val)
			util.LogTrace("[REDUCE] (%d, %s) -> %d (%v)", stateStack.PeekOrZero(), input[idx].LexVal(), val, reduceProduction)
			nodeStack, err = onReduce(reduceProduction, nodeStack) // call `onReduce`
			if err != nil {
				break
			}
			if val == -1 && input[idx].LexVal() == EndMark {
				break
			}
			if reduceProduction.Right[0] != NilMark {
				for range reduceProduction.Right {
					stateStack.Pop()
				}
			}
			if val, ok = LRActionTable[stateStack.PeekOrZero()][reduceProduction.Left]; !ok {
				err = fmt.Errorf("[GOTO] (%d, %s)", stateStack.PeekOrZero(), reduceProduction.Left)
				break
			} else {
				// fmt.Printf("[GOTO] (%d, %s) -> %d\n", stateStack.PeekOrZero(), reduceProduction.Left, val)
				stateStack.Push(val)
			}
		}
	}
	if err != nil {
		return nil, fmt.Errorf("build AST failed|remain:%v|err:%w", input[idx:], err)
	}
	if len(nodeStack) != 1 {
		return nil, fmt.Errorf("build AST failed|stackLen:%d", len(nodeStack))
	}
	return nodeStack[0], nil
}

func NewLRParser(grammar *Grammar) Parser {
	return &LRParser{grammar: grammar}
}
