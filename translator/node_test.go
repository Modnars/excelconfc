package translator

import (
	"testing"
)

func TestRemoveAndCountMark(t *testing.T) {
	s := "[[abc]]"
	count := 0
	ws := "abc]]"
	s, count = removeAndCountMark(s, "[")
	if s != ws || count != 2 {
		t.Errorf("removeAndCountMark return (%s, %d), want (%s, %d)", s, count, ws, 2)
	}

	s = "abc]]"
	ws = "abc]]"
	s, count = removeAndCountMark(s, "[")
	if s != ws || count != 0 {
		t.Errorf("removeAndCountMark return (%s, %d), want (%s, %d)", s, count, ws, 0)
	}
}
