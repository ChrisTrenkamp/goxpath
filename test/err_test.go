package test

import (
	"testing"

	"github.com/ChrisTrenkamp/goxpath/xpath"
)

func TestErr(t *testing.T) {
	xp := `/test/chil::p2`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><p1><p2/></p1></test>`
	_, err := xpath.FromStr(xp, x, nil)
	if err.Error() != "Invalid Axis specifier, chil" {
		t.Error("Invalid error message:" + err.Error())
	}
}
