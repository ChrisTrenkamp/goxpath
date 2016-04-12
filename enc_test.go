package goxpath

import (
	"testing"
)

func TestISO_8859_1(t *testing.T) {
	p := `/test`
	x := `<?xml version="1.0" encoding="iso-8859-1"?><test>test<path>path</path>test2</test>`
	exp := []string{"testpathtest2"}
	execVal(p, x, exp, nil, t)
}
