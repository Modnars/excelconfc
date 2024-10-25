/*
 * @Author: modnarshen
 * @Date: 2024.10.16 11:11:56
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package mcc

import (
	"fmt"
	"strings"
)

const (
	NilMark = "ε"
	EndMark = "$"
)

type Production struct {
	Left   string
	Right  []string
	Number int
}

type Grammar struct {
	Productions []*Production
}

func (p *Production) Read(line string) error {
	parts := strings.Split(line, "->")
	if len(parts) != 2 {
		return fmt.Errorf("invalid production %s", line)
	}
	p.Left = strings.TrimSpace(parts[0])
	for _, part := range strings.Fields(parts[1]) {
		p.Right = append(p.Right, strings.TrimSpace(part))
	}
	return nil
}

func (p *Production) String() string {
	return fmt.Sprintf("No.%02d %s -> %s", p.Number, p.Left, strings.Join(p.Right, " "))
}

func (g *Grammar) Production(prodNo int) *Production {
	if prodNo <= 0 || prodNo > len(g.Productions) {
		return nil
	}
	return g.Productions[prodNo-1]
}

func (g *Grammar) Load(productionLines []string) error {
	productionNum := len(g.Productions)
	for _, productionLine := range productionLines {
		newProduction := &Production{}
		if err := newProduction.Read(productionLine); err != nil {
			return err
		}
		productionNum++
		newProduction.Number = productionNum
		g.Productions = append(g.Productions, newProduction)
	}
	return nil
}

func NewGrammar(productionLines []string) *Grammar {
	newGrammar := &Grammar{}
	if err := newGrammar.Load(productionLines); err != nil {
		panic("load production lines failed: " + err.Error())
	}
	return newGrammar
}
