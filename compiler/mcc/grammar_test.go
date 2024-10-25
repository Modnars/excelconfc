/*
 * @Author: modnarshen
 * @Date: 2024.10.16 15:40:12
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package mcc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringsFields(t *testing.T) {
	testingStr := " a   bc d e	f"
	require.Equal(t, true, strings.Contains(testingStr, "\t"))
	strFields := strings.Fields(testingStr)
	require.Equal(t, 5, len(strFields))
}

func TestProductionRead(t *testing.T) {
	input := "START -> S T A R T"
	production := &Production{}
	err := production.Read(input)
	require.Equal(t, nil, err)
	require.Equal(t, 5, len(production.Right))
	require.Equal(t, "START", production.Left)
	require.Equal(t, []string{"S", "T", "A", "R", "T"}, production.Right)

	input = "S -> ε"
	production = &Production{}
	err = production.Read(input)
	require.Equal(t, nil, err)
	require.Equal(t, 1, len(production.Right))
	require.Equal(t, "S", production.Left)
	require.Equal(t, "ε", production.Right[0])
}
