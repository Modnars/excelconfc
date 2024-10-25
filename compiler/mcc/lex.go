/*
 * @Author: modnarshen
 * @Date: 2024.10.23 16:41:10
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package mcc

type Lex interface {
	LexVal() string
}

type StringLex struct {
	lexVal string
}

func (l *StringLex) LexVal() string {
	return l.lexVal
}

func NewStringLex(lexVal string) *StringLex {
	return &StringLex{lexVal: lexVal}
}
